package healthcheck

import (
	"net/http"
	"strings"
	"sync"

	"github.com/elliotxx/healthcheck/checks"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// NewHandler creates a new health check handler that can be used
// to check if the application is running.
func NewHandler(conf HandlerConfig) gin.HandlerFunc {
	var (
		mu            sync.Mutex
		failureInARow uint32
	)

	return func(c *gin.Context) {
		var (
			eg         errgroup.Group
			httpStatus = http.StatusOK
		)

		statuses := checks.NewCheckStatuses(len(conf.Checks))
		for _, check := range conf.Checks {
			captureCheck := check
			eg.Go(func() error {
				name := captureCheck.Name()

				if _, ok := statuses.Get(name); ok {
					return ErrHealthCheckNamesConflict
				}

				pass := captureCheck.Pass()
				statuses.Set(name, pass)

				if !pass {
					return ErrHealthCheckFailed
				}
				return nil
			})
		}

		// Wait for all the checks to complete.
		mu.Lock()
		if err := eg.Wait(); err != nil {
			httpStatus = http.StatusServiceUnavailable
			failureInARow++

			// Send a notification if the failure threshold is reached.
			if failureInARow >= conf.FailureNotification.Threshold &&
				conf.FailureNotification.Chan != nil {
				conf.FailureNotification.Chan <- err
			}
		} else if failureInARow != 0 && conf.FailureNotification.Chan != nil {
			// Reset the failure counter if all checks pass.
			failureInARow = 0
			conf.FailureNotification.Chan <- nil
		}
		mu.Unlock()

		if _, ok := c.GetQuery("verbose"); ok {
			conf.Verbose = true
		}
		if excludesStr, ok := c.GetQuery("excludes"); ok {
			conf.Excludes = strings.Split(excludesStr, ",")
		}

		c.String(httpStatus, statuses.String(conf.Verbose, conf.Excludes))
	}
}

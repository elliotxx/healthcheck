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

		// Process the request parameters.
		verbose := conf.Verbose
		excludes := conf.Excludes
		if _, ok := c.GetQuery("verbose"); ok {
			verbose = true
		}
		if excludesStr, ok := c.GetQuery("excludes"); ok {
			excludes = strings.Split(excludesStr, ",")
		}

		// Create a new check statuses instance.
		statuses := checks.NewCheckStatuses(len(conf.Checks))

		// Iterate over the list of health checks and execute them
		// concurrently.
		for _, check := range conf.Checks {
			// Capture the check variable to avoid race conditions.
			captureCheck := check

			eg.Go(func() error {
				// Get the name of the check and check if it already
				// exists in the statuses list.
				name := captureCheck.Name()

				if len(excludes) > 0 {
					for _, excludedName := range excludes {
						if excludedName == name {
							return nil
						}
					}
				}

				if _, ok := statuses.Get(name); ok {
					return ErrHealthCheckNamesConflict
				}

				// Execute the check and update the status list.
				pass := captureCheck.Pass()
				statuses.Set(name, pass)

				// If the check fails, return a failure error.
				if !pass {
					return ErrHealthCheckFailed
				}
				return nil
			})
		}

		// Wait for all the checks to complete.
		mu.Lock()
		if err := eg.Wait(); err != nil {
			// If any of the checks fail, set the HTTP status code to service
			// unavailable.
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

		// Return the status response as a string.
		c.String(httpStatus, statuses.String(verbose))
	}
}

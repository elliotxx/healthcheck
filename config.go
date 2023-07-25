package healthcheck

import (
	"net/http"

	"github.com/elliotxx/healthcheck/checks"
)

// Config represents the configuration for the health check.
type Config struct {
	HTTPMethod string
	Endpoint   string
	HandlerConfig
}

// HandlerConfig represents the configuration for the health check
// handler.
type HandlerConfig struct {
	Verbose  bool
	Excludes []string
	Checks   []checks.Check
	FailureNotification
}

// FailureNotification represents the configuration for failure
// notifications.
type FailureNotification struct {
	Threshold uint32
	Chan      chan error
}

// NewDefaultConfig creates a new default health check configuration.
func NewDefaultConfig() Config {
	c := NewDefaultHandlerConfig()
	return Config{
		HTTPMethod:    http.MethodGet,
		Endpoint:      "/healthz",
		HandlerConfig: c,
	}
}

// NewDefaultHandlerConfig creates a new default health check handler
// configuration.
func NewDefaultHandlerConfig() HandlerConfig {
	return HandlerConfig{
		Verbose:             false,
		Excludes:            []string{},
		Checks:              []checks.Check{checks.NewPingCheck()},
		FailureNotification: FailureNotification{Threshold: 1},
	}
}

// NewDefaultConfigFor creates a new default health check
// configuration for the specified checks.
func NewDefaultConfigFor(cs ...checks.Check) Config {
	c := NewDefaultHandlerConfigFor(cs...)
	return Config{
		HTTPMethod:    http.MethodGet,
		Endpoint:      "/healthz",
		HandlerConfig: c,
	}
}

// NewDefaultHandlerConfigFor creates a new default health check
// handler configuration for the specified checks.
func NewDefaultHandlerConfigFor(cs ...checks.Check) HandlerConfig {
	if len(cs) == 0 {
		cs = []checks.Check{}
	}

	return HandlerConfig{
		Verbose:             false,
		Excludes:            []string{},
		Checks:              cs,
		FailureNotification: FailureNotification{Threshold: 1},
	}
}

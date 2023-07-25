package healthcheck

import (
	"net/http"

	"github.com/elliotxx/healthcheck/checks"
)

type Config struct {
	HTTPMethod string
	Endpoint   string
	HandlerConfig
}

type HandlerConfig struct {
	Verbose  bool
	Excludes []string
	Checks   []checks.Check
	FailureNotification
}

type FailureNotification struct {
	Threshold uint32
	Chan      chan error
}

func NewDefaultConfig() Config {
	c := NewDefaultHandlerConfig()
	return Config{
		HTTPMethod:    http.MethodGet,
		Endpoint:      "/healthz",
		HandlerConfig: c,
	}
}

func NewDefaultHandlerConfig() HandlerConfig {
	return HandlerConfig{
		Verbose:             false,
		Excludes:            []string{},
		Checks:              []checks.Check{checks.NewPingCheck()},
		FailureNotification: FailureNotification{Threshold: 1},
	}
}

func NewDefaultConfigFor(cs ...checks.Check) Config {
	c := NewDefaultHandlerConfigFor(cs...)
	return Config{
		HTTPMethod:    http.MethodGet,
		Endpoint:      "/healthz",
		HandlerConfig: c,
	}
}

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

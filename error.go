package healthcheck

import (
	"errors"
)

var (
	ErrHealthCheckFailed             = errors.New("health check failed")
	ErrHealthCheckNamesConflict      = errors.New("health check names conflict")
	ErrConfigsMethodEndpointConflict = errors.New("configs method and endpoint conflict")
	ErrEmptyConfigs                  = errors.New("empty configs")
)

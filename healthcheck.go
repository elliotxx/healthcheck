package healthcheck

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Register registers a default health check handler with the
// specified router group.
func Register(r *gin.RouterGroup) error {
	return registerFor(r, NewDefaultConfig())
}

// RegisterFor registers one or more health check handlers with the
// specified router group and configurations.
func RegisterFor(r *gin.RouterGroup, configs ...Config) error {
	return registerFor(r, configs...)
}

func registerFor(r *gin.RouterGroup, configs ...Config) error {
	if len(configs) == 0 {
		return ErrEmptyConfigs
	}

	existedMethodEndpoints := map[string]struct{}{}
	for _, conf := range configs {
		// check if the method and endpoint are duplicated.
		key := strings.ToLower(conf.HTTPMethod + conf.Endpoint)
		if _, existed := existedMethodEndpoints[key]; existed {
			return ErrConfigsMethodEndpointConflict
		}
		existedMethodEndpoints[key] = struct{}{}

		// register the health check handler.
		r.Handle(conf.HTTPMethod, conf.Endpoint, NewHandler(conf.HandlerConfig))
	}

	return nil
}

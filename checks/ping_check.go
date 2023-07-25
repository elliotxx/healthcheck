package checks

// pingCheck is a simple health check that always returns true.
type pingCheck struct{}

// NewPingCheck creates a new ping health check.
func NewPingCheck() Check {
	return &pingCheck{}
}

// Pass always returns true for the ping health check.
func (c *pingCheck) Pass() bool {
	return true
}

// Name returns the name of the ping health check.
func (c *pingCheck) Name() string {
	return "Ping"
}

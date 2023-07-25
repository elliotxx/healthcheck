package checks

// pingCheck is a simple check that returns true.
type pingCheck struct{}

func NewPingCheck() Check {
	return &pingCheck{}
}

func (c *pingCheck) Pass() bool {
	return true
}

func (c *pingCheck) Name() string {
	return "Ping"
}

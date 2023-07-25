package checks

import (
	"os"
	"regexp"
)

// NewEnvCheck creates a new environment variable check.
func NewEnvCheck(envVariable string) Check {
	return &envCheck{
		envVariable: envVariable,
		regex:       "",
	}
}

// NewEnvRegexCheck creates a new environment variable check with a
// regular expression.
func NewEnvRegexCheck(envVariable string, r string) Check {
	return &envCheck{
		envVariable: envVariable,
		regex:       r,
	}
}

// envCheck is a struct for defining an environment variable check.
type envCheck struct {
	envVariable string
	regex       string
}

// Pass checks if the environment variable passes the check.
func (e *envCheck) Pass() bool {
	envValue := os.Getenv(e.envVariable)
	if envValue == "" {
		return false
	}
	if e.regex != "" {
		matched, _ := regexp.MatchString(e.regex, envValue)
		return matched
	}
	return true
}

// Name returns the name of the environment variable check.
func (e *envCheck) Name() string {
	return "Env-" + e.envVariable
}

package checks

import (
	"os"
	"regexp"
)

func NewEnvCheck(envVariable string) Check {
	return &envCheck{
		envVariable: envVariable,
		regex:       "",
	}
}

func NewEnvRegexCheck(envVariable string, r string) Check {
	return &envCheck{
		envVariable: envVariable,
		regex:       r,
	}
}

type envCheck struct {
	envVariable string
	regex       string
}

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

func (e *envCheck) Name() string {
	return "Env-" + e.envVariable
}

package checks

import (
	"sort"
	"strings"
	"sync"
)

type Check interface {
	Pass() bool
	Name() string
}

type CheckStatuses interface {
	Get(k string) (bool, bool)
	Set(k string, v bool)
	Delete(k string)
	Len() int
	Each(f func(k string, v bool))
	String(verbose bool, excludes []string) string
}

var _ CheckStatuses = &checkStatuses{}

type checkStatuses struct {
	sync.RWMutex
	m map[string]bool
}

func NewCheckStatuses(n int) CheckStatuses {
	return &checkStatuses{
		m: make(map[string]bool, n),
	}
}

func (cs *checkStatuses) Get(k string) (bool, bool) {
	cs.RLock()
	defer cs.RUnlock()
	v, existed := cs.m[k]
	return v, existed
}

func (cs *checkStatuses) Set(k string, v bool) {
	cs.Lock()
	defer cs.Unlock()
	cs.m[k] = v
}

func (cs *checkStatuses) Delete(k string) {
	cs.Lock()
	defer cs.Unlock()
	delete(cs.m, k)
}

func (cs *checkStatuses) Len() int {
	cs.RLock()
	defer cs.RUnlock()

	return len(cs.m)
}

func (cs *checkStatuses) Each(f func(k string, v bool)) {
	cs.RLock()
	defer cs.RUnlock()

	for k, v := range cs.m {
		f(k, v)
	}
}

func (cs *checkStatuses) String(verbose bool, excludes []string) string {
	if len(excludes) > 0 {
		for _, checkName := range excludes {
			cs.Delete(checkName)
		}
	}

	passNames := make([]string, 0, cs.Len())
	failedNames := make([]string, 0, cs.Len())
	allPass := true
	cs.Each(func(name string, pass bool) {
		if pass {
			passNames = append(passNames, name)
		} else {
			failedNames = append(failedNames, name)
			allPass = false
		}
	})

	sort.Strings(passNames)
	sort.Strings(failedNames)

	if verbose {
		var b strings.Builder
		for _, name := range passNames {
			b.WriteString("[+] " + name + " ok\n")
		}

		for _, name := range failedNames {
			b.WriteString("[-] " + name + " fail\n")
		}

		if allPass {
			b.WriteString("health check passed")
		} else {
			b.WriteString("health check failed")
		}

		return b.String()
	}

	if allPass {
		return "OK"
	}
	return "Fail"
}

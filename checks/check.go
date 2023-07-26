package checks

import (
	"sort"
	"strings"
	"sync"
)

// Check represents a health check that can be run to check the status
// of a service.
type Check interface {
	Pass() bool
	Name() string
}

// CheckStatuses is an interface for a thread-safe map of check
// statuses.
type CheckStatuses interface {
	Get(k string) (bool, bool)
	Set(k string, v bool)
	Delete(k string)
	Len() int
	Each(f func(k string, v bool))
	String(verbose bool) string
}

var (
	_ CheckStatuses = &checkStatuses{}

	// stringBuilderPool is a sync.Pool used to reuse strings.Builder
	// instances.
	// Reusing builders can help to avoid frequent memory allocations,
	// which can improve module performance.
	stringBuilderPool = sync.Pool{
		New: func() any {
			return &strings.Builder{}
		},
	}
)

// checkStatuses is a thread-safe implements of CheckStatuses
// interface.
type checkStatuses struct {
	sync.RWMutex
	m map[string]bool
}

// NewCheckStatuses creates a new CheckStatuses instance with the
// specified capacity.
func NewCheckStatuses(n int) CheckStatuses {
	return &checkStatuses{
		m: make(map[string]bool, n),
	}
}

// Get returns the value and existence status for the specified key.
func (cs *checkStatuses) Get(k string) (bool, bool) {
	cs.RLock()
	defer cs.RUnlock()
	v, existed := cs.m[k]
	return v, existed
}

// Set sets the value for the specified key.
func (cs *checkStatuses) Set(k string, v bool) {
	cs.Lock()
	defer cs.Unlock()
	cs.m[k] = v
}

// Delete deletes the specified key from the map.
func (cs *checkStatuses) Delete(k string) {
	cs.Lock()
	defer cs.Unlock()
	delete(cs.m, k)
}

// Len returns the number of items in the map.
func (cs *checkStatuses) Len() int {
	cs.RLock()
	defer cs.RUnlock()

	return len(cs.m)
}

// Each calls the specified function for each key/value pair in the
// map.
func (cs *checkStatuses) Each(f func(k string, v bool)) {
	cs.RLock()
	defer cs.RUnlock()

	for k, v := range cs.m {
		f(k, v)
	}
}

// String returns a string representation of the check statuses.
// If verbose is true, the output includes pass/fail status for each
// check.
func (cs *checkStatuses) String(verbose bool) string {
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
		b := stringBuilderPool.Get().(*strings.Builder)
		defer stringBuilderPool.Put(b)
		defer b.Reset()

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

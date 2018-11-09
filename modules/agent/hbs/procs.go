package hbs

import (
	"sync"
)

var (
	// tags => {1=>name, 2=>cmdline}
	// e.g. 'name=falcon-agent'=>{1=>falcon-agent}
	// e.g. 'cmdline=xx'=>{2=>xx}
	procs     map[string]map[int]string
	procsLock = new(sync.RWMutex)
)

func ReportProcs() map[string]map[int]string {
	procsLock.RLock()
	defer procsLock.RUnlock()
	return procs
}

func CacheReportProcs(value map[string]map[int]string) {
	procsLock.Lock()
	defer procsLock.Unlock()
	procs = value
}

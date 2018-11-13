package hbs

import (
	"sync"
)

var (
	files     map[string]int
	filesLock = new(sync.RWMutex)
)

func ReportFiles() map[string]int {
	filesLock.RLock()
	defer filesLock.RUnlock()
	return files
}

func CacheReportFiles(value map[string]int) {
	filesLock.Lock()
	defer filesLock.Unlock()
	files = value
}

package hbs

import (
	"sync"
)

var (
	paths     []string
	pathsLock = new(sync.RWMutex)
)

func DuPaths() []string {
	pathsLock.RLock()
	defer pathsLock.RUnlock()
	return paths
}

func CacheDuPaths(value []string) {
	pathsLock.Lock()
	defer pathsLock.Unlock()
	paths = value
}

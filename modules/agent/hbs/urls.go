package hbs

import (
	"sync"
)

var (
	urls     map[string]string
	urlsLock = new(sync.RWMutex)
)

func ReportUrls() map[string]string {
	urlsLock.RLock()
	defer urlsLock.RUnlock()
	return urls
}

func CacheReportUrls(value map[string]string) {
	urlsLock.RLock()
	defer urlsLock.RUnlock()
	urls = value
}

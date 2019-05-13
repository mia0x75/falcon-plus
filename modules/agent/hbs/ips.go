package hbs

import (
	"strings"
	"sync"

	"github.com/toolkits/slice"
)

var (
	ips     []string
	ipsLock = new(sync.Mutex)
)

// TrustableIps TODO:
func TrustableIps() []string {
	ipsLock.Lock()
	defer ipsLock.Unlock()
	return ips
}

// CacheTrustableIps TODO:
func CacheTrustableIps(ipStr string) {
	arr := strings.Split(ipStr, ",")
	ipsLock.Lock()
	defer ipsLock.Unlock()
	ips = arr
}

// IsTrustable TODO:
func IsTrustable(remoteAddr string) bool {
	ip := remoteAddr
	idx := strings.LastIndex(remoteAddr, ":")
	if idx > 0 {
		ip = remoteAddr[0:idx]
	}

	if ip == "127.0.0.1" {
		return true
	}

	return slice.ContainsString(TrustableIps(), ip)
}

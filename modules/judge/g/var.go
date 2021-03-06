package g

import (
	"sync"
	"time"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

// SafeStrategyMap TODO:
type SafeStrategyMap struct {
	sync.RWMutex
	// endpoint:metric => [strategy1, strategy2 ...]
	M map[string][]cm.Strategy
}

// SafeExpressionMap TODO:
type SafeExpressionMap struct {
	sync.RWMutex
	// metric:tag1 => [exp1, exp2 ...]
	// metric:tag2 => [exp1, exp2 ...]
	M map[string][]*cm.Expression
}

// SafeEventMap TODO:
type SafeEventMap struct {
	sync.RWMutex
	M map[string]*cm.Event
}

// SafeFilterMap TODO:
type SafeFilterMap struct {
	sync.RWMutex
	M map[string]string
}

// TODO:
var (
	HBSClient     *SingleConnRPCClient
	StrategyMap   = &SafeStrategyMap{M: make(map[string][]cm.Strategy)}
	ExpressionMap = &SafeExpressionMap{M: make(map[string][]*cm.Expression)}
	LastEvents    = &SafeEventMap{M: make(map[string]*cm.Event)}
	FilterMap     = &SafeFilterMap{M: make(map[string]string)}
)

// InitHbsClient TODO:
func InitHbsClient() {
	HBSClient = &SingleConnRPCClient{
		RPCServers:  Config().HBS.Servers,
		Timeout:     time.Duration(Config().HBS.Timeout) * time.Millisecond,
		CallTimeout: time.Duration(3000) * time.Millisecond,
	}
}

// ReInit TODO:
func (cache *SafeStrategyMap) ReInit(m map[string][]cm.Strategy) {
	cache.Lock()
	defer cache.Unlock()
	cache.M = m
}

// Get TODO:
func (cache *SafeStrategyMap) Get() map[string][]cm.Strategy {
	cache.RLock()
	defer cache.RUnlock()
	return cache.M
}

// ReInit TODO:
func (cache *SafeExpressionMap) ReInit(m map[string][]*cm.Expression) {
	cache.Lock()
	defer cache.Unlock()
	cache.M = m
}

// Get TODO:
func (cache *SafeExpressionMap) Get() map[string][]*cm.Expression {
	cache.RLock()
	defer cache.RUnlock()
	return cache.M
}

// Get TODO:
func (cache *SafeEventMap) Get(key string) (*cm.Event, bool) {
	cache.RLock()
	defer cache.RUnlock()
	event, exists := cache.M[key]
	return event, exists
}

// Set TODO:
func (cache *SafeEventMap) Set(key string, event *cm.Event) {
	cache.Lock()
	defer cache.Unlock()
	cache.M[key] = event
}

// ReInit TODO:
func (cache *SafeFilterMap) ReInit(m map[string]string) {
	cache.Lock()
	defer cache.Unlock()
	cache.M = m
}

// Exists TODO:
func (cache *SafeFilterMap) Exists(key string) bool {
	cache.RLock()
	defer cache.RUnlock()
	if _, ok := cache.M[key]; ok {
		return true
	}
	return false
}

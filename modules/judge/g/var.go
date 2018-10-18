package g

import (
	"sync"
	"time"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

type SafeStrategyMap struct {
	sync.RWMutex
	// endpoint:metric => [strategy1, strategy2 ...]
	M map[string][]cmodel.Strategy
}

type SafeExpressionMap struct {
	sync.RWMutex
	// metric:tag1 => [exp1, exp2 ...]
	// metric:tag2 => [exp1, exp2 ...]
	M map[string][]*cmodel.Expression
}

type SafeEventMap struct {
	sync.RWMutex
	M map[string]*cmodel.Event
}

type SafeFilterMap struct {
	sync.RWMutex
	M map[string]string
}

var (
	HbsClient     *SingleConnRpcClient
	StrategyMap   = &SafeStrategyMap{M: make(map[string][]cmodel.Strategy)}
	ExpressionMap = &SafeExpressionMap{M: make(map[string][]*cmodel.Expression)}
	LastEvents    = &SafeEventMap{M: make(map[string]*cmodel.Event)}
	FilterMap     = &SafeFilterMap{M: make(map[string]string)}
)

func InitHbsClient() {
	HbsClient = &SingleConnRpcClient{
		RpcServers: Config().Hbs.Servers,
		Timeout:    time.Duration(Config().Hbs.Timeout) * time.Millisecond,
	}
}

func (this *SafeStrategyMap) ReInit(m map[string][]cmodel.Strategy) {
	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeStrategyMap) Get() map[string][]cmodel.Strategy {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

func (this *SafeExpressionMap) ReInit(m map[string][]*cmodel.Expression) {
	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeExpressionMap) Get() map[string][]*cmodel.Expression {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

func (this *SafeEventMap) Get(key string) (*cmodel.Event, bool) {
	this.RLock()
	defer this.RUnlock()
	event, exists := this.M[key]
	return event, exists
}

func (this *SafeEventMap) Set(key string, event *cmodel.Event) {
	this.Lock()
	defer this.Unlock()
	this.M[key] = event
}

func (this *SafeFilterMap) ReInit(m map[string]string) {
	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeFilterMap) Exists(key string) bool {
	this.RLock()
	defer this.RUnlock()
	if _, ok := this.M[key]; ok {
		return true
	}
	return false
}

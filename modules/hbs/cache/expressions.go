package cache

import (
	"sync"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
)

type SafeExpressionCache struct {
	sync.RWMutex
	L []*cmodel.Expression
}

var ExpressionCache = &SafeExpressionCache{}

func (this *SafeExpressionCache) Get() []*cmodel.Expression {
	this.RLock()
	defer this.RUnlock()
	return this.L
}

func (this *SafeExpressionCache) Init() {
	es, err := db.QueryExpressions()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.L = es
}

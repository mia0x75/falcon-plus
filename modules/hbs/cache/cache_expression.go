package cache

import (
	"sync"

	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
)

type SafeExpressionCache struct {
	sync.RWMutex
	L []*cm.Expression
}

var ExpressionCache = &SafeExpressionCache{}

func (m *SafeExpressionCache) Get() []*cm.Expression {
	m.RLock()
	defer m.RUnlock()
	return m.L
}

func (m *SafeExpressionCache) Init() {
	es, err := db.QueryExpressions()
	if err != nil {
		return
	}

	m.Lock()
	defer m.Unlock()
	m.L = es
}

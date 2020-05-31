package cache

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/api/app/model"
)

// SafeExpressionsMap 线程安全的数据缓存对象
type SafeExpressionsMap struct {
	sync.RWMutex
	M []*model.Expression
}

// ExpressionsMap 群集缓存对象
var ExpressionsMap = &SafeExpressionsMap{}

// Count 返回缓存条数
func (c *SafeExpressionsMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Append 添加元素
func (c *SafeExpressionsMap) Append(item *model.Expression) {
	c.RLock()
	defer c.RUnlock()
	c.M = append(c.M, item)
}

// Remove 删除元素，每次仅删除一个
func (c *SafeExpressionsMap) Remove(f func(*model.Expression) bool) {
	c.RLock()
	defer c.RUnlock()
	for i, host := range c.M {
		if f(host) {
			c.M = append(c.M[:i], c.M[i+1:]...)
			break
		}
	}
}

// Include returns true if one of the element in the slice satisfies the predicate f.
func (c *SafeExpressionsMap) Include(f func(*model.Expression) bool) bool {
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.M {
		if f(v) {
			return true
		}
	}
	return false
}

// Any returns the element if one of the element in the slice satisfies the predicate f.
func (c *SafeExpressionsMap) Any(f func(*model.Expression) bool) *model.Expression {
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.M {
		if f(v) {
			return v
		}
	}
	return nil
}

// All returns all of the slice.
func (c *SafeExpressionsMap) All() []*model.Expression {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeExpressionsMap) Filter(f func(*model.Expression) bool) (L []*model.Expression) {
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.M {
		if f(v) {
			L = append(L, v)
		}
	}
	return
}

// Map returns a new slice containing the results of applying the function f to each string in the original slice.
func (c *SafeExpressionsMap) Map(f func(*model.Expression) *model.Expression) []*model.Expression {
	c.RLock()
	defer c.RUnlock()
	m := make([]*model.Expression, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// Has returns true if a element in the slice.
func (c *SafeExpressionsMap) Has(elem model.Expression) bool {
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.M {
		if v.ID == elem.ID {
			return true
		}
	}
	return false
}

// GetPage 返回一页缓存数据
func (c *SafeExpressionsMap) GetPage(offset, limit int) (Expressions []*model.Expression) {
	c.RLock()
	defer c.RUnlock()
	switch {
	case offset >= len(c.M) || offset < 0:
	case offset+int(limit) >= len(c.M):
		Expressions = c.M[offset:]
	default:
		Expressions = c.M[offset : offset+limit]
	}
	return
}

// Init 缓存初始化
func (c *SafeExpressionsMap) Init() {
	var m []*model.Expression

	if dt := db.Find(&m); dt.Error != nil {
		log.Errorf("[E] An error occurred while caching expressions, error: %s", dt.Error.Error())
		return
	}

	c.Lock()
	defer c.Unlock()

	c.M = m
}

package cache

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/api/app/model"
)

// SafeStrategiesMap 线程安全的数据缓存对象
type SafeStrategiesMap struct {
	sync.RWMutex
	M []*model.Strategy
}

// StrategiesMap 群集缓存对象
var StrategiesMap = &SafeStrategiesMap{}

// Count 返回缓存条数
func (c *SafeStrategiesMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Append 添加元素
func (c *SafeStrategiesMap) Append(item *model.Strategy) {
	c.RLock()
	defer c.RUnlock()
	c.M = append(c.M, item)
}

// Remove 删除元素，每次仅删除一个
func (c *SafeStrategiesMap) Remove(f func(*model.Strategy) bool) {
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
func (c *SafeStrategiesMap) Include(f func(*model.Strategy) bool) bool {
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
func (c *SafeStrategiesMap) Any(f func(*model.Strategy) bool) *model.Strategy {
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
func (c *SafeStrategiesMap) All() []*model.Strategy {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeStrategiesMap) Filter(f func(*model.Strategy) bool) (L []*model.Strategy) {
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
func (c *SafeStrategiesMap) Map(f func(*model.Strategy) *model.Strategy) []*model.Strategy {
	c.RLock()
	defer c.RUnlock()
	m := make([]*model.Strategy, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// Has returns true if a element in the slice.
func (c *SafeStrategiesMap) Has(elem model.Strategy) bool {
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
func (c *SafeStrategiesMap) GetPage(offset, limit int) (Strategies []*model.Strategy) {
	c.RLock()
	defer c.RUnlock()
	switch {
	case offset >= len(c.M) || offset < 0:
	case offset+int(limit) >= len(c.M):
		Strategies = c.M[offset:]
	default:
		Strategies = c.M[offset : offset+limit]
	}
	return
}

// Init 缓存初始化
func (c *SafeStrategiesMap) Init() {
	var m []*model.Strategy

	if dt := db.Find(&m); dt.Error != nil {
		log.Errorf("[E] An error occurred while caching strategies, error: %s", dt.Error.Error())
		return
	}

	c.Lock()
	defer c.Unlock()

	c.M = m
}

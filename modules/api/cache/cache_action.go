package cache

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/api/app/model"
)

// SafeActionsMap 线程安全的数据缓存对象
type SafeActionsMap struct {
	sync.RWMutex
	M []*model.Action
}

// ActionsMap 群集缓存对象
var ActionsMap = &SafeActionsMap{}

// Count 返回缓存条数
func (c *SafeActionsMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Append 添加元素
func (c *SafeActionsMap) Append(item *model.Action) {
	c.RLock()
	defer c.RUnlock()
	c.M = append(c.M, item)
}

// Remove 删除元素，每次仅删除一个
func (c *SafeActionsMap) Remove(f func(*model.Action) bool) {
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
func (c *SafeActionsMap) Include(f func(*model.Action) bool) bool {
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
func (c *SafeActionsMap) Any(f func(*model.Action) bool) *model.Action {
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
func (c *SafeActionsMap) All() []*model.Action {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeActionsMap) Filter(f func(*model.Action) bool) (L []*model.Action) {
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
func (c *SafeActionsMap) Map(f func(*model.Action) *model.Action) []*model.Action {
	c.RLock()
	defer c.RUnlock()
	m := make([]*model.Action, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// GetPage 返回一页缓存数据
func (c *SafeActionsMap) GetPage(offset, limit int) (Actions []*model.Action) {
	c.RLock()
	defer c.RUnlock()
	switch {
	case offset >= len(c.M) || offset < 0:
	case offset+int(limit) >= len(c.M):
		Actions = c.M[offset:]
	default:
		Actions = c.M[offset : offset+limit]
	}
	return
}

// Init 缓存初始化
func (c *SafeActionsMap) Init() {
	var m []*model.Action

	if dt := db.Find(&m); dt.Error != nil {
		log.Errorf("[E] An error occurred while caching actions, error: %s", dt.Error.Error())
		return
	}

	c.Lock()
	defer c.Unlock()

	c.M = m
}

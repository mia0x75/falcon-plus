package cache

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/api/app/model"
)

// SafeEdgesMap 线程安全的数据缓存对象
type SafeEdgesMap struct {
	sync.RWMutex
	M []*model.Edge
}

// EdgesMap 群集缓存对象
var EdgesMap = &SafeEdgesMap{}

// Count 返回缓存条数
func (c *SafeEdgesMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Append 添加元素
func (c *SafeEdgesMap) Append(item *model.Edge) {
	c.RLock()
	defer c.RUnlock()
	c.M = append(c.M, item)
}

// Remove 删除元素，每次仅删除一个
func (c *SafeEdgesMap) Remove(f func(*model.Edge) bool) {
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
func (c *SafeEdgesMap) Include(f func(*model.Edge) bool) bool {
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
func (c *SafeEdgesMap) Any(f func(*model.Edge) bool) *model.Edge {
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
func (c *SafeEdgesMap) All() []*model.Edge {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeEdgesMap) Filter(f func(*model.Edge) bool) (L []*model.Edge) {
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
func (c *SafeEdgesMap) Map(f func(*model.Edge) *model.Edge) []*model.Edge {
	c.RLock()
	defer c.RUnlock()
	m := make([]*model.Edge, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// Has returns true if a element in the slice.
func (c *SafeEdgesMap) Has(elem model.Edge) bool {
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.M {
		if v.AncestorID == elem.AncestorID && v.DescendantID == elem.DescendantID && v.Type == elem.Type {
			return true
		}
	}
	return false
}

// GetPage 返回一页缓存数据
func (c *SafeEdgesMap) GetPage(offset, limit int) (Edges []*model.Edge) {
	c.RLock()
	defer c.RUnlock()
	switch {
	case offset >= len(c.M) || offset < 0:
	case offset+int(limit) >= len(c.M):
		Edges = c.M[offset:]
	default:
		Edges = c.M[offset : offset+limit]
	}
	return
}

// Init 缓存初始化
func (c *SafeEdgesMap) Init() {
	var m []*model.Edge

	if dt := db.Find(&m); dt.Error != nil {
		log.Errorf("[E] An error occurred while caching edges, error: %s", dt.Error.Error())
		return
	}

	c.Lock()
	defer c.Unlock()

	c.M = m
}

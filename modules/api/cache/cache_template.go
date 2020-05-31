package cache

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/api/app/model"
)

// SafeTemplatesMap 线程安全的数据缓存对象
type SafeTemplatesMap struct {
	sync.RWMutex
	M []*model.Template
}

// TemplatesMap 群集缓存对象
var TemplatesMap = &SafeTemplatesMap{}

// Count 返回缓存条数
func (c *SafeTemplatesMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Append 添加元素
func (c *SafeTemplatesMap) Append(item *model.Template) {
	c.RLock()
	defer c.RUnlock()
	c.M = append(c.M, item)
}

// Remove 删除元素，每次仅删除一个
func (c *SafeTemplatesMap) Remove(f func(*model.Template) bool) {
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
func (c *SafeTemplatesMap) Include(f func(*model.Template) bool) bool {
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
func (c *SafeTemplatesMap) Any(f func(*model.Template) bool) *model.Template {
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
func (c *SafeTemplatesMap) All() []*model.Template {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeTemplatesMap) Filter(f func(*model.Template) bool) (L []*model.Template) {
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
func (c *SafeTemplatesMap) Map(f func(*model.Template) *model.Template) []*model.Template {
	c.RLock()
	defer c.RUnlock()
	m := make([]*model.Template, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// Has returns true if a element in the slice.
func (c *SafeTemplatesMap) Has(elem model.Template) bool {
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.M {
		if v.ID == elem.ID || v.Name == elem.Name {
			return true
		}
	}
	return false
}

// GetPage 返回一页缓存数据
func (c *SafeTemplatesMap) GetPage(offset, limit int) (Templates []*model.Template) {
	c.RLock()
	defer c.RUnlock()
	switch {
	case offset >= len(c.M) || offset < 0:
	case offset+int(limit) >= len(c.M):
		Templates = c.M[offset:]
	default:
		Templates = c.M[offset : offset+limit]
	}
	return
}

// Init 缓存初始化
func (c *SafeTemplatesMap) Init() {
	var m []*model.Template

	if dt := db.Find(&m); dt.Error != nil {
		log.Errorf("[E] An error occurred while caching templates, error: %s", dt.Error.Error())
		return
	}

	c.Lock()
	defer c.Unlock()

	c.M = m
}

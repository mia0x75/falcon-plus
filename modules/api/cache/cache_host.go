package cache

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/api/app/model"
)

// SafeHostsMap 线程安全的数据缓存对象
type SafeHostsMap struct {
	sync.RWMutex
	M []*model.Host
}

// HostsMap 群集缓存对象
var HostsMap = &SafeHostsMap{}

// Count 返回缓存条数
func (c *SafeHostsMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Append 添加元素
func (c *SafeHostsMap) Append(item *model.Host) {
	c.RLock()
	defer c.RUnlock()
	c.M = append(c.M, item)
}

// Remove 删除元素，每次仅删除一个
func (c *SafeHostsMap) Remove(f func(*model.Host) bool) {
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
func (c *SafeHostsMap) Include(f func(*model.Host) bool) bool {
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
func (c *SafeHostsMap) Any(f func(*model.Host) bool) *model.Host {
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
func (c *SafeHostsMap) All() []*model.Host {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeHostsMap) Filter(f func(*model.Host) bool) (L []*model.Host) {
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
func (c *SafeHostsMap) Map(f func(*model.Host) *model.Host) []*model.Host {
	c.RLock()
	defer c.RUnlock()
	m := make([]*model.Host, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// Has returns true if a element in the slice.
func (c *SafeHostsMap) Has(elem model.Host) bool {
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.M {
		if v.ID == elem.ID || v.Hostname == elem.Hostname {
			return true
		}
	}
	return false
}

// GetPage 返回一页缓存数据
func (c *SafeHostsMap) GetPage(offset, limit int) (Hosts []*model.Host) {
	c.RLock()
	defer c.RUnlock()
	switch {
	case offset >= len(c.M) || offset < 0:
	case offset+int(limit) >= len(c.M):
		Hosts = c.M[offset:]
	default:
		Hosts = c.M[offset : offset+limit]
	}
	return
}

// Init 缓存初始化
func (c *SafeHostsMap) Init() {
	var m []*model.Host

	if dt := db.Find(&m); dt.Error != nil {
		log.Errorf("[E] An error occurred while caching hosts, error: %s", dt.Error.Error())
		return
	}

	c.Lock()
	defer c.Unlock()

	c.M = m
}

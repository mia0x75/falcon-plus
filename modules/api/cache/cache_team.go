package cache

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/api/app/model"
)

// SafeTeamsMap 线程安全的数据缓存对象
type SafeTeamsMap struct {
	sync.RWMutex
	M []*model.Team
}

// TeamsMap 群集缓存对象
var TeamsMap = &SafeTeamsMap{}

// Count 返回缓存条数
func (c *SafeTeamsMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Append 添加元素
func (c *SafeTeamsMap) Append(item *model.Team) {
	c.RLock()
	defer c.RUnlock()
	c.M = append(c.M, item)
}

// Remove 删除元素，每次仅删除一个
func (c *SafeTeamsMap) Remove(f func(*model.Team) bool) {
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
func (c *SafeTeamsMap) Include(f func(*model.Team) bool) bool {
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
func (c *SafeTeamsMap) Any(f func(*model.Team) bool) *model.Team {
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
func (c *SafeTeamsMap) All() []*model.Team {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeTeamsMap) Filter(f func(*model.Team) bool) (L []*model.Team) {
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
func (c *SafeTeamsMap) Map(f func(*model.Team) *model.Team) []*model.Team {
	c.RLock()
	defer c.RUnlock()
	m := make([]*model.Team, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// Has returns true if a element in the slice.
func (c *SafeTeamsMap) Has(elem model.Team) bool {
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
func (c *SafeTeamsMap) GetPage(offset, limit int) (Teams []*model.Team) {
	c.RLock()
	defer c.RUnlock()
	switch {
	case offset >= len(c.M) || offset < 0:
	case offset+int(limit) >= len(c.M):
		Teams = c.M[offset:]
	default:
		Teams = c.M[offset : offset+limit]
	}
	return
}

// Init 缓存初始化
func (c *SafeTeamsMap) Init() {
	var m []*model.Team

	if dt := db.Find(&m); dt.Error != nil {
		log.Errorf("[E] An error occurred while caching teams, error: %s", dt.Error.Error())
		return
	}

	c.Lock()
	defer c.Unlock()

	c.M = m
}

package cache

import (
	"sync"

	"github.com/open-falcon/falcon-plus/modules/hbs/db"
)

// 一个机器可能在多个group下，做一个map缓存hostid与groupid的对应关系
type SafeHostGroupsMap struct {
	sync.RWMutex
	M map[int][]int
}

var HostGroupsMap = &SafeHostGroupsMap{M: make(map[int][]int)}

func (m *SafeHostGroupsMap) GetGroupIds(hid int) ([]int, bool) {
	m.RLock()
	defer m.RUnlock()
	gids, exists := m.M[hid]
	return gids, exists
}

func (m *SafeHostGroupsMap) Init() {
	groups, err := db.QueryHostGroups()
	if err != nil {
		return
	}

	m.Lock()
	defer m.Unlock()
	m.M = groups
}

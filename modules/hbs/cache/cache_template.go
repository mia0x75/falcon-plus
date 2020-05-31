package cache

import (
	"sync"

	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
)

// 一个HostGroup对应多个Template
type SafeGroupTemplates struct {
	sync.RWMutex
	M map[int][]int
}

var GroupTemplates = &SafeGroupTemplates{M: make(map[int][]int)}

func (m *SafeGroupTemplates) GetTemplateIds(gid int) ([]int, bool) {
	m.RLock()
	defer m.RUnlock()
	templateIds, exists := m.M[gid]
	return templateIds, exists
}

func (m *SafeGroupTemplates) Init() {
	templates, err := db.QueryGroupTemplates()
	if err != nil {
		return
	}

	m.Lock()
	defer m.Unlock()
	m.M = templates
}

type SafeTemplateCache struct {
	sync.RWMutex
	M map[int]*cm.Template
}

var TemplateCache = &SafeTemplateCache{M: make(map[int]*cm.Template)}

func (m *SafeTemplateCache) GetMap() map[int]*cm.Template {
	m.RLock()
	defer m.RUnlock()
	return m.M
}

func (m *SafeTemplateCache) Init() {
	templates, err := db.QueryTemplates()
	if err != nil {
		return
	}

	m.Lock()
	defer m.Unlock()
	m.M = templates
}

type SafeHostTemplateIDs struct {
	sync.RWMutex
	M map[int][]int
}

// TODO:
var HostTemplateIDs = &SafeHostTemplateIDs{M: make(map[int][]int)}

func (m *SafeHostTemplateIDs) GetMap() map[int][]int {
	m.RLock()
	defer m.RUnlock()
	return m.M
}

func (m *SafeHostTemplateIDs) Init() {
	templateIDs, err := db.QueryHostTemplateIDs()
	if err != nil {
		return
	}

	m.Lock()
	defer m.Unlock()
	m.M = templateIDs
}

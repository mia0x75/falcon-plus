package store

import (
	"container/list"
	"errors"
	"hash/crc32"
	"sync"

	log "github.com/sirupsen/logrus"

	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
)

var GraphItems *GraphItemMap

type GraphItemMap struct {
	sync.RWMutex
	A    []map[string]*SafeLinkedList
	Size int
}

func (m *GraphItemMap) Get(key string) (*SafeLinkedList, bool) {
	m.RLock()
	defer m.RUnlock()
	idx := hashKey(key) % uint32(m.Size)
	val, ok := m.A[idx][key]
	return val, ok
}

// Remove method remove key from GraphItemMap, return true if exists
func (m *GraphItemMap) Remove(key string) bool {
	m.Lock()
	defer m.Unlock()
	idx := hashKey(key) % uint32(m.Size)
	_, exists := m.A[idx][key]
	if !exists {
		return false
	}

	delete(m.A[idx], key)
	return true
}

func (m *GraphItemMap) Getitems(idx int) map[string]*SafeLinkedList {
	m.RLock()
	defer m.RUnlock()
	items := m.A[idx]
	m.A[idx] = make(map[string]*SafeLinkedList)
	return items
}

func (m *GraphItemMap) Set(key string, val *SafeLinkedList) {
	m.Lock()
	defer m.Unlock()
	idx := hashKey(key) % uint32(m.Size)
	m.A[idx][key] = val
}

func (m *GraphItemMap) Len() int {
	m.RLock()
	defer m.RUnlock()
	var l int
	for i := 0; i < m.Size; i++ {
		l += len(m.A[i])
	}
	return l
}

func (m *GraphItemMap) First(key string) *cm.GraphItem {
	m.RLock()
	defer m.RUnlock()
	idx := hashKey(key) % uint32(m.Size)
	sl, ok := m.A[idx][key]
	if !ok {
		return nil
	}

	first := sl.Front()
	if first == nil {
		return nil
	}

	return first.Value.(*cm.GraphItem)
}

func (m *GraphItemMap) PushAll(key string, items []*cm.GraphItem) error {
	m.Lock()
	defer m.Unlock()
	idx := hashKey(key) % uint32(m.Size)
	sl, ok := m.A[idx][key]
	if !ok {
		return errors.New("not exist")
	}
	sl.PushAll(items)
	return nil
}

func (m *GraphItemMap) GetFlag(key string) (uint32, error) {
	m.Lock()
	defer m.Unlock()
	idx := hashKey(key) % uint32(m.Size)
	sl, ok := m.A[idx][key]
	if !ok {
		return 0, errors.New("not exist")
	}
	return sl.Flag, nil
}

func (m *GraphItemMap) SetFlag(key string, flag uint32) error {
	m.Lock()
	defer m.Unlock()
	idx := hashKey(key) % uint32(m.Size)
	sl, ok := m.A[idx][key]
	if !ok {
		return errors.New("not exist")
	}
	sl.Flag = flag
	return nil
}

func (m *GraphItemMap) PopAll(key string) []*cm.GraphItem {
	m.Lock()
	defer m.Unlock()
	idx := hashKey(key) % uint32(m.Size)
	sl, ok := m.A[idx][key]
	if !ok {
		return []*cm.GraphItem{}
	}
	return sl.PopAll()
}

func (m *GraphItemMap) FetchAll(key string) ([]*cm.GraphItem, uint32) {
	m.RLock()
	defer m.RUnlock()
	idx := hashKey(key) % uint32(m.Size)
	sl, ok := m.A[idx][key]
	if !ok {
		return []*cm.GraphItem{}, 0
	}

	return sl.FetchAll()
}

func hashKey(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func getWts(key string, now int64) int64 {
	interval := int64(g.CACHE_TIME)
	return now + interval - (int64(hashKey(key)) % interval)
}

func (m *GraphItemMap) PushFront(key string,
	item *cm.GraphItem, md5 string, cfg *g.GlobalConfig) {
	if linkedList, exists := m.Get(key); exists {
		linkedList.PushFront(item)
	} else {
		safeList := &SafeLinkedList{L: list.New()}
		safeList.L.PushFront(item)

		if cfg.Migrate.Enabled && !g.IsRrdFileExist(g.RrdFileName(
			cfg.RRD.Storage, md5, item.DsType, item.Step)) {
			safeList.Flag = g.GRAPH_F_MISS
		}
		m.Set(key, safeList)
	}
}

func (m *GraphItemMap) KeysByIndex(idx int) []string {
	m.RLock()
	defer m.RUnlock()

	count := len(m.A[idx])
	if count == 0 {
		return []string{}
	}

	keys := make([]string, 0, count)
	for key := range m.A[idx] {
		keys = append(keys, key)
	}

	return keys
}

func (m *GraphItemMap) Back(key string) *cm.GraphItem {
	m.RLock()
	defer m.RUnlock()
	idx := hashKey(key) % uint32(m.Size)
	L, ok := m.A[idx][key]
	if !ok {
		return nil
	}

	back := L.Back()
	if back == nil {
		return nil
	}

	return back.Value.(*cm.GraphItem)
}

// 指定key对应的Item数量
func (m *GraphItemMap) ItemCnt(key string) int {
	m.RLock()
	defer m.RUnlock()
	idx := hashKey(key) % uint32(m.Size)
	L, ok := m.A[idx][key]
	if !ok {
		return 0
	}
	return L.Len()
}

func init() {
	size := g.CACHE_TIME / g.FLUSH_DISK_STEP
	if size < 0 {
		log.Panicf("[P] store.init, bad size %d", size)
	}

	GraphItems = &GraphItemMap{
		A:    make([]map[string]*SafeLinkedList, size),
		Size: size,
	}
	for i := 0; i < size; i++ {
		GraphItems.A[i] = make(map[string]*SafeLinkedList)
	}
}

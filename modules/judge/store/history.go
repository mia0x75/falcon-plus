package store

import (
	"container/list"
	"sync"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

type JudgeItemMap struct {
	sync.RWMutex
	M map[string]*SafeLinkedList
}

func NewJudgeItemMap() *JudgeItemMap {
	return &JudgeItemMap{M: make(map[string]*SafeLinkedList)}
}

func (m *JudgeItemMap) Get(key string) (*SafeLinkedList, bool) {
	m.RLock()
	defer m.RUnlock()
	val, ok := m.M[key]
	return val, ok
}

func (m *JudgeItemMap) Set(key string, val *SafeLinkedList) {
	m.Lock()
	defer m.Unlock()
	m.M[key] = val
}

func (m *JudgeItemMap) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.M)
}

func (m *JudgeItemMap) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.M, key)
}

func (m *JudgeItemMap) BatchDelete(keys []string) {
	count := len(keys)
	if count == 0 {
		return
	}

	m.Lock()
	defer m.Unlock()
	for i := 0; i < count; i++ {
		delete(m.M, keys[i])
	}
}

func (m *JudgeItemMap) CleanStale(before int64) {
	keys := []string{}

	m.RLock()
	for key, L := range m.M {
		front := L.Front()
		if front == nil {
			continue
		}

		if front.Value.(*cm.JudgeItem).Timestamp < before {
			keys = append(keys, key)
		}
	}
	m.RUnlock()

	m.BatchDelete(keys)
}

func (m *JudgeItemMap) PushFrontAndMaintain(key string, val *cm.JudgeItem, maxCount int, now int64) {
	if linkedList, exists := m.Get(key); exists {
		needJudge := linkedList.PushFrontAndMaintain(val, maxCount)
		if needJudge {
			Judge(linkedList, val, now)
		}
	} else {
		NL := list.New()
		NL.PushFront(val)
		safeList := &SafeLinkedList{L: NL}
		m.Set(key, safeList)
		Judge(safeList, val, now)
	}
}

// 这是个线程不安全的大Map，需要提前初始化好
var HistoryBigMap = make(map[string]*JudgeItemMap)

func InitHistoryBigMap() {
	arr := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			HistoryBigMap[arr[i]+arr[j]] = NewJudgeItemMap()
		}
	}
}

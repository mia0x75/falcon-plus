package store

import (
	"container/list"
	"sync"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

type SafeLinkedList struct {
	sync.RWMutex
	L *list.List
}

func (m *SafeLinkedList) ToSlice() []*cm.JudgeItem {
	m.RLock()
	defer m.RUnlock()
	sz := m.L.Len()
	if sz == 0 {
		return []*cm.JudgeItem{}
	}

	ret := make([]*cm.JudgeItem, 0, sz)
	for e := m.L.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*cm.JudgeItem))
	}
	return ret
}

// @param limit 至多返回这些，如果不够，有多少返回多少
// @return bool isEnough
func (m *SafeLinkedList) HistoryData(limit int) ([]*cm.HistoryData, bool) {
	if limit < 1 {
		// 其实limit不合法，此处也返回false吧，上层代码要注意
		// 因为false通常使上层代码进入异常分支，这样就统一了
		return []*cm.HistoryData{}, false
	}

	size := m.Len()
	if size == 0 {
		return []*cm.HistoryData{}, false
	}

	firstElement := m.Front()
	firstItem := firstElement.Value.(*cm.JudgeItem)

	var vs []*cm.HistoryData
	isEnough := true

	judgeType := firstItem.JudgeType[0]
	if judgeType == 'G' || judgeType == 'g' {
		if size < limit {
			// 有多少获取多少
			limit = size
			isEnough = false
		}
		vs = make([]*cm.HistoryData, limit)
		vs[0] = &cm.HistoryData{Timestamp: firstItem.Timestamp, Value: firstItem.Value}
		i := 1
		currentElement := firstElement
		for i < limit {
			nextElement := currentElement.Next()
			vs[i] = &cm.HistoryData{
				Timestamp: nextElement.Value.(*cm.JudgeItem).Timestamp,
				Value:     nextElement.Value.(*cm.JudgeItem).Value,
			}
			i++
			currentElement = nextElement
		}
	} else {
		if size < limit+1 {
			isEnough = false
			limit = size - 1
		}

		vs = make([]*cm.HistoryData, limit)

		i := 0
		currentElement := firstElement
		for i < limit {
			nextElement := currentElement.Next()
			diffVal := currentElement.Value.(*cm.JudgeItem).Value - nextElement.Value.(*cm.JudgeItem).Value
			diffTs := currentElement.Value.(*cm.JudgeItem).Timestamp - nextElement.Value.(*cm.JudgeItem).Timestamp
			vs[i] = &cm.HistoryData{
				Timestamp: currentElement.Value.(*cm.JudgeItem).Timestamp,
				Value:     diffVal / float64(diffTs),
			}
			i++
			currentElement = nextElement
		}
	}

	return vs, isEnough
}

func (m *SafeLinkedList) PushFront(v interface{}) *list.Element {
	m.Lock()
	defer m.Unlock()
	return m.L.PushFront(v)
}

// @return needJudge 如果是false不需要做judge，因为新上来的数据不合法
func (m *SafeLinkedList) PushFrontAndMaintain(v *cm.JudgeItem, maxCount int) bool {
	m.Lock()
	defer m.Unlock()

	sz := m.L.Len()
	if sz > 0 {
		// 新push上来的数据有可能重复了，或者timestamp不对，这种数据要丢掉
		if v.Timestamp <= m.L.Front().Value.(*cm.JudgeItem).Timestamp || v.Timestamp <= 0 {
			return false
		}
	}

	m.L.PushFront(v)

	sz++
	if sz <= maxCount {
		return true
	}

	del := sz - maxCount
	for i := 0; i < del; i++ {
		m.L.Remove(m.L.Back())
	}

	return true
}

func (m *SafeLinkedList) Front() *list.Element {
	m.RLock()
	defer m.RUnlock()
	return m.L.Front()
}

func (m *SafeLinkedList) Len() int {
	m.RLock()
	defer m.RUnlock()
	return m.L.Len()
}

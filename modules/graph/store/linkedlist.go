package store

import (
	"container/list"
	"sync"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

type SafeLinkedList struct {
	sync.RWMutex
	Flag uint32
	L    *list.List
}

// 新创建SafeLinkedList容器
func NewSafeLinkedList() *SafeLinkedList {
	return &SafeLinkedList{L: list.New()}
}

func (m *SafeLinkedList) PushFront(v interface{}) *list.Element {
	m.Lock()
	defer m.Unlock()
	return m.L.PushFront(v)
}

func (m *SafeLinkedList) Front() *list.Element {
	m.RLock()
	defer m.RUnlock()
	return m.L.Front()
}

func (m *SafeLinkedList) PopBack() *list.Element {
	m.Lock()
	defer m.Unlock()

	back := m.L.Back()
	if back != nil {
		m.L.Remove(back)
	}

	return back
}

func (m *SafeLinkedList) Back() *list.Element {
	m.Lock()
	defer m.Unlock()

	return m.L.Back()
}

func (m *SafeLinkedList) Len() int {
	m.RLock()
	defer m.RUnlock()
	return m.L.Len()
}

// remain参数表示要给linkedlist中留几个元素
// 在cron中刷磁盘的时候要留一个，用于创建数据库索引
// 在程序退出的时候要一个不留的全部刷到磁盘
func (m *SafeLinkedList) PopAll() []*cm.GraphItem {
	m.Lock()
	defer m.Unlock()

	size := m.L.Len()
	if size <= 0 {
		return []*cm.GraphItem{}
	}

	ret := make([]*cm.GraphItem, 0, size)

	for i := 0; i < size; i++ {
		item := m.L.Back()
		ret = append(ret, item.Value.(*cm.GraphItem))
		m.L.Remove(item)
	}

	return ret
}

// restore PushAll
func (m *SafeLinkedList) PushAll(items []*cm.GraphItem) {
	m.Lock()
	defer m.Unlock()

	size := len(items)
	if size > 0 {
		for i := size - 1; i >= 0; i-- {
			m.L.PushBack(items[i])
		}
	}
}

// return为倒叙的?
func (m *SafeLinkedList) FetchAll() ([]*cm.GraphItem, uint32) {
	m.Lock()
	defer m.Unlock()
	count := m.L.Len()
	ret := make([]*cm.GraphItem, 0, count)

	p := m.L.Back()
	for p != nil {
		ret = append(ret, p.Value.(*cm.GraphItem))
		p = p.Prev()
	}

	return ret, m.Flag
}

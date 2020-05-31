package sender

import (
	"container/list"
	"sync"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

type SafeLinkedList struct {
	sync.RWMutex
	L *list.List
}

func NewSafeLinkedList() *SafeLinkedList {
	return &SafeLinkedList{L: list.New()}
}

func (this *SafeLinkedList) PopBack(limit int) []*cm.JSONMetaData {
	this.RLock()
	defer this.RUnlock()
	sz := this.L.Len()
	if sz == 0 {
		return []*cm.JSONMetaData{}
	}

	if sz < limit {
		limit = sz
	}

	ret := make([]*cm.JSONMetaData, 0, limit)
	for i := 0; i < limit; i++ {
		e := this.L.Back()
		ret = append(ret, e.Value.(*cm.JSONMetaData))
		this.L.Remove(e)
	}

	return ret
}

func (this *SafeLinkedList) PushFront(v interface{}) *list.Element {
	this.Lock()
	defer this.Unlock()
	return this.L.PushFront(v)
}

func (this *SafeLinkedList) Front() *list.Element {
	this.RLock()
	defer this.RUnlock()
	return this.L.Front()
}

func (this *SafeLinkedList) Len() int {
	this.RLock()
	defer this.RUnlock()
	return this.L.Len()
}

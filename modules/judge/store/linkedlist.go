package store

import (
	"container/list"
	"sync"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

type SafeLinkedList struct {
	sync.RWMutex
	L *list.List
}

func (this *SafeLinkedList) ToSlice() []*cmodel.JudgeItem {
	this.RLock()
	defer this.RUnlock()
	sz := this.L.Len()
	if sz == 0 {
		return []*cmodel.JudgeItem{}
	}

	ret := make([]*cmodel.JudgeItem, 0, sz)
	for e := this.L.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*cmodel.JudgeItem))
	}
	return ret
}

// @param limit 至多返回这些，如果不够，有多少返回多少
// @return bool isEnough
func (this *SafeLinkedList) HistoryData(limit int) ([]*cmodel.HistoryData, bool) {
	if limit < 1 {
		// 其实limit不合法，此处也返回false吧，上层代码要注意
		// 因为false通常使上层代码进入异常分支，这样就统一了
		return []*cmodel.HistoryData{}, false
	}

	size := this.Len()
	if size == 0 {
		return []*cmodel.HistoryData{}, false
	}

	firstElement := this.Front()
	firstItem := firstElement.Value.(*cmodel.JudgeItem)

	var vs []*cmodel.HistoryData
	isEnough := true

	judgeType := firstItem.JudgeType[0]
	if judgeType == 'G' || judgeType == 'g' {
		if size < limit {
			// 有多少获取多少
			limit = size
			isEnough = false
		}
		vs = make([]*cmodel.HistoryData, limit)
		vs[0] = &cmodel.HistoryData{Timestamp: firstItem.Timestamp, Value: firstItem.Value}
		i := 1
		currentElement := firstElement
		for i < limit {
			nextElement := currentElement.Next()
			vs[i] = &cmodel.HistoryData{
				Timestamp: nextElement.Value.(*cmodel.JudgeItem).Timestamp,
				Value:     nextElement.Value.(*cmodel.JudgeItem).Value,
			}
			i++
			currentElement = nextElement
		}
	} else {
		if size < limit+1 {
			isEnough = false
			limit = size - 1
		}

		vs = make([]*cmodel.HistoryData, limit)

		i := 0
		currentElement := firstElement
		for i < limit {
			nextElement := currentElement.Next()
			diffVal := currentElement.Value.(*cmodel.JudgeItem).Value - nextElement.Value.(*cmodel.JudgeItem).Value
			diffTs := currentElement.Value.(*cmodel.JudgeItem).Timestamp - nextElement.Value.(*cmodel.JudgeItem).Timestamp
			vs[i] = &cmodel.HistoryData{
				Timestamp: currentElement.Value.(*cmodel.JudgeItem).Timestamp,
				Value:     diffVal / float64(diffTs),
			}
			i++
			currentElement = nextElement
		}
	}

	return vs, isEnough
}

func (this *SafeLinkedList) PushFront(v interface{}) *list.Element {
	this.Lock()
	defer this.Unlock()
	return this.L.PushFront(v)
}

// @return needJudge 如果是false不需要做judge，因为新上来的数据不合法
func (this *SafeLinkedList) PushFrontAndMaintain(v *cmodel.JudgeItem, maxCount int) bool {
	this.Lock()
	defer this.Unlock()

	sz := this.L.Len()
	if sz > 0 {
		// 新push上来的数据有可能重复了，或者timestamp不对，这种数据要丢掉
		if v.Timestamp <= this.L.Front().Value.(*cmodel.JudgeItem).Timestamp || v.Timestamp <= 0 {
			return false
		}
	}

	this.L.PushFront(v)

	sz++
	if sz <= maxCount {
		return true
	}

	del := sz - maxCount
	for i := 0; i < del; i++ {
		this.L.Remove(this.L.Back())
	}

	return true
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

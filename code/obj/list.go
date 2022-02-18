package obj

import (
	"container/list"
	"sync"
)

// ListValue list value
type ListValue struct {
	sync.RWMutex
	list.List
}

// NewString create string object
func NewList() *Obj {
	return &Obj{
		T:     ObjList,
		Value: new(ListValue),
	}
}

func (list *ListValue) appendLeft(next *ListValue) {
	if list == next {
		return
	}
	list.Lock()
	defer list.Unlock()
	next.RLock()
	defer next.RUnlock()
	for node := next.Back(); node != nil; node = node.Prev() {
		list.PushFront(node.Value)
	}
}

func (list *ListValue) appendRight(next *ListValue) {
	if list == next {
		return
	}
	list.Lock()
	defer list.Unlock()
	next.RLock()
	defer next.RUnlock()
	for node := next.Front(); node != nil; node = node.Next() {
		list.PushBack(node.Value)
	}
}

package obj

import "sync"

// Objs objects
type Objs struct {
	sync.RWMutex
	data map[string]*Obj
}

// New new objects
func New() *Objs {
	return &Objs{data: make(map[string]*Obj)}
}

// Set set object
func (objs *Objs) Set(key string, value *Obj) {
	objs.Lock()
	objs.data[key] = value
	objs.Unlock()
}

// MergeLeft merge list to left
func (objs *Objs) MergeLeft(key string, obj *Obj) {
	objs.Lock()
	defer objs.Unlock()
	if src, ok := objs.data[key]; !ok {
		objs.data[key] = obj
	} else {
		src.List().appendLeft(obj.List())
	}
}

func (objs *Objs) MergeRight(key string, obj *Obj) {
	objs.Lock()
	defer objs.Unlock()
	if src, ok := objs.data[key]; !ok {
		objs.data[key] = obj
	} else {
		src.List().appendRight(obj.List())
	}
}

// Lookup lookup object
func (objs *Objs) Lookup(key string) *Obj {
	objs.RLock()
	defer objs.RUnlock()
	return objs.data[key]
}

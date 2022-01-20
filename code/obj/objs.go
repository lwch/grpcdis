package obj

import "sync"

// Objs objects
type Objs struct {
	sync.RWMutex
	data map[string]*Obj // TODO: change to hash struct
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

// Lookup lookup object
func (objs *Objs) Lookup(key string) *Obj {
	objs.RLock()
	defer objs.RUnlock()
	return objs.data[key]
}

package obj

import "sync"

const blocks = 255

// Objs objects
type Objs struct {
	locks [blocks]sync.RWMutex
	data  [blocks]map[string]*Obj
}

// New new objects
func New() *Objs {
	ret := new(Objs)
	for i := 0; i < blocks; i++ {
		ret.data[i] = make(map[string]*Obj)
	}
	return ret
}

// Bernstein hash
func hash(str string) uint64 {
	hash := uint64(5381)
	for i := 0; i < len(str); i++ {
		hash = (hash << 5) + hash + uint64(str[i])
	}
	return hash
}

// Set set object
func (objs *Objs) Set(key string, value *Obj) {
	idx := uint8(key[0]) % blocks
	objs.locks[idx].Lock()
	objs.data[idx][key] = value
	objs.locks[idx].Unlock()
}

// Lookup lookup object
func (objs *Objs) Lookup(key string) *Obj {
	idx := uint8(key[0]) % blocks
	objs.locks[idx].RLock()
	defer objs.locks[idx].RUnlock()
	return objs.data[idx][key]
}

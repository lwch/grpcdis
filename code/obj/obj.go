package obj

import (
	"sync"
	"time"
)

// Type object type
type Type byte

const (
	// ObjString string object
	ObjString Type = iota
	// ObjList list object
	ObjList
	// ObjHash hash object
	ObjHash
	// ObjSet set object
	ObjSet
	// ObjSortedSet sorted set object
	ObjSortedSet
)

// Obj object
type Obj struct {
	sync.RWMutex
	T        Type
	Value    interface{}
	deadline time.Time
}

// Valid object is valid
func (obj *Obj) Valid() bool {
	return time.Now().Before(obj.deadline)
}

// String convert object to string value
func (obj *Obj) String() StringValue {
	obj.RLock()
	defer obj.RUnlock()
	return obj.Value.(StringValue)
}

// List convert object to list value
func (obj *Obj) List() *ListValue {
	obj.RLock()
	defer obj.RUnlock()
	return obj.Value.(*ListValue)
}

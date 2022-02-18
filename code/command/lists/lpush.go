package lists

import (
	"fmt"

	"github.com/lwch/goredis/code/command"
	"github.com/lwch/goredis/code/obj"
)

// Lpush lpush command
type Lpush struct {
	objs *obj.Objs
}

// NewLpush new lpush command
func NewLpush(objs *obj.Objs) *Lpush {
	return &Lpush{objs: objs}
}

// Name lpush command name
func (lpush *Lpush) Name() string {
	return "lpush"
}

// Argc lpush command argc
func (lpush *Lpush) Argc() int {
	return -3
}

// Flags lpush command flags
func (lpush *Lpush) Flags() []command.Flag {
	return []command.Flag{
		command.FlagWrite,
		command.FlagDenyOOM,
		command.FlagFast,
	}
}

// FirstKey lpush command first key index
func (lpush *Lpush) FirstKey() int {
	return 1
}

// LastKey lpush command last key index
func (lpush *Lpush) LastKey() int {
	return 1
}

// StepCount lpush command key step count
func (lpush *Lpush) StepCount() int {
	return 1
}

// Run run lpush command
func (lpush *Lpush) Run(argv [][]byte, w *command.PipeWriter) error {
	key := string(argv[0])
	o := lpush.objs.Lookup(key)
	if o == nil {
		o = obj.NewList()
	}
	list := o.List()
	list.Lock()
	for i := 1; i < len(argv); i++ {
		list.PushFront(string(argv[i]))
	}
	list.Unlock()
	lpush.objs.MergeLeft(key, o)
	w.Write([]byte(fmt.Sprintf(":%d\r\n", len(argv)-1)))
	return nil
}

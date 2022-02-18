package lists

import (
	"fmt"

	"github.com/lwch/goredis/code/command"
	"github.com/lwch/goredis/code/obj"
)

// Rpush rpush command
type Rpush struct {
	objs *obj.Objs
}

// NewRpush new rpush command
func NewRpush(objs *obj.Objs) *Rpush {
	return &Rpush{objs: objs}
}

// Name rpush command name
func (rpush *Rpush) Name() string {
	return "rpush"
}

// Argc rpush command argc
func (rpush *Rpush) Argc() int {
	return -3
}

// Flags rpush command flags
func (rpush *Rpush) Flags() []command.Flag {
	return []command.Flag{
		command.FlagWrite,
		command.FlagDenyOOM,
		command.FlagFast,
	}
}

// FirstKey rpush command first key index
func (rpush *Rpush) FirstKey() int {
	return 1
}

// LastKey rpush command last key index
func (rpush *Rpush) LastKey() int {
	return 1
}

// StepCount rpush command key step count
func (rpush *Rpush) StepCount() int {
	return 1
}

// Run run rpush command
func (rpush *Rpush) Run(argv [][]byte, w *command.PipeWriter) error {
	key := string(argv[0])
	o := rpush.objs.Lookup(key)
	if o == nil {
		o = obj.NewList()
	}
	list := o.List()
	list.Lock()
	for i := 1; i < len(argv); i++ {
		list.PushBack(string(argv[i]))
	}
	list.Unlock()
	rpush.objs.MergeRight(key, o)
	w.Write([]byte(fmt.Sprintf(":%d\r\n", len(argv)-1)))
	return nil
}

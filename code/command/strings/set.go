package strings

import (
	"github.com/lwch/goredis/code/command"
	"github.com/lwch/goredis/code/obj"
)

// Set set command
type Set struct {
	objs *obj.Objs
}

// NewSet new set command
func NewSet(objs *obj.Objs) *Set {
	return &Set{objs: objs}
}

// Name set command name
func (set *Set) Name() string {
	return "set"
}

// Argc set command argc
func (set *Set) Argc() int {
	return 3
}

// Flags set command flags
func (set *Set) Flags() []command.Flag {
	return []command.Flag{
		command.FlagWrite,
		command.FlagDenyOOM,
	}
}

// FirstKey set command first key index
func (set *Set) FirstKey() int {
	return 1
}

// LastKey set command last key index
func (set *Set) LastKey() int {
	return 1
}

// StepCount set command key step count
func (set *Set) StepCount() int {
	return 1
}

// Run run set command
func (set *Set) Run(argv [][]byte, w *command.LockWriter) error {
	key := string(argv[0])
	value := string(argv[1])
	set.objs.Set(key, obj.NewString(value))
	w.Lock()
	_, err := w.Write([]byte("+OK\r\n"))
	w.Unlock()
	return err
}

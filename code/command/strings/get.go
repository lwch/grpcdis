package strings

import (
	"fmt"
	"io"
	"strings"

	"github.com/lwch/goredis/code/command"
	"github.com/lwch/goredis/code/obj"
)

// Get get command
type Get struct {
	objs *obj.Objs
}

// NewGet new get command
func NewGet(objs *obj.Objs) *Get {
	return &Get{objs: objs}
}

// Name get command name
func (get *Get) Name() string {
	return "get"
}

// Argc get command argc
func (get *Get) Argc() int {
	return 2
}

// Flags get command flags
func (get *Get) Flags() []command.Flag {
	return []command.Flag{
		command.FlagReadOnly,
		command.FlagFast,
	}
}

// FirstKey get command first key index
func (get *Get) FirstKey() int {
	return 1
}

// LastKey get command last key index
func (get *Get) LastKey() int {
	return 1
}

// StepCount get command key step count
func (get *Get) StepCount() int {
	return 1
}

// Run run get command
func (get *Get) Run(argv [][]byte, w *command.LockWriter) error {
	key := string(argv[0])
	o := get.objs.Lookup(key)
	if o == nil {
		w.Lock()
		_, err := w.Write([]byte("$-1\r\n"))
		w.Unlock()
		return err
	}
	if o.T != obj.ObjString {
		w.Lock()
		_, err := fmt.Fprintf(w, "-WRONGTYPE wrong type of key %s\r\n", key)
		w.Unlock()
		return err
	}
	value := o.String()
	w.Lock()
	defer w.Unlock()
	_, err := fmt.Fprintf(w, "$%d\r\n", len(value))
	if err != nil {
		return err
	}
	_, err = io.Copy(w, strings.NewReader(string(value)))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\r\n"))
	return err
}

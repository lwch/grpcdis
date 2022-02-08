package strings

import (
	"bytes"
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
func (get *Get) Run(argv [][]byte, w *command.PipeWriter) error {
	key := string(argv[0])
	o := get.objs.Lookup(key)
	if o == nil {
		w.Write([]byte("$-1\r\n"))
		return nil
	}
	if o.T != obj.ObjString {
		w.Write([]byte(fmt.Sprintf("-WRONGTYPE wrong type of key %s\r\n", key)))
		return nil
	}
	value := o.String()
	var buf bytes.Buffer
	_, err := fmt.Fprintf(&buf, "$%d\r\n", len(value))
	if err != nil {
		return err
	}
	_, err = io.Copy(&buf, strings.NewReader(string(value)))
	if err != nil {
		return err
	}
	_, err = buf.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	w.Write(buf.Bytes())
	return nil
}

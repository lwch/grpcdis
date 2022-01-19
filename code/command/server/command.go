package server

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/lwch/goredis/code/command"
)

// Command command command
type Command struct {
	sync.RWMutex
	cmds map[string]command.Command
}

// NewCommand new command command
func NewCommand() *Command {
	cmd := &Command{
		cmds: make(map[string]command.Command),
	}
	cmd.Add(cmd)
	return cmd
}

// Name command command name
func (cmd *Command) Name() string {
	return "command"
}

// Argc command command argc
func (cmd *Command) Argc() int {
	return 1
}

// Flags command command flags
func (cmd *Command) Flags() []command.Flag {
	return []command.Flag{
		command.FlagLoading,
	}
}

// FirstKey command command first key index
func (cmd *Command) FirstKey() int {
	return 0
}

// LastKey command command last key index
func (cmd *Command) LastKey() int {
	return 0
}

// StepCount command command key step count
func (cmd *Command) StepCount() int {
	return 0
}

// Add add command
func (cmd *Command) Add(in command.Command) {
	cmd.Lock()
	cmd.cmds[in.Name()] = in
	cmd.Unlock()
}

// Run run command command
func (cmd *Command) Run(argv [][]byte, w *command.LockWriter) error {
	cmds := make([]command.Command, 0, len(cmd.cmds))
	cmd.RLock()
	for _, c := range cmd.cmds {
		cmds = append(cmds, c)
	}
	cmd.RUnlock()
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].Name() < cmds[j].Name()
	})
	w.Lock()
	defer w.Unlock()
	// commands count
	_, err := fmt.Fprintf(w, "*%d\r\n", len(cmds))
	if err != nil {
		return err
	}
	for _, c := range cmds {
		name := c.Name()
		flags := c.Flags()
		var strFlags string
		for _, flag := range flags {
			strFlags += "+" + flag.String() + "\r\n"
		}
		// name, argc, flags, first key, last key, step count
		_, err = fmt.Fprintf(w, "*%d\r\n$%d\r\n%s\r\n:%d\r\n*%d\r\n%s:%d\r\n:%d\r\n:%d\r\n",
			6, len(name), name, c.Argc(), len(flags), strFlags,
			c.FirstKey(), c.LastKey(), c.StepCount())
		if err != nil {
			return err
		}
	}
	return nil
}

// Lookup lookup command by name
func (cmd *Command) Lookup(name string) command.Command {
	cmd.RLock()
	defer cmd.RUnlock()
	return cmd.cmds[strings.ToLower(name)]
}

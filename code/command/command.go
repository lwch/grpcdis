package command

import (
	"io"
	"sync"
)

// LockWriter writer with lock
type LockWriter struct {
	sync.RWMutex
	io.Writer
}

// NewWriter new writer with lock
func NewWriter(w io.Writer) *LockWriter {
	return &LockWriter{Writer: w}
}

// Command command
type Command interface {
	Name() string
	Argc() int // negative is required arguments count
	Flags() []Flag
	FirstKey() int
	LastKey() int
	StepCount() int
	Run([][]byte, *LockWriter) error
}

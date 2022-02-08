package command

import (
	"io"

	"github.com/lwch/logging"
)

// PipeWriter writer with lock
type PipeWriter struct {
	io.Writer
	chWrite chan []byte
}

// NewWriter new writer with pipe
func NewWriter(w io.Writer) *PipeWriter {
	pw := &PipeWriter{
		Writer:  w,
		chWrite: make(chan []byte, 1024),
	}
	go pw.write()
	return pw
}

// Write write data
func (lw *PipeWriter) Write(data []byte) {
	lw.chWrite <- data
}

func (lw *PipeWriter) write() {
	for {
		data := <-lw.chWrite
		_, err := lw.Writer.Write(data)
		if err != nil {
			logging.Error("write: %v", err)
			return
		}
	}
}

// Command command
type Command interface {
	Name() string
	Argc() int // negative is required arguments count
	Flags() []Flag
	FirstKey() int
	LastKey() int
	StepCount() int
	Run([][]byte, *PipeWriter) error
}

package xpty

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var ErrUnsupported = errors.New("platform not supported by xpty")

func Open() (Master, Slave, error) {
	return open()
}

type Master interface {
	io.ReadWriteCloser
	GetSize() (int, int, error)
	SetSize(row, col int) error
}

type Slave interface {
	io.ReadWriteCloser
	Start(cmd *Cmd) (*os.Process, error)
}

type Cmd struct {
	Path string
	Args []string
}

type SizeError struct {
	Row, Col int
}

func (e *SizeError) Error() string {
	return fmt.Sprintf("attempt to set invalid terminal winsize (%d, %d)", e.Row, e.Col)
}

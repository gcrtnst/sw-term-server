package xpty

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var ErrUnsupported = errors.New("platform not supported by xpty")

func Open() (Terminal, error) {
	return open()
}

type Terminal interface {
	io.ReadWriteCloser
	Session(Size) (Session, error)
}

type Session interface {
	StartProcess(cmd Cmd) (*os.Process, error)
	GetSize() (Size, error)
	SetSize(Size) error
	Close() error
}

type Size struct {
	Row, Col int
}

type Cmd struct {
	Path string
	Args []string
}

type SizeError struct {
	Size Size
}

func (e *SizeError) Error() string {
	return fmt.Sprintf("attempt to set invalid terminal winsize (%d, %d)", e.Size.Row, e.Size.Col)
}

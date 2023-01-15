package xpty

import (
	"bytes"
	"errors"
	"io"
	"os"
)

var (
	ErrMockTerminalOpen    = errors.New("xpty: mock terminal open")
	ErrMockTerminalNotOpen = errors.New("xpty: mock terminal not open")
	ErrMockSessionOpen     = errors.New("xpty: mock session open")
	ErrMockSessionNotOpen  = errors.New("xpty: mock session not open")
)

type MockTerminal struct {
	ErrOpen          error
	ErrSession       error
	ErrStartProcess  error
	ErrGetSize       error
	ErrSetSize       error
	ErrCloseSession  error
	ErrCloseTerminal error
	PID              int

	Size         Size
	Cmd          Cmd
	OpenTerminal bool
	OpenSession  bool

	ir *io.PipeReader
	iw *io.PipeWriter
	ob *bytes.Buffer
}

func (t *MockTerminal) Open() (Terminal, error) {
	if t.ErrOpen != nil {
		return nil, t.ErrOpen
	}

	if t.OpenTerminal {
		panic(ErrMockTerminalOpen)
	}

	t.ir, t.iw = io.Pipe()
	t.ob = new(bytes.Buffer)
	t.OpenTerminal = true
	t.OpenSession = false
	return t, nil
}

func (t *MockTerminal) Computer() *MockComputer {
	return &MockComputer{t: t}
}

func (t *MockTerminal) Read(p []byte) (int, error) {
	return t.ir.Read(p)
}

func (t *MockTerminal) Write(p []byte) (int, error) {
	return t.ob.Write(p)
}

func (t *MockTerminal) Close() error {
	if t.ErrCloseTerminal != nil {
		return t.ErrCloseTerminal
	}

	if !t.OpenTerminal {
		panic(ErrMockTerminalNotOpen)
	}
	if t.OpenSession {
		panic(ErrMockSessionOpen)
	}

	_ = t.iw.CloseWithError(os.ErrClosed)
	t.OpenTerminal = false
	return nil
}

func (t *MockTerminal) Session(size Size) (Session, error) {
	if t.ErrSession != nil {
		return nil, t.ErrSession
	}

	if !t.OpenTerminal {
		panic(ErrMockTerminalNotOpen)
	}
	if t.OpenSession {
		panic(ErrMockSessionOpen)
	}

	t.Size = size
	t.OpenSession = true

	s := &MockSession{T: t}
	return s, nil
}

type MockSession struct {
	T *MockTerminal
}

func (s *MockSession) StartProcess(cmd Cmd) (*os.Process, error) {
	if s.T.ErrStartProcess != nil {
		return nil, s.T.ErrStartProcess
	}

	if !s.T.OpenSession {
		panic(ErrMockSessionNotOpen)
	}

	proc, err := os.FindProcess(s.T.PID)
	if err != nil {
		panic(err)
	}

	s.T.Cmd = cmd
	return proc, nil
}

func (s *MockSession) GetSize() (Size, error) {
	if s.T.ErrGetSize != nil {
		return Size{}, s.T.ErrGetSize
	}

	if !s.T.OpenSession {
		panic(ErrMockSessionNotOpen)
	}

	return s.T.Size, nil
}

func (s *MockSession) SetSize(size Size) error {
	if s.T.ErrSetSize != nil {
		return s.T.ErrSetSize
	}

	if !s.T.OpenSession {
		panic(ErrMockSessionNotOpen)
	}

	s.T.Size = size
	return nil
}

func (s *MockSession) Close() error {
	if s.T.ErrCloseSession != nil {
		return s.T.ErrCloseSession
	}

	if !s.T.OpenSession {
		panic(ErrMockSessionNotOpen)
	}
	s.T.OpenSession = false

	return nil
}

type MockComputer struct {
	t *MockTerminal
}

func (c *MockComputer) Read(p []byte) (int, error) {
	return c.t.ob.Read(p)
}

func (c *MockComputer) Write(p []byte) (int, error) {
	if !c.t.OpenTerminal {
		return 0, ErrMockTerminalNotOpen
	}
	return c.t.iw.Write(p)
}

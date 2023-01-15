package xpty

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestMockTerminalOpen(t *testing.T) {
	errDummy := errors.New("dummy error")

	tt := []struct {
		name             string
		t                *MockTerminal
		wantRecover      any
		wantPipe         bool
		wantOpenTerminal bool
		wantOpenSession  bool
		wantT            bool
		wantErr          error
	}{
		{
			name: "ErrOpen",
			t: &MockTerminal{
				ErrOpen:      errDummy,
				OpenTerminal: false,
				OpenSession:  false,
			},
			wantRecover:      nil,
			wantPipe:         false,
			wantOpenTerminal: false,
			wantOpenSession:  false,
			wantT:            false,
			wantErr:          errDummy,
		},
		{
			name: "AlreadyOpen",
			t: &MockTerminal{
				OpenTerminal: true,
				OpenSession:  false,
			},
			wantRecover:      ErrMockTerminalOpen,
			wantPipe:         false,
			wantOpenTerminal: true,
			wantOpenSession:  false,
		},
		{
			name: "NormalOpen",
			t: &MockTerminal{
				OpenTerminal: false,
				OpenSession:  true,
			},
			wantRecover:      nil,
			wantPipe:         true,
			wantOpenTerminal: true,
			wantOpenSession:  false,
			wantT:            true,
			wantErr:          nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var gotRecover any
			func() {
				defer func() {
					gotRecover = recover()
				}()

				var wantT Terminal
				if tc.wantT {
					wantT = tc.t
				}

				gotT, gotErr := tc.t.Open()
				if gotT != wantT {
					t.Errorf("t: expected %p, got %p", wantT, gotT)
				}
				if gotErr != tc.wantErr {
					t.Errorf("err: expected %#v, got %#v", tc.wantErr, gotErr)
				}
			}()

			gotPipe := tc.t.ir != nil && tc.t.iw != nil && tc.t.ob != nil

			if gotRecover != tc.wantRecover {
				panic(gotRecover)
			}
			if gotPipe != tc.wantPipe {
				t.Errorf("t pipe: expected %t, got %t", tc.wantPipe, gotPipe)
			}
			if tc.t.OpenTerminal != tc.wantOpenTerminal {
				t.Errorf("t.openTerminal: expected %t, got %t", tc.wantOpenTerminal, tc.t.OpenTerminal)
			}
			if tc.t.OpenSession != tc.wantOpenSession {
				t.Errorf("t.openSession: expected %t, got %t", tc.wantOpenSession, tc.t.OpenSession)
			}
		})
	}
}

func TestMockTerminalClose(t *testing.T) {
	var ir *io.PipeReader
	var iw *io.PipeWriter
	errDummy := errors.New("dummy error")

	type testCase struct {
		name             string
		t                *MockTerminal
		wantRecover      any
		wantOpenTerminal bool
		wantCloseInput   bool
		wantErr          error
	}

	tt := []testCase{}

	ir, iw = io.Pipe()
	tt = append(tt, testCase{
		name: "ErrCloseTerminal",
		t: &MockTerminal{
			ErrCloseTerminal: errDummy,
			ir:               ir,
			iw:               iw,
			ob:               new(bytes.Buffer),
			OpenTerminal:     true,
			OpenSession:      false,
		},
		wantRecover:      nil,
		wantOpenTerminal: true,
		wantCloseInput:   false,
		wantErr:          errDummy,
	})

	ir, iw = io.Pipe()
	tt = append(tt, testCase{
		name: "TerminalNotOpen",
		t: &MockTerminal{
			ErrCloseTerminal: nil,
			ir:               ir,
			iw:               iw,
			ob:               new(bytes.Buffer),
			OpenTerminal:     false,
			OpenSession:      false,
		},
		wantRecover:      ErrMockTerminalNotOpen,
		wantOpenTerminal: false,
	})

	ir, iw = io.Pipe()
	tt = append(tt, testCase{
		name: "SessionOpen",
		t: &MockTerminal{
			ErrCloseTerminal: nil,
			ir:               ir,
			iw:               iw,
			ob:               new(bytes.Buffer),
			OpenTerminal:     true,
			OpenSession:      true,
		},
		wantRecover:      ErrMockSessionOpen,
		wantOpenTerminal: true,
	})

	ir, iw = io.Pipe()
	tt = append(tt, testCase{
		name: "Normal",
		t: &MockTerminal{
			ErrCloseTerminal: nil,
			ir:               ir,
			iw:               iw,
			ob:               new(bytes.Buffer),
			OpenTerminal:     true,
			OpenSession:      false,
		},
		wantRecover:      nil,
		wantOpenTerminal: false,
		wantCloseInput:   true,
		wantErr:          nil,
	})

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var gotRecover any
			func() {
				defer func() {
					gotRecover = recover()
				}()

				gotErr := tc.t.Close()

				done := make(chan struct{})
				go func() {
					_, _ = tc.t.iw.Write([]byte{})
					close(done)
				}()

				var err error
				_, err = tc.t.Read([]byte{})
				gotCloseInput := err == os.ErrClosed
				<-done

				if gotCloseInput != tc.wantCloseInput {
					t.Errorf("t.iw close: expected %t, got %t", tc.wantCloseInput, gotCloseInput)
				}
				if gotErr != tc.wantErr {
					t.Errorf("err: expected %#v, got %#v", tc.wantErr, gotErr)
				}
			}()

			if gotRecover != tc.wantRecover {
				panic(gotRecover)
			}
			if tc.t.OpenTerminal != tc.wantOpenTerminal {
				t.Errorf("t.openTerminal: expected %t, got %t", tc.wantOpenTerminal, tc.t.OpenTerminal)
			}
		})
	}
}

func TestMockTerminalSession(t *testing.T) {
	errDummy := errors.New("dummy error")

	tt := []struct {
		name            string
		t               *MockTerminal
		inSize          Size
		wantRecover     any
		wantSize        Size
		wantOpenSession bool
		wantS           bool
		wantErr         error
	}{
		{
			name: "ErrSession",
			t: &MockTerminal{
				ErrSession:   errDummy,
				Size:         Size{},
				OpenTerminal: true,
				OpenSession:  false,
			},
			inSize:          Size{Row: 30, Col: 120},
			wantRecover:     nil,
			wantSize:        Size{},
			wantOpenSession: false,
			wantS:           false,
			wantErr:         errDummy,
		},
		{
			name: "TerminalNotOpen",
			t: &MockTerminal{
				ErrSession:   nil,
				Size:         Size{},
				OpenTerminal: false,
				OpenSession:  false,
			},
			inSize:          Size{Row: 30, Col: 120},
			wantRecover:     ErrMockTerminalNotOpen,
			wantSize:        Size{},
			wantOpenSession: false,
		},
		{
			name: "SessionOpen",
			t: &MockTerminal{
				ErrSession:   nil,
				Size:         Size{},
				OpenTerminal: true,
				OpenSession:  true,
			},
			inSize:          Size{Row: 30, Col: 120},
			wantRecover:     ErrMockSessionOpen,
			wantSize:        Size{},
			wantOpenSession: true,
		},
		{
			name: "Normal",
			t: &MockTerminal{
				ErrSession:   nil,
				Size:         Size{},
				OpenTerminal: true,
				OpenSession:  false,
			},
			inSize:          Size{Row: 30, Col: 120},
			wantRecover:     nil,
			wantSize:        Size{Row: 30, Col: 120},
			wantOpenSession: true,
			wantS:           true,
			wantErr:         nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var gotRecover any
			func() {
				defer func() {
					gotRecover = recover()
				}()

				var wantS Session
				if tc.wantS {
					wantS = &MockSession{T: tc.t}
				}

				gotS, gotErr := tc.t.Session(tc.inSize)

				if !reflect.DeepEqual(gotS, wantS) {
					t.Errorf("s: expected %#v, got %#v", wantS, gotS)
				}
				if gotErr != tc.wantErr {
					t.Errorf("err: expected %#v, got %#v", tc.wantErr, gotErr)
				}
			}()

			if gotRecover != tc.wantRecover {
				panic(gotRecover)
			}
			if tc.t.Size != tc.wantSize {
				t.Errorf("t.Size: expected %#v, got %#v", tc.wantSize, tc.t.Size)
			}
			if tc.t.OpenSession != tc.wantOpenSession {
				t.Errorf("t.openSession: expected %t, got %t", tc.wantOpenSession, tc.t.OpenSession)
			}
		})
	}
}

func TestMockSessionStartProcess(t *testing.T) {
	errDummy := errors.New("dummy error")
	pid := os.Getpid()

	tt := []struct {
		name        string
		t           *MockTerminal
		inCmd       Cmd
		wantRecover any
		wantCmd     Cmd
		wantPID     int
		wantErr     error
	}{
		{
			name: "ErrStartProcess",
			t: &MockTerminal{
				ErrStartProcess: errDummy,
				PID:             pid,
				Cmd:             Cmd{},
				OpenTerminal:    true,
				OpenSession:     true,
			},
			inCmd: Cmd{
				Path: "bash",
				Args: []string{"--version"},
			},
			wantRecover: nil,
			wantCmd:     Cmd{},
			wantPID:     0,
			wantErr:     errDummy,
		},
		{
			name: "SessionNotOpen",
			t: &MockTerminal{
				ErrStartProcess: nil,
				PID:             pid,
				Cmd:             Cmd{},
				OpenTerminal:    true,
				OpenSession:     false,
			},
			inCmd: Cmd{
				Path: "bash",
				Args: []string{"--version"},
			},
			wantRecover: ErrMockSessionNotOpen,
			wantCmd:     Cmd{},
		},
		{
			name: "Normal",
			t: &MockTerminal{
				ErrStartProcess: nil,
				PID:             pid,
				Cmd:             Cmd{},
				OpenTerminal:    true,
				OpenSession:     true,
			},
			inCmd: Cmd{
				Path: "bash",
				Args: []string{"--version"},
			},
			wantRecover: nil,
			wantCmd: Cmd{
				Path: "bash",
				Args: []string{"--version"},
			},
			wantPID: pid,
			wantErr: nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var gotRecover any
			func() {
				defer func() {
					gotRecover = recover()
				}()

				s := &MockSession{T: tc.t}
				gotProc, gotErr := s.StartProcess(tc.inCmd)

				gotPID := 0
				if gotProc != nil {
					gotPID = gotProc.Pid
				}

				if gotPID != tc.wantPID {
					t.Errorf("proc.Pid: expected %d, got %d", tc.wantPID, gotPID)
				}
				if gotErr != tc.wantErr {
					t.Errorf("err: expected %#v, got %#v", tc.wantErr, gotErr)
				}
			}()

			if gotRecover != tc.wantRecover {
				panic(gotRecover)
			}
			if !reflect.DeepEqual(tc.t.Cmd, tc.wantCmd) {
				t.Errorf("t.Cmd: expected %#v, got %#v", tc.wantCmd, tc.t.Cmd)
			}
		})
	}
}

func TestMockSessionGetSize(t *testing.T) {
	errDummy := errors.New("dummy error")

	tt := []struct {
		name        string
		t           *MockTerminal
		wantRecover any
		wantSize    Size
		wantErr     error
	}{
		{
			name: "ErrGetSize",
			t: &MockTerminal{
				ErrGetSize:   errDummy,
				Size:         Size{Row: 30, Col: 120},
				OpenTerminal: true,
				OpenSession:  true,
			},
			wantRecover: nil,
			wantSize:    Size{},
			wantErr:     errDummy,
		},
		{
			name: "SessionNotOpen",
			t: &MockTerminal{
				ErrGetSize:   nil,
				Size:         Size{Row: 30, Col: 120},
				OpenTerminal: true,
				OpenSession:  false,
			},
			wantRecover: ErrMockSessionNotOpen,
		},
		{
			name: "Normal",
			t: &MockTerminal{
				ErrGetSize:   nil,
				Size:         Size{Row: 30, Col: 120},
				OpenTerminal: true,
				OpenSession:  true,
			},
			wantRecover: nil,
			wantSize:    Size{Row: 30, Col: 120},
			wantErr:     nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var gotRecover any
			func() {
				defer func() {
					gotRecover = recover()
				}()

				s := &MockSession{T: tc.t}
				gotSize, gotErr := s.GetSize()

				if gotSize != tc.wantSize {
					t.Errorf("size: expected %#v, got %#v", tc.wantSize, gotSize)
				}
				if gotErr != tc.wantErr {
					t.Errorf("err: expected %#v, got %#v", tc.wantErr, gotErr)
				}
			}()

			if gotRecover != tc.wantRecover {
				panic(gotRecover)
			}
		})
	}
}

func TestMockSessionSetSize(t *testing.T) {
	errDummy := errors.New("dummy error")

	tt := []struct {
		name        string
		t           *MockTerminal
		inSize      Size
		wantRecover any
		wantSize    Size
		wantErr     error
	}{
		{
			name: "ErrSetSize",
			t: &MockTerminal{
				ErrSetSize:   errDummy,
				Size:         Size{},
				OpenTerminal: true,
				OpenSession:  true,
			},
			inSize:      Size{Row: 30, Col: 120},
			wantRecover: nil,
			wantSize:    Size{},
			wantErr:     errDummy,
		},
		{
			name: "SessionNotOpen",
			t: &MockTerminal{
				ErrSetSize:   nil,
				Size:         Size{},
				OpenTerminal: true,
				OpenSession:  false,
			},
			inSize:      Size{Row: 30, Col: 120},
			wantRecover: ErrMockSessionNotOpen,
		},
		{
			name: "Normal",
			t: &MockTerminal{
				ErrSetSize:   nil,
				Size:         Size{},
				OpenTerminal: true,
				OpenSession:  true,
			},
			inSize:      Size{Row: 30, Col: 120},
			wantRecover: nil,
			wantSize:    Size{Row: 30, Col: 120},
			wantErr:     nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var gotRecover any
			func() {
				defer func() {
					gotRecover = recover()
				}()

				s := &MockSession{T: tc.t}
				gotErr := s.SetSize(tc.inSize)

				if gotErr != tc.wantErr {
					t.Errorf("err: expected %#v, got %#v", tc.wantErr, gotErr)
				}
			}()

			if gotRecover != tc.wantRecover {
				panic(gotRecover)
			}
			if tc.t.Size != tc.wantSize {
				t.Errorf("t.Size: expected %#v, got %#v", tc.wantSize, tc.t.Size)
			}
		})
	}
}

func TestMockSessionClose(t *testing.T) {
	errDummy := errors.New("dummy error")

	tt := []struct {
		name            string
		t               *MockTerminal
		wantRecover     any
		wantOpenSession bool
		wantErr         error
	}{
		{
			name: "ErrCloseSession",
			t: &MockTerminal{
				ErrCloseSession: errDummy,
				OpenTerminal:    true,
				OpenSession:     true,
			},
			wantRecover:     nil,
			wantOpenSession: true,
			wantErr:         errDummy,
		},
		{
			name: "SessionNotOpen",
			t: &MockTerminal{
				ErrCloseSession: nil,
				OpenTerminal:    true,
				OpenSession:     false,
			},
			wantRecover:     ErrMockSessionNotOpen,
			wantOpenSession: false,
		},
		{
			name: "Normal",
			t: &MockTerminal{
				ErrCloseSession: nil,
				OpenTerminal:    true,
				OpenSession:     true,
			},
			wantRecover:     nil,
			wantOpenSession: false,
			wantErr:         nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var gotRecover any
			func() {
				defer func() {
					gotRecover = recover()
				}()

				s := &MockSession{T: tc.t}
				gotErr := s.Close()

				if gotErr != tc.wantErr {
					t.Errorf("err: expected %#v, got %#v", tc.wantErr, gotErr)
				}
			}()

			if gotRecover != tc.wantRecover {
				panic(gotRecover)
			}
			if tc.t.OpenSession != tc.wantOpenSession {
				t.Errorf("t.openSession: expected %t, got %t", tc.wantOpenSession, tc.t.OpenSession)
			}
		})
	}
}

func TestMockComputerWrite(t *testing.T) {
	in := []byte("foo")

	mt := &MockTerminal{}
	_, _ = mt.Open()

	mc := mt.Computer()
	done := make(chan struct{})
	go func() {
		_, _ = mc.Write(in)
		close(done)
	}()

	p := make([]byte, len(in))
	n, err := mt.Read(p)
	if n != len(in) {
		t.Errorf("n: expected %d, got %d", len(in), n)
	}
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if !bytes.Equal(p, in) {
		t.Errorf("p: expected %#v, got %#v", in, p)
	}
	<-done
}

func TestMockComputerWriteNotOpen(t *testing.T) {
	mt := &MockTerminal{}
	mc := mt.Computer()
	_, err := mc.Write(nil)
	if err != ErrMockTerminalNotOpen {
		t.Errorf("err: expected %#v, got %#v", ErrMockTerminalNotOpen, err)
	}
}

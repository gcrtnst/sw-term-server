package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/gcrtnst/sw-term-server/internal/vterm"
	"github.com/gcrtnst/sw-term-server/internal/xpty"
)

func TestNewTerm(t *testing.T) {
	errDummy := errors.New("dummy error")
	pid := os.Getpid()

	tt := []struct {
		name               string
		pt                 *xpty.MockTerminal
		inCfg              TermConfig
		wantPTSize         xpty.Size
		wantPTCmd          xpty.Cmd
		wantPTOpenTerminal bool
		wantPTOpenSession  bool
		wantErr            error
		wantVTRows         int
		wantVTCols         int
	}{
		{
			name: "Normal",
			pt: &xpty.MockTerminal{
				PID: pid,
			},
			inCfg: TermConfig{
				Row: 30,
				Col: 120,
				Cmd: xpty.Cmd{
					Path: "bash",
					Args: []string{"--version"},
				},
			},
			wantPTSize: xpty.Size{
				Row: 30,
				Col: 120,
			},
			wantPTCmd: xpty.Cmd{
				Path: "bash",
				Args: []string{"--version"},
			},
			wantPTOpenTerminal: true,
			wantPTOpenSession:  true,
			wantErr:            nil,
			wantVTRows:         30,
			wantVTCols:         120,
		},
		{
			name: "ErrOpen",
			pt: &xpty.MockTerminal{
				ErrOpen: errDummy,
				PID:     pid,
			},
			inCfg: TermConfig{
				Row: 30,
				Col: 120,
				Cmd: xpty.Cmd{
					Path: "bash",
					Args: []string{"--version"},
				},
			},
			wantPTSize:         xpty.Size{},
			wantPTCmd:          xpty.Cmd{},
			wantPTOpenTerminal: false,
			wantPTOpenSession:  false,
			wantErr:            errDummy,
		},
		{
			name: "ErrSession",
			pt: &xpty.MockTerminal{
				ErrSession: errDummy,
				PID:        pid,
			},
			inCfg: TermConfig{
				Row: 30,
				Col: 120,
				Cmd: xpty.Cmd{
					Path: "bash",
					Args: []string{"--version"},
				},
			},
			wantPTSize:         xpty.Size{},
			wantPTCmd:          xpty.Cmd{},
			wantPTOpenTerminal: false,
			wantPTOpenSession:  false,
			wantErr:            errDummy,
		},
		{
			name: "ErrStartProcess",
			pt: &xpty.MockTerminal{
				ErrStartProcess: errDummy,
				PID:             pid,
			},
			inCfg: TermConfig{
				Row: 30,
				Col: 120,
				Cmd: xpty.Cmd{
					Path: "bash",
					Args: []string{"--version"},
				},
			},
			wantPTSize: xpty.Size{
				Row: 30,
				Col: 120,
			},
			wantPTCmd:          xpty.Cmd{},
			wantPTOpenTerminal: false,
			wantPTOpenSession:  false,
			wantErr:            errDummy,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.inCfg.Open = tc.pt.Open

			gotTerm, gotErr := NewTerm(tc.inCfg)
			if gotErr != tc.wantErr {
				t.Errorf("err: expected %#v, got %#v", tc.wantErr, gotErr)
			}

			gotTermNonNil := gotTerm != nil
			wantTermNonNil := tc.wantErr == nil
			if gotTermNonNil != wantTermNonNil {
				t.Errorf("term non-nil: expected %t, got %t", wantTermNonNil, gotTermNonNil)
			}

			if gotTerm != nil {
				if gotTerm.pt != tc.pt {
					t.Errorf("term.pt: expected %#v, got %#v", tc.pt, gotTerm.pt)
				}
				if ms, ok := gotTerm.ps.(*xpty.MockSession); !ok || ms.T != tc.pt {
					wantPS := &xpty.MockSession{T: tc.pt}
					t.Errorf("term.ps: expected %#v, got %#v", wantPS, gotTerm.ps)
				}

				gotVTRows, gotVTCols := gotTerm.vt.GetSize()
				if gotVTRows != tc.wantVTRows || gotVTCols != tc.wantVTCols {
					t.Errorf("term.vt size: expected (%d, %d), got (%d, %d)", tc.wantVTRows, tc.wantVTCols, gotVTRows, gotVTCols)
				}

				gotVTUTF8 := gotTerm.vt.GetUTF8()
				const wantVTUTF8 = true
				if gotVTUTF8 != wantVTUTF8 {
					t.Errorf("term.vt utf8: expected %t, got %t", wantVTUTF8, gotVTUTF8)
				}

				select {
				case <-gotTerm.di:
					t.Errorf("term.di: closed")
				default:
				}

				select {
				case <-gotTerm.do:
					t.Errorf("term.do: closed")
				default:
				}
			}

			if tc.pt.Size != tc.wantPTSize {
				t.Errorf("pt size: expected %#v, got %#v", tc.wantPTSize, tc.pt.Size)
			}
			if !reflect.DeepEqual(tc.pt.Cmd, tc.wantPTCmd) {
				t.Errorf("pt cmd: expected %#v, got %#v", tc.wantPTCmd, tc.pt.Cmd)
			}
			if tc.pt.OpenTerminal != tc.wantPTOpenTerminal {
				t.Errorf("pt open terminal: expected %t, got %t", tc.wantPTOpenTerminal, tc.pt.OpenTerminal)
			}
			if tc.pt.OpenSession != tc.wantPTOpenSession {
				t.Errorf("pt open terminal: expected %t, got %t", tc.wantPTOpenSession, tc.pt.OpenSession)
			}
		})
	}
}

func TestTermKeyboard(t *testing.T) {
	pid := os.Getpid()

	tt := []struct {
		name    string
		inKey   Key
		inMod   vterm.Modifier
		wantOK  bool
		wantOut []byte
	}{
		{
			name:    "NormalA",
			inKey:   "A",
			inMod:   vterm.ModNone,
			wantOK:  true,
			wantOut: []byte("A"),
		},
		{
			name:    "ModA",
			inKey:   "A",
			inMod:   vterm.ModShift | vterm.ModAlt | vterm.ModCtrl,
			wantOK:  true,
			wantOut: []byte("\x1B[65;7u"),
		},
		{
			name:    "NormalEnter",
			inKey:   "Enter",
			inMod:   vterm.ModNone,
			wantOK:  true,
			wantOut: []byte("\r"),
		},
		{
			name:    "ModEnter",
			inKey:   "Enter",
			inMod:   vterm.ModShift | vterm.ModAlt | vterm.ModCtrl,
			wantOK:  true,
			wantOut: []byte("\x1B[13;8u"),
		},
		{
			name:    "Invalid",
			inKey:   "",
			inMod:   vterm.ModNone,
			wantOK:  false,
			wantOut: []byte(""),
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mt := &xpty.MockTerminal{PID: pid}
			cfg := TermConfig{
				Open: mt.Open,
				Row:  30,
				Col:  120,
				Cmd: xpty.Cmd{
					Path: "bash",
					Args: []string{"--version"},
				},
			}

			term, err := NewTerm(cfg)
			if err != nil {
				t.Errorf("new: %s", err.Error())
			}
			term.pc = nil

			gotOK := term.Keyboard(tc.inKey, tc.inMod)
			err = term.Close()
			if err != nil {
				panic(err)
			}

			mc := mt.Computer()
			gotOut, _ := io.ReadAll(mc)

			if gotOK != tc.wantOK {
				t.Errorf("ok: expected %t, got %t", tc.wantOK, gotOK)
			}
			if !bytes.Equal(gotOut, tc.wantOut) {
				t.Errorf("out: expected %#v, got %#v", tc.wantOut, gotOut)
			}
		})
	}
}

func TestTermClose(t *testing.T) {
	pt := &xpty.MockTerminal{
		PID: os.Getpid(),
	}
	cfg := TermConfig{
		Open: pt.Open,
		Row:  30,
		Col:  120,
		Cmd: xpty.Cmd{
			Path: "bash",
			Args: []string{"--version"},
		},
	}

	term, err := NewTerm(cfg)
	if err != nil {
		t.Errorf("new: %s", err.Error())
	}
	term.pc = nil

	err = term.Close()
	if err != nil {
		t.Errorf("close 1: %s", err.Error())
	}

	err = term.Close()
	if err != nil {
		t.Errorf("close 2: %s", err.Error())
	}
}

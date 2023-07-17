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

func TestTermSlotKeyboard(t *testing.T) {
	errDummy := errors.New("dummy error")
	pid := os.Getpid()

	tt := []struct {
		name               string
		inStart            bool
		inMTErrOpen        error
		inKey              Key
		inMod              vterm.Modifier
		wantErr            error
		wantMTOpenTerminal bool
		wantOut            []byte
	}{
		{
			name:               "Normal",
			inStart:            true,
			inMTErrOpen:        errDummy,
			inKey:              "A",
			inMod:              vterm.ModAlt | vterm.ModCtrl,
			wantErr:            nil,
			wantMTOpenTerminal: true,
			wantOut:            []byte("\x1B[65;7u"),
		},
		{
			name:               "Start",
			inStart:            false,
			inMTErrOpen:        nil,
			inKey:              "A",
			inMod:              vterm.ModAlt | vterm.ModCtrl,
			wantErr:            nil,
			wantMTOpenTerminal: true,
			wantOut:            []byte("\x1B[65;7u"),
		},
		{
			name:               "StartError",
			inStart:            false,
			inMTErrOpen:        errDummy,
			inKey:              "A",
			inMod:              vterm.ModAlt | vterm.ModCtrl,
			wantErr:            errDummy,
			wantMTOpenTerminal: false,
			wantOut:            []byte{},
		},
		{
			name:               "Invalid",
			inStart:            true,
			inMTErrOpen:        errDummy,
			inKey:              "",
			inMod:              vterm.ModAlt | vterm.ModCtrl,
			wantErr:            ErrInvalidKey,
			wantMTOpenTerminal: true,
			wantOut:            []byte{},
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
			slot := NewTermSlot(cfg)

			if tc.inStart {
				err := slot.start()
				if err != nil {
					t.Fatal(err)
				}
			}

			mt.ErrOpen = tc.inMTErrOpen
			gotErr := slot.Keyboard(tc.inKey, tc.inMod)
			gotMTOpenTerminal := mt.OpenTerminal

			if slot.term != nil {
				slot.term.pc = nil
			}
			slot.Stop()

			gotOut := []byte{}
			if gotMTOpenTerminal {
				mc := mt.Computer()
				gotOut, _ = io.ReadAll(mc)
			}

			if gotErr != tc.wantErr {
				t.Errorf("err: expected %#v, got %#v", tc.wantErr, gotErr)
			}
			if gotMTOpenTerminal != tc.wantMTOpenTerminal {
				t.Errorf("mt open: expected %t, got %t", tc.wantMTOpenTerminal, gotMTOpenTerminal)
			}
			if !bytes.Equal(gotOut, tc.wantOut) {
				t.Errorf("mt out: expected %#v, got %#v", tc.wantOut, gotOut)
			}
		})
	}
}

func TestTermSlotCapture(t *testing.T) {
	errDummy := errors.New("dummy error")
	pid := os.Getpid()

	fg := vterm.NewColorIndexed(7)
	fg.Type |= vterm.ColorDefaultFG
	bg := vterm.NewColorIndexed(0)
	bg.Type |= vterm.ColorDefaultBG

	tt := []struct {
		name               string
		inStart            bool
		inIn               []byte
		inMTErrOpen        error
		wantSS             vterm.ScreenShot
		wantErr            error
		wantMTOpenTerminal bool
	}{
		{
			name:        "Normal",
			inStart:     true,
			inIn:        []byte("ABCDEF"),
			inMTErrOpen: errDummy,
			wantSS: vterm.ScreenShot{
				Stride: 3,
				Cell: []vterm.Cell{
					{
						Runes: []rune{'A'},
						Width: 1,
						FG:    fg,
						BG:    bg,
					},
					{
						Runes: []rune{'B'},
						Width: 1,
						FG:    fg,
						BG:    bg,
					},
					{
						Runes: []rune{'C'},
						Width: 1,
						FG:    fg,
						BG:    bg,
					},
					{
						Runes: []rune{'D'},
						Width: 1,
						FG:    fg,
						BG:    bg,
					},
					{
						Runes: []rune{'E'},
						Width: 1,
						FG:    fg,
						BG:    bg,
					},
					{
						Runes: []rune{'F'},
						Width: 1,
						FG:    fg,
						BG:    bg,
					},
				},
				CursorPos: vterm.Pos{
					Row: 1,
					Col: 2,
				},
				CursorVisible: true,
				CursorBlink:   true,
				CursorShape:   vterm.CursorShapeBlock,
			},
			wantErr:            nil,
			wantMTOpenTerminal: true,
		},
		{
			name:        "Start",
			inStart:     false,
			inMTErrOpen: nil,
			wantSS: vterm.ScreenShot{
				Stride: 3,
				Cell: []vterm.Cell{
					{
						Runes: []rune{},
						Width: 1,
						FG:    fg,
						BG:    bg,
					},
					{
						Runes: []rune{},
						Width: 1,
						FG:    fg,
						BG:    bg,
					},
					{
						Runes: []rune{},
						Width: 1,
						FG:    fg,
						BG:    bg,
					},
					{
						Runes: []rune{},
						Width: 1,
						FG:    fg,
						BG:    bg,
					},
					{
						Runes: []rune{},
						Width: 1,
						FG:    fg,
						BG:    bg,
					},
					{
						Runes: []rune{},
						Width: 1,
						FG:    fg,
						BG:    bg,
					},
				},
				CursorPos: vterm.Pos{
					Row: 0,
					Col: 0,
				},
				CursorVisible: true,
				CursorBlink:   true,
				CursorShape:   vterm.CursorShapeBlock,
			},
			wantErr:            nil,
			wantMTOpenTerminal: true,
		},
		{
			name:               "StartError",
			inStart:            false,
			inMTErrOpen:        errDummy,
			wantSS:             vterm.ScreenShot{},
			wantErr:            errDummy,
			wantMTOpenTerminal: false,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mt := &xpty.MockTerminal{PID: pid}
			mc := mt.Computer()
			cfg := TermConfig{
				Open: mt.Open,
				Row:  2,
				Col:  3,
				Cmd: xpty.Cmd{
					Path: "bash",
					Args: []string{"--version"},
				},
			}
			slot := NewTermSlot(cfg)

			if tc.inStart {
				var err error

				err = slot.start()
				if err != nil {
					t.Fatal(err)
				}

				_, err = mc.Write(tc.inIn)
				if err != nil {
					t.Fatal(err)
				}

				_, err = mc.Write([]byte{})
				if err != nil {
					t.Fatal(err)
				}
			}

			mt.ErrOpen = tc.inMTErrOpen
			gotSS, gotErr := slot.Capture()
			gotMTOpenTerminal := mt.OpenTerminal

			if slot.term != nil {
				slot.term.pc = nil
			}
			slot.Stop()

			if !reflect.DeepEqual(gotSS, tc.wantSS) {
				t.Errorf("ss: expected %#v, got %#v", tc.wantSS, gotSS)
			}
			if gotErr != tc.wantErr {
				t.Errorf("err: expected %#v, got %#v", tc.wantErr, gotErr)
			}
			if gotMTOpenTerminal != tc.wantMTOpenTerminal {
				t.Errorf("mt open: expected %t, got %t", tc.wantMTOpenTerminal, gotMTOpenTerminal)
			}
		})
	}
}

func TestTermSlotStop(t *testing.T) {
	pid := os.Getpid()
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

	slot := NewTermSlot(cfg)
	err := slot.start()
	if err != nil {
		t.Fatal(err)
	}
	slot.term.pc = nil

	slot.Stop()
	if mt.OpenTerminal {
		t.Errorf("slot.term open")
	}
	if slot.term != nil {
		t.Errorf("slot.term non-nil")
	}
}

func TestTermSlotStopStop(t *testing.T) {
	pid := os.Getpid()
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

	slot := NewTermSlot(cfg)
	slot.Stop()
	if mt.OpenTerminal {
		t.Errorf("slot.term open")
	}
	if slot.term != nil {
		t.Errorf("slot.term non-nil")
	}
}

func TestTermSlotStart(t *testing.T) {
	errDummy := errors.New("dummy error")
	pid := os.Getpid()

	tt := []struct {
		name            string
		inErrOpen1      error
		inErrOpen2      error
		wantErr1        error
		wantTermNonNil1 bool
		wantErr2        error
		wantTermNonNil2 bool
	}{
		{
			name:            "Normal",
			inErrOpen1:      nil,
			inErrOpen2:      nil,
			wantErr1:        nil,
			wantTermNonNil1: true,
			wantErr2:        nil,
			wantTermNonNil2: true,
		},
		{
			name:            "Err1",
			inErrOpen1:      errDummy,
			inErrOpen2:      nil,
			wantErr1:        errDummy,
			wantTermNonNil1: false,
			wantErr2:        nil,
			wantTermNonNil2: true,
		},
		{
			name:            "Err2",
			inErrOpen1:      nil,
			inErrOpen2:      errDummy,
			wantErr1:        nil,
			wantTermNonNil1: true,
			wantErr2:        nil,
			wantTermNonNil2: true,
		},
		{
			name:            "Err3",
			inErrOpen1:      errDummy,
			inErrOpen2:      errDummy,
			wantErr1:        errDummy,
			wantTermNonNil1: false,
			wantErr2:        errDummy,
			wantTermNonNil2: false,
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
			slot := NewTermSlot(cfg)

			mt.ErrOpen = tc.inErrOpen1
			gotErr1 := slot.start()
			if gotErr1 != tc.wantErr1 {
				t.Errorf("err1: expected %#v, got %#v", tc.wantErr1, gotErr1)
			}

			gotTerm1 := slot.term
			gotTermNonNil1 := slot.term != nil
			if gotTermNonNil1 != tc.wantTermNonNil1 {
				t.Errorf("slot.term non-nil 1: expected %t, got %t", tc.wantTermNonNil1, gotTermNonNil1)
			}

			if slot.term != nil && (mt.Size.Row != cfg.Row || mt.Size.Col != cfg.Col) {
				t.Errorf("slot.term size: expected (%d, %d), got (%d, %d)", cfg.Row, cfg.Col, mt.Size.Row, mt.Size.Col)
			}
			if slot.term != nil && !reflect.DeepEqual(mt.Cmd, cfg.Cmd) {
				t.Errorf("slot.term cmd: expected %#v, got %#v", cfg.Cmd, mt.Cmd)
			}

			mt.ErrOpen = tc.inErrOpen2
			gotErr2 := slot.start()
			if gotErr2 != tc.wantErr2 {
				t.Errorf("err2: expected %#v, got %#v", tc.wantErr2, gotErr2)
			}

			gotTerm2 := slot.term
			gotTermNonNil2 := slot.term != nil
			if gotTermNonNil2 != tc.wantTermNonNil2 {
				t.Errorf("slot.term non-nil 2: expected %t, got %t", tc.wantTermNonNil2, gotTermNonNil2)
			}

			if slot.term != nil && (mt.Size.Row != cfg.Row || mt.Size.Col != cfg.Col) {
				t.Errorf("slot.term size: expected (%d, %d), got (%d, %d)", cfg.Row, cfg.Col, mt.Size.Row, mt.Size.Col)
			}
			if slot.term != nil && !reflect.DeepEqual(mt.Cmd, cfg.Cmd) {
				t.Errorf("slot.term cmd: expected %#v, got %#v", cfg.Cmd, mt.Cmd)
			}

			if gotTerm1 != nil && gotTerm2 != nil && gotTerm1 != gotTerm2 {
				t.Errorf("slot.term reinitialized")
			}
		})
	}
}

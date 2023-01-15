package vterm

import (
	"bytes"
	"io"
	"testing"
)

func TestVTermKeyboardRune(t *testing.T) {
	tt := []struct {
		name    string
		inR     rune
		inMod   Modifier
		wantOut []byte
	}{
		{
			name:    "NormalSpace",
			inR:     ' ',
			inMod:   0,
			wantOut: []byte(" "),
		},
		{
			name:    "ShiftSpace",
			inR:     ' ',
			inMod:   ModShift,
			wantOut: []byte("\x1B[32;2u"),
		},
		{
			name:    "AltSpace",
			inR:     ' ',
			inMod:   ModAlt,
			wantOut: []byte("\x1B "),
		},
		{
			name:    "CtrlSpace",
			inR:     ' ',
			inMod:   ModCtrl,
			wantOut: []byte("\x00"),
		},
		{
			name:    "Invalid",
			inR:     -1,
			inMod:   0,
			wantOut: []byte{},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			vt := New(30, 120)
			out := vt.Output()

			var gotOut []byte
			done := make(chan struct{})
			go func() {
				gotOut, _ = io.ReadAll(out)
				close(done)
			}()

			vt.KeyboardRune(tc.inR, tc.inMod)
			_ = out.Close()
			<-done

			if !bytes.Equal(gotOut, tc.wantOut) {
				t.Errorf("expected %#v, got %#v", tc.wantOut, gotOut)
			}
		})
	}
}

func TestVTermKeyboardKey(t *testing.T) {
	tt := []struct {
		name    string
		inKey   Key
		inMod   Modifier
		wantOut []byte
	}{
		{
			name:    "NormalEnter",
			inKey:   KeyEnter,
			inMod:   0,
			wantOut: []byte("\r"),
		},
		{
			name:    "ShiftTab",
			inKey:   KeyEnter,
			inMod:   ModShift,
			wantOut: []byte("\x1B[13;2u"),
		},
		{
			name:    "AltTab",
			inKey:   KeyEnter,
			inMod:   ModAlt,
			wantOut: []byte("\x1B\r"),
		},
		{
			name:    "CtrlTab",
			inKey:   KeyEnter,
			inMod:   ModCtrl,
			wantOut: []byte("\x1B[13;5u"),
		},
		{
			name:    "Invalid",
			inKey:   -1,
			inMod:   0,
			wantOut: []byte{},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			vt := New(30, 120)
			out := vt.Output()

			var gotOut []byte
			done := make(chan struct{})
			go func() {
				gotOut, _ = io.ReadAll(out)
				close(done)
			}()

			vt.KeyboardKey(tc.inKey, tc.inMod)
			_ = out.Close()
			<-done

			if !bytes.Equal(gotOut, tc.wantOut) {
				t.Errorf("expected %#v, got %#v", tc.wantOut, gotOut)
			}
		})
	}
}

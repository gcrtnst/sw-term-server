package vterm

import (
	"bytes"
	"io"
	"testing"
)

func TestInputWrite(t *testing.T) {
	tt := []struct {
		name     string
		inInP    []byte
		wantOutP []byte
	}{
		{
			name:     "Empty",
			inInP:    []byte{},
			wantOutP: []byte{},
		},
		{
			name:     "DA",
			inInP:    []byte{0x1B, ' ', 'F', 0x1B, '[', '0', 'c'}, // S7C1T, DA
			wantOutP: []byte{0x1B, '[', '?', '1', ';', '2', 'c'},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			vt := New(30, 120)
			in := vt.Input()
			out := vt.Output()

			var gotOutP []byte
			var gotOutErr error
			done := make(chan struct{})
			go func() {
				gotOutP, gotOutErr = io.ReadAll(out)
				close(done)
			}()

			gotInN, gotInErr := in.Write(tc.inInP)
			_ = out.Close()
			<-done

			if gotInN != len(tc.inInP) {
				t.Errorf("in n: expected %d, got %d", len(tc.inInP), gotInN)
			}
			if gotInErr != nil {
				t.Errorf("in err: %v", gotInErr)
			}
			if !bytes.Equal(gotOutP, tc.wantOutP) {
				t.Errorf("out b: expected %#v, got %#v", tc.wantOutP, gotOutP)
			}
			if gotOutErr != nil {
				t.Errorf("out err: %v", gotOutErr)
			}
		})
	}
}

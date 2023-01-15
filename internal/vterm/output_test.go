package vterm

import (
	"bytes"
	"io"
	"testing"
)

func TestOutputMultiple(t *testing.T) {
	vt := New(30, 120)
	out := vt.Output()

	var gotP []byte
	var gotErr error
	done := make(chan struct{})
	go func() {
		gotP, gotErr = io.ReadAll(out)
		close(done)
	}()

	vt.KeyboardRune('a', ModNone)
	vt.KeyboardRune('b', ModNone)
	_ = out.Close()
	<-done

	wantP := []byte("ab")
	if !bytes.Equal(gotP, wantP) {
		t.Errorf("p: expected %#v, got %#v", wantP, gotP)
	}
	if gotErr != nil {
		t.Errorf("err: %v", gotErr)
	}
}

func TestOutputClose(t *testing.T) {
	vt := New(30, 120)
	out := vt.Output()

	errClose := out.Close()
	if errClose != nil {
		t.Fatalf("err: %v", errClose)
	}

	vt.KeyboardRune('a', ModNone)

	p := make([]byte, 1)
	n, errRead := out.Read(p)
	if n != 0 {
		t.Errorf("n: expected 0, got %v", n)
	}
	if errRead != io.EOF {
		t.Errorf("errRead: expected io.EOF, got %#v", errRead)
	}
}

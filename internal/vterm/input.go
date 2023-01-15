package vterm

// #include <vterm.h>
import "C"
import (
	"errors"
	"unsafe"
)

var ErrInputTooLarge = errors.New("vterm: input data too large")

type Input struct {
	vt *VTerm
}

func (in *Input) Write(p []byte) (int, error) {
	if len(p) <= 0 {
		return 0, nil
	}

	c_bytes := (*C.char)(unsafe.Pointer(&p[0]))
	c_len, ok := go2cSize(len(p))
	if !ok {
		return 0, ErrInputTooLarge
	}

	in.vt.mu.Lock()
	defer in.vt.mu.Unlock()

	c_ret := C.vterm_input_write(in.vt.vt, c_bytes, c_len)
	if c_ret != c_len {
		panic("vterm: vterm_input_write did not consume enough data")
	}

	in.vt.flush()

	return len(p), nil
}

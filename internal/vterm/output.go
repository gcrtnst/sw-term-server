package vterm

// #include <vterm.h>
import "C"
import (
	"io"
	"unsafe"
)

type Output struct {
	r *io.PipeReader
	w *io.PipeWriter
}

func newOutput() *Output {
	r, w := io.Pipe()
	return &Output{r: r, w: w}
}

func (out *Output) Read(p []byte) (int, error) {
	return out.r.Read(p)
}

func (out *Output) Close() error {
	return out.w.Close()
}

func (out *Output) flush(vt *C.VTerm) {
	cur := C.vterm_output_get_buffer_current(vt)
	if cur <= 0 {
		return
	}

	buf := make([]byte, cur)
	ptr := (*C.char)(unsafe.Pointer(&buf[0]))
	got := C.vterm_output_read(vt, ptr, cur)
	if got <= 0 {
		return
	}
	buf = buf[:got]

	_, _ = out.w.Write(buf)
}

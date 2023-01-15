package vterm

// #cgo pkg-config: vterm
// #include <vterm.h>
import "C"
import (
	"runtime"
	"sync"
)

type VTerm struct {
	mu  sync.Mutex
	vt  *C.VTerm
	out *Output
}

func New(rows, cols int) *VTerm {
	if rows < 1 {
		rows = 1
	}
	if cols < 1 {
		cols = 1
	}

	c_rows, _ := go2cInt(rows)
	c_cols, _ := go2cInt(cols)
	c_vt := C.vterm_new(c_rows, c_cols)

	vt := &VTerm{
		vt:  c_vt,
		out: newOutput(),
	}
	runtime.SetFinalizer(vt, (*VTerm).free)

	vt.Screen().init()
	return vt
}

func (vt *VTerm) Input() *Input {
	return &Input{vt: vt}
}

func (vt *VTerm) Output() *Output {
	return vt.out
}

func (vt *VTerm) Screen() *Screen {
	return &Screen{vt: vt}
}

func (vt *VTerm) GetSize() (int, int) {
	vt.mu.Lock()
	defer vt.mu.Unlock()

	return vt.size()
}

func (vt *VTerm) SetSize(rows, cols int) {
	vt.mu.Lock()
	defer vt.mu.Unlock()

	if rows < 1 {
		rows = 1
	}
	if cols < 1 {
		cols = 1
	}

	c_rows, _ := go2cInt(rows)
	c_cols, _ := go2cInt(cols)
	C.vterm_set_size(vt.vt, c_rows, c_cols)
}

func (vt *VTerm) GetUTF8() bool {
	vt.mu.Lock()
	defer vt.mu.Unlock()

	c := C.vterm_get_utf8(vt.vt)
	return c != 0
}

func (vt *VTerm) SetUTF8(isUTF8 bool) {
	vt.mu.Lock()
	defer vt.mu.Unlock()

	c := go2cBool(isUTF8)
	C.vterm_set_utf8(vt.vt, c)
}

func (vt *VTerm) size() (int, int) {
	var c_rows, c_cols C.int
	C.vterm_get_size(vt.vt, &c_rows, &c_cols)

	rows, _ := c2goInt(c_rows)
	cols, _ := c2goInt(c_cols)
	return rows, cols
}

func (vt *VTerm) flush() {
	vt.out.flush(vt.vt)
}

func (vt *VTerm) free() {
	vt.Screen().free()

	C.vterm_free(vt.vt)
	vt.vt = nil
	vt.out = nil
}

type Pos struct {
	Row int
	Col int
}

func newPosFromC(pos C.VTermPos) Pos {
	row, _ := c2goInt(pos.row)
	col, _ := c2goInt(pos.col)
	return Pos{Row: row, Col: col}
}

func (pos Pos) toC() C.VTermPos {
	row, _ := go2cInt(pos.Row)
	col, _ := go2cInt(pos.Col)
	return C.VTermPos{row: row, col: col}
}

package vterm

// #include <stdlib.h>
// #include <string.h>
// #include <vterm.h>
// #include <cgo_vterm_screen.h>
import "C"
import "unicode"

type Screen struct {
	vt *VTerm
}

func (scr *Screen) SetReflow(reflow bool) {
	scr.vt.mu.Lock()
	defer scr.vt.mu.Unlock()

	c_screen := scr.obtain()
	c_reflow := C.bool(reflow)
	C.vterm_screen_enable_reflow(c_screen, c_reflow)
}

func (scr *Screen) SetAltScreen(altscreen bool) {
	scr.vt.mu.Lock()
	defer scr.vt.mu.Unlock()

	c_screen := scr.obtain()
	c_altscreen := go2cBool(altscreen)
	C.vterm_screen_enable_altscreen(c_screen, c_altscreen)
}

func (scr *Screen) SetDefaultColor(fg, bg Color) {
	scr.vt.mu.Lock()
	defer scr.vt.mu.Unlock()

	c_screen := scr.obtain()
	c_fg := fg.toC()
	c_bg := bg.toC()
	C.vterm_screen_set_default_colors(c_screen, &c_fg, &c_bg)
}

func (scr *Screen) Capture() ScreenShot {
	scr.vt.mu.Lock()
	defer scr.vt.mu.Unlock()

	rows, cols := scr.vt.size()
	if rows < 0 {
		panic("row < 0")
	}
	if cols < 0 {
		panic("col < 0")
	}

	cell := make([]Cell, rows*cols)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			pos := Pos{Row: row, Col: col}
			c, ok := scr.cell(pos)
			if !ok {
				panic("cell not found")
			}
			cell[row*cols+col] = c
		}
	}

	c_user := scr.cbdata()
	return ScreenShot{
		Stride:        cols,
		Cell:          cell,
		CursorPos:     newPosFromC(c_user.cursor_pos),
		CursorVisible: c_user.cursor_visible != 0,
		CursorBlink:   c_user.cursor_blink != 0,
		CursorShape:   CursorShape(c_user.cursor_shape),
	}
}

func (scr *Screen) Cell(pos Pos) (Cell, bool) {
	scr.vt.mu.Lock()
	defer scr.vt.mu.Unlock()

	return scr.cell(pos)
}

func (scr *Screen) CursorPos() Pos {
	scr.vt.mu.Lock()
	defer scr.vt.mu.Unlock()

	c_user := scr.cbdata()
	return newPosFromC(c_user.cursor_pos)
}

func (scr *Screen) CursorVisible() bool {
	scr.vt.mu.Lock()
	defer scr.vt.mu.Unlock()

	c_user := scr.cbdata()
	return c_user.cursor_visible != 0
}

func (scr *Screen) CursorBlink() bool {
	scr.vt.mu.Lock()
	defer scr.vt.mu.Unlock()

	c_user := scr.cbdata()
	return c_user.cursor_blink != 0
}

func (scr *Screen) CursorShape() CursorShape {
	scr.vt.mu.Lock()
	defer scr.vt.mu.Unlock()

	c_user := scr.cbdata()
	return CursorShape(c_user.cursor_shape)
}

func (scr *Screen) cell(pos Pos) (Cell, bool) {
	c_screen := scr.obtain()
	c_pos := pos.toC()
	c_cell := new(C.VTermScreenCell)
	c_ok := C.vterm_screen_get_cell(c_screen, c_pos, c_cell)
	if c_ok == 0 {
		return Cell{}, false
	}

	cell := newCellFromC(c_cell)
	return cell, true
}

func (scr *Screen) cbdata() *C.CGoVTermScreenUser {
	c_screen := scr.obtain()
	c_user := C.vterm_screen_get_cbdata(c_screen)
	return (*C.CGoVTermScreenUser)(c_user)
}

func (scr *Screen) init() {
	c_user := C.malloc(C.sizeof_CGoVTermScreenUser)
	_ = C.memset(c_user, 0, C.sizeof_CGoVTermScreenUser)

	c_screen := scr.obtain()
	C.vterm_screen_set_callbacks(c_screen, &C.cgo_vterm_screen_user_callbacks, c_user)
	C.vterm_screen_set_damage_merge(c_screen, C.VTERM_DAMAGE_CELL)
	C.vterm_screen_reset(c_screen, 1)
}

func (scr *Screen) free() {
	c_screen := scr.obtain()
	c_user := C.vterm_screen_get_cbdata(c_screen)
	C.vterm_screen_set_callbacks(c_screen, nil, nil)
	C.free(c_user)
}

func (scr *Screen) obtain() *C.VTermScreen {
	return C.vterm_obtain_screen(scr.vt.vt)
}

type ScreenShot struct {
	Stride int
	Cell   []Cell

	CursorPos     Pos
	CursorVisible bool
	CursorBlink   bool
	CursorShape   CursorShape
}

func (ss ScreenShot) Size() (int, int) {
	if ss.Stride <= 0 || len(ss.Cell)%ss.Stride != 0 {
		return 0, 0
	}

	row := len(ss.Cell) / ss.Stride
	col := ss.Stride
	return row, col
}

func (ss ScreenShot) At(pos Pos) Cell {
	if pos.Row < 0 || pos.Col < 0 {
		return Cell{}
	}

	rows, cols := ss.Size()
	if rows <= pos.Row || cols <= pos.Col {
		return Cell{}
	}

	idx := pos.Row*ss.Stride + pos.Col
	return ss.Cell[idx]
}

type Cell struct {
	Runes  []rune
	Width  int
	Attrs  CellAttrs
	FG, BG Color
}

func newCellFromC(cell *C.VTermScreenCell) Cell {
	runes := make([]rune, 0, len(cell.chars))
	for i := 0; i < len(cell.chars); i++ {
		c := cell.chars[i]
		if c <= 0 || unicode.MaxRune < c {
			break
		}

		r := rune(c)
		runes = append(runes, r)
	}

	return Cell{
		Runes: runes,
		Width: int(cell.width),
		Attrs: newCellAttrsFromC(cell.attrs),
		FG:    newColorFromC(cell.fg),
		BG:    newColorFromC(cell.bg),
	}
}

type CellAttrs struct {
	Bold      bool
	Underline Underline
	Italic    bool
	Blink     bool
	Reverse   bool
	Conceal   bool
	Strike    bool
	Font      int
	DWL       bool
	DHL       DHL
	Small     bool
	Baseline  Baseline
}

func newCellAttrsFromC(attrs C.VTermScreenCellAttrs) CellAttrs {
	return CellAttrs{
		Bold:      C.cgo_vterm_screen_attrs_bold(attrs) != 0,
		Underline: Underline(C.cgo_vterm_screen_attrs_underline(attrs)),
		Italic:    C.cgo_vterm_screen_attrs_italic(attrs) != 0,
		Blink:     C.cgo_vterm_screen_attrs_blink(attrs) != 0,
		Reverse:   C.cgo_vterm_screen_attrs_reverse(attrs) != 0,
		Conceal:   C.cgo_vterm_screen_attrs_conceal(attrs) != 0,
		Strike:    C.cgo_vterm_screen_attrs_strike(attrs) != 0,
		Font:      int(C.cgo_vterm_screen_attrs_font(attrs)),
		DWL:       C.cgo_vterm_screen_attrs_dwl(attrs) != 0,
		DHL:       DHL(C.cgo_vterm_screen_attrs_dhl(attrs)),
		Small:     C.cgo_vterm_screen_attrs_small(attrs) != 0,
		Baseline:  Baseline(C.cgo_vterm_screen_attrs_baseline(attrs)),
	}
}

type Underline uint8

const (
	UnderlineOff    Underline = C.VTERM_UNDERLINE_OFF
	UnderlineSingle Underline = C.VTERM_UNDERLINE_SINGLE
	UnderlineDouble Underline = C.VTERM_UNDERLINE_DOUBLE
	UnderlineCurly  Underline = C.VTERM_UNDERLINE_CURLY
)

type DHL uint8

const (
	DHLOff    DHL = 0
	DHLTop    DHL = 1
	DHLBottom DHL = 2
)

type Baseline uint8

const (
	BaselineNormal Baseline = C.VTERM_BASELINE_NORMAL
	BaselineRaise  Baseline = C.VTERM_BASELINE_RAISE
	BaselineLower  Baseline = C.VTERM_BASELINE_LOWER
)

type CursorShape uint8

const (
	CursorShapeBlock     CursorShape = C.VTERM_PROP_CURSORSHAPE_BLOCK
	CursorShapeUnderline CursorShape = C.VTERM_PROP_CURSORSHAPE_UNDERLINE
	CursorShapeBarLeft   CursorShape = C.VTERM_PROP_CURSORSHAPE_BAR_LEFT
)

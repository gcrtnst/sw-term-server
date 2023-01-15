package vterm

// #include <vterm.h>
// #include <cgo_vterm_color.h>
import "C"

type ColorType uint8

const (
	ColorRGB         ColorType = C.VTERM_COLOR_RGB
	ColorIndexed     ColorType = C.VTERM_COLOR_INDEXED
	ColorTypeMask    ColorType = C.VTERM_COLOR_TYPE_MASK
	ColorDefaultFG   ColorType = C.VTERM_COLOR_DEFAULT_FG
	ColorDefaultBG   ColorType = C.VTERM_COLOR_DEFAULT_BG
	ColorDefaultMask ColorType = C.VTERM_COLOR_DEFAULT_MASK
)

type Color struct {
	Type             ColorType
	Red, Green, Blue uint8
	Idx              uint8
}

func NewColorRGB(red, green, blue uint8) Color {
	return Color{
		Type:  ColorRGB,
		Red:   red,
		Green: green,
		Blue:  blue,
	}
}

func NewColorIndexed(idx uint8) Color {
	return Color{
		Type: ColorIndexed,
		Idx:  idx,
	}
}

func (col Color) IsIndexed() bool {
	return (col.Type & ColorTypeMask) == ColorIndexed
}

func (col Color) IsRGB() bool {
	return (col.Type & ColorTypeMask) == ColorRGB
}

func (col Color) IsDefaultFG() bool {
	return (col.Type & ColorDefaultFG) != 0
}

func (col Color) IsDefaultBG() bool {
	return (col.Type & ColorDefaultBG) != 0
}

func (a Color) Equal(b Color) bool {
	ca := a.toC()
	cb := b.toC()
	return C.vterm_color_is_equal(&ca, &cb) != 0
}

func newColorFromC(col C.VTermColor) Color {
	return Color{
		Type:  ColorType(C.cgo_vterm_color_type(&col)),
		Red:   uint8(C.cgo_vterm_color_red(&col)),
		Green: uint8(C.cgo_vterm_color_green(&col)),
		Blue:  uint8(C.cgo_vterm_color_blue(&col)),
		Idx:   uint8(C.cgo_vterm_color_idx(&col)),
	}
}

func (col Color) toC() C.VTermColor {
	c_type := C.uint8_t(col.Type)
	c_red := C.uint8_t(col.Red)
	c_green := C.uint8_t(col.Green)
	c_blue := C.uint8_t(col.Blue)
	c_idx := C.uint8_t(col.Idx)

	var c_col C.VTermColor
	C.cgo_vterm_color(&c_col, c_type, c_red, c_green, c_blue, c_idx)
	return c_col
}

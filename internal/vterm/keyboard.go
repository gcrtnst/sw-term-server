package vterm

// #include <vterm.h>
import "C"
import "unicode"

type Modifier uint8

const (
	ModNone  Modifier = C.VTERM_MOD_NONE
	ModShift Modifier = C.VTERM_MOD_SHIFT
	ModAlt   Modifier = C.VTERM_MOD_ALT
	ModCtrl  Modifier = C.VTERM_MOD_CTRL
	ModAll   Modifier = C.VTERM_ALL_MODS_MASK
)

type Key int

const (
	KeyNone        Key = C.VTERM_KEY_NONE
	KeyEnter       Key = C.VTERM_KEY_ENTER
	KeyTab         Key = C.VTERM_KEY_TAB
	KeyBackspace   Key = C.VTERM_KEY_BACKSPACE
	KeyEscape      Key = C.VTERM_KEY_ESCAPE
	KeyUp          Key = C.VTERM_KEY_UP
	KeyDown        Key = C.VTERM_KEY_DOWN
	KeyLeft        Key = C.VTERM_KEY_LEFT
	KeyRight       Key = C.VTERM_KEY_RIGHT
	KeyIns         Key = C.VTERM_KEY_INS
	KeyDel         Key = C.VTERM_KEY_DEL
	KeyHome        Key = C.VTERM_KEY_HOME
	KeyEnd         Key = C.VTERM_KEY_END
	KeyPageup      Key = C.VTERM_KEY_PAGEUP
	KeyPagedown    Key = C.VTERM_KEY_PAGEDOWN
	KeyFunction0   Key = C.VTERM_KEY_FUNCTION_0
	KeyFunction1   Key = KeyFunction0 + 1
	KeyFunction2   Key = KeyFunction0 + 2
	KeyFunction3   Key = KeyFunction0 + 3
	KeyFunction4   Key = KeyFunction0 + 4
	KeyFunction5   Key = KeyFunction0 + 5
	KeyFunction6   Key = KeyFunction0 + 6
	KeyFunction7   Key = KeyFunction0 + 7
	KeyFunction8   Key = KeyFunction0 + 8
	KeyFunction9   Key = KeyFunction0 + 9
	KeyFunction10  Key = KeyFunction0 + 10
	KeyFunction11  Key = KeyFunction0 + 11
	KeyFunction12  Key = KeyFunction0 + 12
	KeyFunctionMax Key = C.VTERM_KEY_FUNCTION_MAX
	KeyKP0         Key = C.VTERM_KEY_KP_0
	KeyKP1         Key = C.VTERM_KEY_KP_1
	KeyKP2         Key = C.VTERM_KEY_KP_2
	KeyKP3         Key = C.VTERM_KEY_KP_3
	KeyKP4         Key = C.VTERM_KEY_KP_4
	KeyKP5         Key = C.VTERM_KEY_KP_5
	KeyKP6         Key = C.VTERM_KEY_KP_6
	KeyKP7         Key = C.VTERM_KEY_KP_7
	KeyKP8         Key = C.VTERM_KEY_KP_8
	KeyKP9         Key = C.VTERM_KEY_KP_9
	KeyKPMult      Key = C.VTERM_KEY_KP_MULT
	KeyKPPlus      Key = C.VTERM_KEY_KP_PLUS
	KeyKPComma     Key = C.VTERM_KEY_KP_COMMA
	KeyKPMinus     Key = C.VTERM_KEY_KP_MINUS
	KeyKPPeriod    Key = C.VTERM_KEY_KP_PERIOD
	KeyKPDivide    Key = C.VTERM_KEY_KP_DIVIDE
	KeyKPEnter     Key = C.VTERM_KEY_KP_ENTER
	KeyKPEqual     Key = C.VTERM_KEY_KP_EQUAL
	KeyMax         Key = C.VTERM_KEY_MAX
)

func (vt *VTerm) KeyboardRune(r rune, mod Modifier) {
	if r < 0 || unicode.MaxRune < r {
		return
	}

	vt.mu.Lock()
	defer vt.mu.Unlock()

	c_c := C.uint32_t(r)
	c_mod := C.VTermModifier(mod & ModAll)
	C.vterm_keyboard_unichar(vt.vt, c_c, c_mod)

	vt.flush()
}

func (vt *VTerm) KeyboardKey(key Key, mod Modifier) {
	if key < KeyNone || KeyMax < key {
		return
	}

	vt.mu.Lock()
	defer vt.mu.Unlock()

	c_key := C.VTermKey(key)
	c_mod := C.VTermModifier(mod & ModAll)
	C.vterm_keyboard_key(vt.vt, c_key, c_mod)

	vt.flush()
}

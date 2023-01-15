package main

import (
	"unicode/utf8"

	"github.com/gcrtnst/sw-term-server/internal/vterm"
)

type Key string

var vtermKeyMap = map[Key]vterm.Key{
	"Enter":      vterm.KeyEnter,
	"Tab":        vterm.KeyTab,
	"Backspace":  vterm.KeyBackspace,
	"Escape":     vterm.KeyEscape,
	"ArrowUp":    vterm.KeyUp,
	"ArrowDown":  vterm.KeyDown,
	"ArrowLeft":  vterm.KeyLeft,
	"ArrowRight": vterm.KeyRight,
	"Insert":     vterm.KeyIns,
	"Delete":     vterm.KeyDel,
	"Home":       vterm.KeyHome,
	"End":        vterm.KeyEnd,
	"PageUp":     vterm.KeyPageup,
	"PageDown":   vterm.KeyPagedown,
	"F1":         vterm.KeyFunction1,
	"F2":         vterm.KeyFunction2,
	"F3":         vterm.KeyFunction3,
	"F4":         vterm.KeyFunction4,
	"F5":         vterm.KeyFunction5,
	"F6":         vterm.KeyFunction6,
	"F7":         vterm.KeyFunction7,
	"F8":         vterm.KeyFunction8,
	"F9":         vterm.KeyFunction9,
	"F10":        vterm.KeyFunction10,
	"F11":        vterm.KeyFunction11,
	"F12":        vterm.KeyFunction12,
	"KP0":        vterm.KeyKP0,
	"KP1":        vterm.KeyKP1,
	"KP2":        vterm.KeyKP2,
	"KP3":        vterm.KeyKP3,
	"KP4":        vterm.KeyKP4,
	"KP5":        vterm.KeyKP5,
	"KP6":        vterm.KeyKP6,
	"KP7":        vterm.KeyKP7,
	"KP8":        vterm.KeyKP8,
	"KP9":        vterm.KeyKP9,
	"KP*":        vterm.KeyKPMult,
	"KP+":        vterm.KeyKPPlus,
	"KP,":        vterm.KeyKPComma,
	"KP-":        vterm.KeyKPMinus,
	"KP.":        vterm.KeyKPPeriod,
	"KP/":        vterm.KeyKPDivide,
	"KPEnter":    vterm.KeyKPEnter,
	"KP=":        vterm.KeyKPEqual,
}

func (k Key) Rune() (rune, bool) {
	r, size := utf8.DecodeRuneInString(string(k))
	if (r == utf8.RuneError && (size == 0 || size == 1)) || (size != len(k)) {
		return utf8.RuneError, false
	}
	return r, true
}

func (k Key) VTermKey() (vterm.Key, bool) {
	vk, ok := vtermKeyMap[k]
	return vk, ok
}

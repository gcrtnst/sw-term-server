package main

import (
	"testing"

	"github.com/gcrtnst/sw-term-server/internal/vterm"
)

func TestKeyVTermRune(t *testing.T) {
	tt := []struct {
		name     string
		inK      Key
		wantRune rune
		wantOK   bool
	}{
		{
			name:     "NormalSingleByte",
			inK:      "A",
			wantRune: 'A',
			wantOK:   true,
		},
		{
			name:     "NormalMultiByte",
			inK:      "あ",
			wantRune: 'あ',
			wantOK:   true,
		},
		{
			name:     "NormalInvalid",
			inK:      "\uFFFD",
			wantRune: 0xFFFD,
			wantOK:   true,
		},
		{
			name:     "ErrorEmpty",
			inK:      "",
			wantRune: 0xFFFD,
			wantOK:   false,
		},
		{
			name:     "ErrorInvalid",
			inK:      "\x80",
			wantRune: 0xFFFD,
			wantOK:   false,
		},
		{
			name:     "ErrorMultiRune",
			inK:      "AA",
			wantRune: 0xFFFD,
			wantOK:   false,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotRune, gotOK := tc.inK.Rune()
			if gotRune != tc.wantRune {
				t.Errorf("rune: expected %U, got %U", tc.wantRune, gotRune)
			}
			if gotOK != tc.wantOK {
				t.Errorf("ok: expected %t, got %t", tc.wantOK, gotOK)
			}
		})
	}
}

func TestKeyVTermKey(t *testing.T) {
	tt := []struct {
		name   string
		inK    Key
		wantVK vterm.Key
		wantOK bool
	}{
		{
			name:   "Normal",
			inK:    "Enter",
			wantVK: vterm.KeyEnter,
			wantOK: true,
		},
		{
			name:   "Error",
			inK:    "",
			wantVK: 0,
			wantOK: false,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotVK, gotOK := tc.inK.VTermKey()
			if gotVK != tc.wantVK {
				t.Errorf("vk: expected %d, got %d", tc.wantVK, gotVK)
			}
			if gotOK != tc.wantOK {
				t.Errorf("ok: expected %t, got %t", tc.wantOK, gotOK)
			}
		})
	}
}

package main

import (
	"reflect"
	"testing"

	"github.com/gcrtnst/sw-term-server/internal/vterm"
)

func TestEncodeScreenShot(t *testing.T) {
	tt := []struct {
		name string
		in   vterm.ScreenShot
		want [][]string
	}{
		{
			name: "Zero",
			in:   vterm.ScreenShot{},
			want: [][]string{
				{"#sw-term/screen"},
				{"screen", "0", "0"},
				{"cursor", "0", "0", "0", "0", "0"},
			},
		},
		{
			name: "CursorPos",
			in: vterm.ScreenShot{
				CursorPos: vterm.Pos{Row: 1, Col: 2},
			},
			want: [][]string{
				{"#sw-term/screen"},
				{"screen", "0", "0"},
				{"cursor", "0", "0", "0", "1", "2"},
			},
		},
		{
			name: "CursorVisible",
			in: vterm.ScreenShot{
				CursorVisible: true,
			},
			want: [][]string{
				{"#sw-term/screen"},
				{"screen", "0", "0"},
				{"cursor", "1", "0", "0", "0", "0"},
			},
		},
		{
			name: "CursorBlink",
			in: vterm.ScreenShot{
				CursorBlink: true,
			},
			want: [][]string{
				{"#sw-term/screen"},
				{"screen", "0", "0"},
				{"cursor", "0", "1", "0", "0", "0"},
			},
		},
		{
			name: "CursorShapeBlock",
			in: vterm.ScreenShot{
				CursorShape: vterm.CursorShapeBlock,
			},
			want: [][]string{
				{"#sw-term/screen"},
				{"screen", "0", "0"},
				{"cursor", "0", "0", "1", "0", "0"},
			},
		},
		{
			name: "CursorShapeUnderline",
			in: vterm.ScreenShot{
				CursorShape: vterm.CursorShapeUnderline,
			},
			want: [][]string{
				{"#sw-term/screen"},
				{"screen", "0", "0"},
				{"cursor", "0", "0", "2", "0", "0"},
			},
		},
		{
			name: "CursorShapeBarLeft",
			in: vterm.ScreenShot{
				CursorShape: vterm.CursorShapeBarLeft,
			},
			want: [][]string{
				{"#sw-term/screen"},
				{"screen", "0", "0"},
				{"cursor", "0", "0", "3", "0", "0"},
			},
		},
		{
			name: "Blank",
			in: vterm.ScreenShot{
				Stride: 2,
				Cell: []vterm.Cell{
					{
						Runes: []rune{},
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultFG,
							Idx:  7,
						},
						BG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
							Idx:  0,
						},
					},
					{
						Runes: []rune{},
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultFG,
							Idx:  7,
						},
						BG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
							Idx:  0,
						},
					},
					{
						Runes: []rune{},
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultFG,
							Idx:  7,
						},
						BG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
							Idx:  0,
						},
					},
					{
						Runes: []rune{},
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultFG,
							Idx:  7,
						},
						BG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
							Idx:  0,
						},
					},
					{
						Runes: []rune{},
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultFG,
							Idx:  7,
						},
						BG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
							Idx:  0,
						},
					},
					{
						Runes: []rune{},
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultFG,
							Idx:  7,
						},
						BG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
							Idx:  0,
						},
					},
				},
			},
			want: [][]string{
				{"#sw-term/screen"},
				{"screen", "3", "2"},
				{"cursor", "0", "0", "0", "0", "0"},
			},
		},
		{
			name: "Cell",
			in: vterm.ScreenShot{
				Stride: 2,
				Cell: []vterm.Cell{
					{
						Runes: []rune{'A'},
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG: vterm.Color{
							Type: vterm.ColorIndexed,
							Idx:  7,
						},
						BG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
							Idx:  0,
						},
					},
					{
						Runes: []rune{'B'},
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG: vterm.Color{
							Type: vterm.ColorIndexed,
							Idx:  6,
						},
						BG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
							Idx:  0,
						},
					},
					{
						Runes: []rune{'C'},
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG: vterm.Color{
							Type: vterm.ColorIndexed,
							Idx:  5,
						},
						BG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
							Idx:  0,
						},
					},
					{
						Runes: []rune{'D'},
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG: vterm.Color{
							Type: vterm.ColorIndexed,
							Idx:  4,
						},
						BG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
							Idx:  0,
						},
					},
					{
						Runes: []rune{'E'},
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG: vterm.Color{
							Type: vterm.ColorIndexed,
							Idx:  3,
						},
						BG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
							Idx:  0,
						},
					},
					{
						Runes: []rune{'F'},
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG: vterm.Color{
							Type: vterm.ColorIndexed,
							Idx:  2,
						},
						BG: vterm.Color{
							Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
							Idx:  0,
						},
					},
				},
			},
			want: [][]string{
				{"#sw-term/screen"},
				{"screen", "3", "2"},
				{"cursor", "0", "0", "0", "0", "0"},
				{"cell", "0", "0", "char", "1", "A"},
				{"cell", "0", "0", "color", "fg", "idx", "7"},
				{"cell", "0", "1", "char", "1", "B"},
				{"cell", "0", "1", "color", "fg", "idx", "6"},
				{"cell", "1", "0", "char", "1", "C"},
				{"cell", "1", "0", "color", "fg", "idx", "5"},
				{"cell", "1", "1", "char", "1", "D"},
				{"cell", "1", "1", "color", "fg", "idx", "4"},
				{"cell", "2", "0", "char", "1", "E"},
				{"cell", "2", "0", "color", "fg", "idx", "3"},
				{"cell", "2", "1", "char", "1", "F"},
				{"cell", "2", "1", "color", "fg", "idx", "2"},
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := EncodeScreenShot(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("expected %#v, got %#v", tc.want, got)
			}
		})
	}
}

func TestEncodeCell(t *testing.T) {
	tt := []struct {
		name string
		in   vterm.Cell
		want [][]string
	}{
		{
			name: "Default",
			in: vterm.Cell{
				Runes: []rune{},
				Width: 1,
				FG: vterm.Color{
					Type: vterm.ColorIndexed | vterm.ColorDefaultFG,
					Idx:  7,
				},
				BG: vterm.Color{
					Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
					Idx:  0,
				},
			},
			want: [][]string{},
		},
		{
			name: "DefaultWidth",
			in: vterm.Cell{
				Runes: []rune{'A'},
				Width: 0,
				FG: vterm.Color{
					Type: vterm.ColorIndexed | vterm.ColorDefaultFG,
					Idx:  7,
				},
				BG: vterm.Color{
					Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
					Idx:  0,
				},
			},
			want: [][]string{},
		},
		{
			name: "Char",
			in: vterm.Cell{
				Runes: []rune{'A'},
				Width: 1,
				FG: vterm.Color{
					Type: vterm.ColorIndexed | vterm.ColorDefaultFG,
					Idx:  7,
				},
				BG: vterm.Color{
					Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
					Idx:  0,
				},
			},
			want: [][]string{
				{"char", "1", "A"},
			},
		},
		{
			name: "Attr",
			in: vterm.Cell{
				Runes: []rune{},
				Width: 1,
				Attrs: vterm.CellAttrs{
					Bold:      true,
					Underline: vterm.UnderlineSingle,
					Italic:    true,
				},
				FG: vterm.Color{
					Type: vterm.ColorIndexed | vterm.ColorDefaultFG,
					Idx:  7,
				},
				BG: vterm.Color{
					Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
					Idx:  0,
				},
			},
			want: [][]string{
				{"attr", "bold"},
				{"attr", "underline", "1"},
				{"attr", "italic"},
			},
		},
		{
			name: "Color",
			in: vterm.Cell{
				Runes: []rune{},
				Width: 1,
				FG: vterm.Color{
					Type: vterm.ColorIndexed,
					Idx:  7,
				},
				BG: vterm.Color{
					Type: vterm.ColorIndexed,
					Idx:  0,
				},
			},
			want: [][]string{
				{"color", "fg", "idx", "7"},
				{"color", "bg", "idx", "0"},
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := encodeCell(tc.in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("expected %#v, got %#v", tc.want, got)
			}
		})
	}
}

func TestEncodeCellAttrs(t *testing.T) {
	tt := []struct {
		name string
		in   vterm.CellAttrs
		want [][]string
	}{
		{
			name: "Default",
			in:   vterm.CellAttrs{},
			want: [][]string{},
		},
		{
			name: "Bold",
			in:   vterm.CellAttrs{Bold: true},
			want: [][]string{{"bold"}},
		},
		{
			name: "UnderlineSingle",
			in:   vterm.CellAttrs{Underline: vterm.UnderlineSingle},
			want: [][]string{{"underline", "1"}},
		},
		{
			name: "UnderlineDouble",
			in:   vterm.CellAttrs{Underline: vterm.UnderlineDouble},
			want: [][]string{{"underline", "2"}},
		},
		{
			name: "UnderlineCurly",
			in:   vterm.CellAttrs{Underline: vterm.UnderlineCurly},
			want: [][]string{{"underline", "3"}},
		},
		{
			name: "Italic",
			in:   vterm.CellAttrs{Italic: true},
			want: [][]string{{"italic"}},
		},
		{
			name: "Blink",
			in:   vterm.CellAttrs{Blink: true},
			want: [][]string{{"blink"}},
		},
		{
			name: "Reverse",
			in:   vterm.CellAttrs{Reverse: true},
			want: [][]string{{"reverse"}},
		},
		{
			name: "Conceal",
			in:   vterm.CellAttrs{Conceal: true},
			want: [][]string{{"conceal"}},
		},
		{
			name: "Strike",
			in:   vterm.CellAttrs{Strike: true},
			want: [][]string{{"strike"}},
		},
		{
			name: "Font",
			in:   vterm.CellAttrs{Font: 9},
			want: [][]string{{"font", "9"}},
		},
		{
			name: "DWL",
			in:   vterm.CellAttrs{DWL: true},
			want: [][]string{{"dwl"}},
		},
		{
			name: "DHLTop",
			in:   vterm.CellAttrs{DHL: vterm.DHLTop},
			want: [][]string{{"dhl", "1"}},
		},
		{
			name: "DHLBottom",
			in:   vterm.CellAttrs{DHL: vterm.DHLBottom},
			want: [][]string{{"dhl", "2"}},
		},
		{
			name: "Small",
			in:   vterm.CellAttrs{Small: true},
			want: [][]string{{"small"}},
		},
		{
			name: "BaselineRaise",
			in:   vterm.CellAttrs{Baseline: vterm.BaselineRaise},
			want: [][]string{{"baseline", "1"}},
		},
		{
			name: "BaselineLower",
			in:   vterm.CellAttrs{Baseline: vterm.BaselineLower},
			want: [][]string{{"baseline", "2"}},
		},
		{
			name: "All",
			in: vterm.CellAttrs{
				Bold:      true,
				Underline: vterm.UnderlineCurly,
				Italic:    true,
				Blink:     true,
				Reverse:   true,
				Conceal:   true,
				Strike:    true,
				Font:      9,
				DWL:       true,
				DHL:       vterm.DHLTop,
				Small:     true,
				Baseline:  vterm.BaselineLower,
			},
			want: [][]string{
				{"bold"},
				{"underline", "3"},
				{"italic"},
				{"blink"},
				{"reverse"},
				{"conceal"},
				{"strike"},
				{"font", "9"},
				{"dwl"},
				{"dhl", "1"},
				{"small"},
				{"baseline", "2"},
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := encodeCellAttrs(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("expected %#v, got %#v", tc.want, got)
			}
		})
	}
}

func TestEncodeCellColor(t *testing.T) {
	tt := []struct {
		name       string
		inFG, inBG vterm.Color
		want       [][]string
	}{
		{
			name: "Default",
			inFG: vterm.Color{
				Type: vterm.ColorIndexed | vterm.ColorDefaultFG,
				Idx:  7,
			},
			inBG: vterm.Color{
				Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
				Idx:  0,
			},
			want: [][]string{},
		},
		{
			name: "FG",
			inFG: vterm.Color{
				Type: vterm.ColorIndexed,
				Idx:  7,
			},
			inBG: vterm.Color{
				Type: vterm.ColorIndexed | vterm.ColorDefaultBG,
				Idx:  0,
			},
			want: [][]string{{"fg", "idx", "7"}},
		},
		{
			name: "BG",
			inFG: vterm.Color{
				Type: vterm.ColorIndexed | vterm.ColorDefaultFG,
				Idx:  7,
			},
			inBG: vterm.Color{
				Type: vterm.ColorIndexed,
				Idx:  0,
			},
			want: [][]string{{"bg", "idx", "0"}},
		},
		{
			name: "All",
			inFG: vterm.Color{
				Type: vterm.ColorIndexed,
				Idx:  7,
			},
			inBG: vterm.Color{
				Type: vterm.ColorIndexed,
				Idx:  0,
			},
			want: [][]string{
				{"fg", "idx", "7"},
				{"bg", "idx", "0"},
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := encodeCellColor(tc.inFG, tc.inBG)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("expected %#v, got %#v", tc.want, got)
			}
		})
	}
}

func TestEncodeColor(t *testing.T) {
	tt := []struct {
		name string
		in   vterm.Color
		want []string
	}{
		{
			name: "RGB",
			in: vterm.Color{
				Type:  vterm.ColorRGB | ^vterm.ColorTypeMask,
				Red:   11,
				Green: 12,
				Blue:  13,
				Idx:   14,
			},
			want: []string{"rgb", "11", "12", "13"},
		},
		{
			name: "Indexed",
			in: vterm.Color{
				Type:  vterm.ColorIndexed | ^vterm.ColorTypeMask,
				Red:   11,
				Green: 12,
				Blue:  13,
				Idx:   14,
			},
			want: []string{"idx", "14"},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := encodeColor(tc.in)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("expected %#v, got %#v", tc.want, got)
			}
		})
	}
}

func TestEncodeBool(t *testing.T) {
	tt := []struct {
		name string
		in   bool
		want string
	}{
		{
			name: "False",
			in:   false,
			want: "0",
		},
		{
			name: "True",
			in:   true,
			want: "1",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := encodeBool(tc.in)
			if got != tc.want {
				t.Errorf("expected %#v, got %#v", tc.want, got)
			}
		})
	}
}

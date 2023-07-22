package main

import (
	"bytes"
	"testing"

	"github.com/gcrtnst/sw-term-server/internal/vterm"
)

func TestEncodeScreenShot(t *testing.T) {
	tt := []struct {
		name string
		in   vterm.ScreenShot
		want []byte
	}{
		{
			name: "Zero",
			in:   vterm.ScreenShot{},
			want: []byte{
				0x00,                                           // CursorVisible
				0x00,                                           // CursorBlink
				0x00,                                           // CursorShape
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // CursorPos.Row
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // CursorPos.Col
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Stride
			},
		},
		{
			name: "CursorVisible",
			in: vterm.ScreenShot{
				CursorVisible: true,
			},
			want: []byte{
				0x01,                                           // CursorVisible
				0x00,                                           // CursorBlink
				0x00,                                           // CursorShape
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // CursorPos.Row
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // CursorPos.Col
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Stride
			},
		},
		{
			name: "CursorBlink",
			in: vterm.ScreenShot{
				CursorBlink: true,
			},
			want: []byte{
				0x00,                                           // CursorVisible
				0x01,                                           // CursorBlink
				0x00,                                           // CursorShape
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // CursorPos.Row
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // CursorPos.Col
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Stride
			},
		},
		{
			name: "CursorShape",
			in: vterm.ScreenShot{
				CursorShape: vterm.CursorShapeBarLeft,
			},
			want: []byte{
				0x00,                                           // CursorVisible
				0x00,                                           // CursorBlink
				0x03,                                           // CursorShape
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // CursorPos.Row
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // CursorPos.Col
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Stride
			},
		},
		{
			name: "CursorPos",
			in: vterm.ScreenShot{
				CursorPos: vterm.Pos{
					Row: 1,
					Col: -2,
				},
			},
			want: []byte{
				0x00,                                           // CursorVisible
				0x00,                                           // CursorBlink
				0x00,                                           // CursorShape
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // CursorPos.Row
				0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // CursorPos.Col
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Stride
			},
		},
		{
			name: "Cell",
			in: vterm.ScreenShot{
				Stride: 5,
				Cell: []vterm.Cell{
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG:    vterm.NewColorIndexed(7),
						BG:    vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune{0x0041, 0x030A},
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG:    vterm.NewColorIndexed(7),
						BG:    vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune("あ"),
						Width: 2,
						Attrs: vterm.CellAttrs{},
						FG:    vterm.NewColorIndexed(7),
						BG:    vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune{},
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG:    vterm.NewColorIndexed(7),
						BG:    vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{
							Bold: true,
						},
						FG: vterm.NewColorIndexed(7),
						BG: vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{
							Underline: vterm.UnderlineCurly,
						},
						FG: vterm.NewColorIndexed(7),
						BG: vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{
							Italic: true,
						},
						FG: vterm.NewColorIndexed(7),
						BG: vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{
							Blink: true,
						},
						FG: vterm.NewColorIndexed(7),
						BG: vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{
							Reverse: true,
						},
						FG: vterm.NewColorIndexed(7),
						BG: vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{
							Conceal: true,
						},
						FG: vterm.NewColorIndexed(7),
						BG: vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{
							Strike: true,
						},
						FG: vterm.NewColorIndexed(7),
						BG: vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{
							Font: 9,
						},
						FG: vterm.NewColorIndexed(7),
						BG: vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{
							DWL: true,
						},
						FG: vterm.NewColorIndexed(7),
						BG: vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{
							DHL: vterm.DHLBottom,
						},
						FG: vterm.NewColorIndexed(7),
						BG: vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{
							Small: true,
						},
						FG: vterm.NewColorIndexed(7),
						BG: vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{
							Baseline: vterm.BaselineLower,
						},
						FG: vterm.NewColorIndexed(7),
						BG: vterm.NewColorIndexed(0),
					},
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG:    vterm.NewColorRGB(0, 1, 2),
						BG:    vterm.NewColorRGB(253, 254, 255),
					},
					{
						Runes: []rune("A"),
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
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG: vterm.Color{
							Type:  vterm.ColorRGB | vterm.ColorDefaultFG,
							Red:   0,
							Green: 1,
							Blue:  2,
						},
						BG: vterm.Color{
							Type:  vterm.ColorRGB | vterm.ColorDefaultBG,
							Red:   253,
							Green: 254,
							Blue:  255,
						},
					},
					{
						Runes: []rune("A"),
						Width: 1,
						Attrs: vterm.CellAttrs{},
						FG:    vterm.NewColorIndexed(7),
						BG:    vterm.NewColorIndexed(0),
					},
				},
			},
			want: []byte{
				0x00,                                           // CursorVisible
				0x00,                                           // CursorBlink
				0x00,                                           // CursorShape
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // CursorPos.Row
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // CursorPos.Col
				0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Stride

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				0x41, 0xCC, 0x8A, // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x02,                                           // Cell[i].Width
				0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				0xE3, 0x81, 0x82, // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)

				0x01,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x03,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x01,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x01,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x01,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x01,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x01,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x09,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x01,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x02,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x01,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x02,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x00,                                           // Cell[i].FG.Type
				0x00,                                           // Cell[i].FG.Red
				0x01,                                           // Cell[i].FG.Green
				0x02,                                           // Cell[i].FG.Blue
				0x00,                                           // Cell[i].BG.Type
				0xFD,                                           // Cell[i].BG.Red
				0xFE,                                           // Cell[i].BG.Green
				0xFF,                                           // Cell[i].BG.Blue
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x03,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x05,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x02,                                           // Cell[i].FG.Type
				0x00,                                           // Cell[i].FG.Red
				0x01,                                           // Cell[i].FG.Green
				0x02,                                           // Cell[i].FG.Blue
				0x04,                                           // Cell[i].BG.Type
				0xFD,                                           // Cell[i].BG.Red
				0xFE,                                           // Cell[i].BG.Green
				0xFF,                                           // Cell[i].BG.Blue
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)

				0x00,                                           // Cell[i].Attrs.Bold
				0x00,                                           // Cell[i].Attrs.Underline
				0x00,                                           // Cell[i].Attrs.Italic
				0x00,                                           // Cell[i].Attrs.Blink
				0x00,                                           // Cell[i].Attrs.Reverse
				0x00,                                           // Cell[i].Attrs.Conceal
				0x00,                                           // Cell[i].Attrs.Strike
				0x00,                                           // Cell[i].Attrs.Font
				0x00,                                           // Cell[i].Attrs.DWL
				0x00,                                           // Cell[i].Attrs.DHL
				0x00,                                           // Cell[i].Attrs.Small
				0x00,                                           // Cell[i].Attrs.Baseline
				0x01,                                           // Cell[i].FG.Type
				0x07,                                           // Cell[i].FG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].BG.Type
				0x00,                                           // Cell[i].BG.Idx
				0x00,                                           // padding
				0x00,                                           // padding
				0x01,                                           // Cell[i].Width
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // len(Cell[i].Rune)
				'A', // string(Cell[i].Rune)
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := EncodeScreenShot(tc.in)
			if !bytes.Equal(got, tc.want) {
				t.Errorf("expected %X, got %X", tc.want, got)
			}
		})
	}
}

func TestEscapeZero(t *testing.T) {
	tt := []struct {
		name string
		in   []byte
		want []byte
	}{
		{
			name: "Nil",
			in:   nil,
			want: []byte{},
		},
		{
			name: "Empty",
			in:   []byte{},
			want: []byte{},
		},
		{
			name: "0x00",
			in:   []byte{0x00},
			want: []byte{0x01},
		},
		{
			name: "0x5C",
			in:   []byte{0x5C},
			want: []byte{0x5D},
		},
		{
			name: "0xFE",
			in:   []byte{0xFE},
			want: []byte{0xFF, 0xFE},
		},
		{
			name: "0xFF",
			in:   []byte{0xFF},
			want: []byte{0xFF, 0xFF},
		},
		{
			name: "Multiple",
			in:   []byte{0x00, 0x5C, 0xFE, 0xFF},
			want: []byte{0x01, 0x5D, 0xFF, 0xFE, 0xFF, 0xFF},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := EscapeZero(tc.in)
			if !bytes.Equal(got, tc.want) {
				t.Errorf("expected %X, got %X", tc.want, got)
			}
		})
	}
}

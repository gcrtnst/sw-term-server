package vterm

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestScreenSetReflow(t *testing.T) {
	tt := []struct {
		inReflow   bool
		wantString string
	}{
		{
			inReflow:   false,
			wantString: "",
		},
		{
			inReflow:   true,
			wantString: "A",
		},
	}

	for _, tc := range tt {
		tc := tc

		name := "Disable"
		if tc.inReflow {
			name = "Enable"
		}

		t.Run(name, func(t *testing.T) {
			vt := New(5, 10)
			_ = vt.Output().Close()

			vt.Screen().SetReflow(tc.inReflow)
			_, _ = vt.Input().Write([]byte("AAAAAAAAAAAA"))
			vt.SetSize(5, 15)
			gotCell, _ := vt.Screen().Cell(Pos{Row: 0, Col: 10})
			gotString := string(gotCell.Runes)

			if gotString != tc.wantString {
				t.Errorf("expected %#v, got %#v", tc.wantString, gotString)
			}
		})
	}
}

func TestScreenSetAltScreen(t *testing.T) {
	tt := []struct {
		inAltScreen bool
		wantString  string
	}{
		{
			inAltScreen: false,
			wantString:  "A",
		},
		{
			inAltScreen: true,
			wantString:  "",
		},
	}

	for _, tc := range tt {
		tc := tc

		name := "Disable"
		if tc.inAltScreen {
			name = "Enable"
		}

		t.Run(name, func(t *testing.T) {
			vt := New(30, 120)
			_ = vt.Output().Close()

			vt.Screen().SetAltScreen(tc.inAltScreen)
			_, _ = vt.Input().Write([]byte("A\x1B[?1047h"))
			gotCell, _ := vt.Screen().Cell(Pos{Row: 0, Col: 0})
			gotString := string(gotCell.Runes)

			if gotString != tc.wantString {
				t.Errorf("expected %#v, got %#v", tc.wantString, gotString)
			}
		})
	}
}

func TestScreenSetDefaultColor(t *testing.T) {
	inFG := NewColorRGB(1, 2, 3)
	inBG := NewColorRGB(4, 5, 6)
	wantFG, wantBG := inFG, inBG
	wantFG.Type |= ColorDefaultFG
	wantBG.Type |= ColorDefaultBG

	vt := New(30, 120)
	_ = vt.Output().Close()
	scr := vt.Screen()
	scr.SetDefaultColor(inFG, inBG)
	cell, _ := scr.Cell(Pos{Row: 0, Col: 0})
	if !cell.FG.Equal(wantFG) {
		t.Errorf("fg: expected %#v, got %#v", wantFG, cell.FG)
	}
	if !cell.BG.Equal(wantBG) {
		t.Errorf("bg: expected %#v, got %#v", wantBG, cell.BG)
	}
}

func TestScreenSetPaletteColor(t *testing.T) {
	vt := New(30, 120)
	_ = vt.Output().Close()

	scr := vt.Screen()
	scr.SetPaletteColor(0, NewColorRGB(0, 0, 0))
	scr.SetPaletteColor(1, NewColorRGB(1, 1, 1))
	scr.SetPaletteColor(2, NewColorRGB(2, 2, 2))
	scr.SetPaletteColor(3, NewColorRGB(3, 3, 3))
	scr.SetPaletteColor(4, NewColorRGB(4, 4, 4))
	scr.SetPaletteColor(5, NewColorRGB(5, 5, 5))
	scr.SetPaletteColor(6, NewColorRGB(6, 6, 6))
	scr.SetPaletteColor(7, NewColorRGB(7, 7, 7))
	scr.SetPaletteColor(8, NewColorRGB(8, 8, 8))
	scr.SetPaletteColor(9, NewColorRGB(9, 9, 9))
	scr.SetPaletteColor(10, NewColorRGB(10, 10, 10))
	scr.SetPaletteColor(11, NewColorRGB(11, 11, 11))
	scr.SetPaletteColor(12, NewColorRGB(12, 12, 12))
	scr.SetPaletteColor(13, NewColorRGB(13, 13, 13))
	scr.SetPaletteColor(14, NewColorRGB(14, 14, 14))
	scr.SetPaletteColor(15, NewColorRGB(15, 15, 15))

	tt := []struct {
		inIdx   uint8
		wantCol Color
	}{
		{
			inIdx:   0,
			wantCol: NewColorRGB(0, 0, 0),
		},
		{
			inIdx:   1,
			wantCol: NewColorRGB(1, 1, 1),
		},
		{
			inIdx:   2,
			wantCol: NewColorRGB(2, 2, 2),
		},
		{
			inIdx:   3,
			wantCol: NewColorRGB(3, 3, 3),
		},
		{
			inIdx:   4,
			wantCol: NewColorRGB(4, 4, 4),
		},
		{
			inIdx:   5,
			wantCol: NewColorRGB(5, 5, 5),
		},
		{
			inIdx:   6,
			wantCol: NewColorRGB(6, 6, 6),
		},
		{
			inIdx:   7,
			wantCol: NewColorRGB(7, 7, 7),
		},
		{
			inIdx:   8,
			wantCol: NewColorRGB(8, 8, 8),
		},
		{
			inIdx:   9,
			wantCol: NewColorRGB(9, 9, 9),
		},
		{
			inIdx:   10,
			wantCol: NewColorRGB(10, 10, 10),
		},
		{
			inIdx:   11,
			wantCol: NewColorRGB(11, 11, 11),
		},
		{
			inIdx:   12,
			wantCol: NewColorRGB(12, 12, 12),
		},
		{
			inIdx:   13,
			wantCol: NewColorRGB(13, 13, 13),
		},
		{
			inIdx:   14,
			wantCol: NewColorRGB(14, 14, 14),
		},
		{
			inIdx:   15,
			wantCol: NewColorRGB(15, 15, 15),
		},
	}

	for _, tc := range tt {
		tc := tc
		name := fmt.Sprintf("No%d", tc.inIdx)
		t.Run(name, func(t *testing.T) {
			gotCol := scr.ConvertColorToRGB(NewColorIndexed(tc.inIdx))
			if !gotCol.Equal(tc.wantCol) {
				t.Errorf("expected %#v, got %#v", tc.wantCol, gotCol)
			}
		})
	}
}

func TestScreenCapture(t *testing.T) {
	fg := NewColorIndexed(7)
	fg.Type |= ColorDefaultFG
	bg := NewColorIndexed(0)
	bg.Type |= ColorDefaultBG

	vt := New(3, 4)
	_ = vt.Output().Close()
	vt.Screen().SetDefaultColor(fg, bg)
	_, _ = vt.Input().Write([]byte("123456789AB"))
	_, _ = vt.Input().Write([]byte("\x1B[2;3H")) // CUP
	_, _ = vt.Input().Write([]byte("\x1B[?25h")) // DECTCEM
	_, _ = vt.Input().Write([]byte("\x1B[6 q"))  // DECSCUSR
	got := vt.Screen().Capture()

	want := ScreenShot{
		Stride: 4,
		Cell: []Cell{
			{Runes: []rune("1"), Width: 1, FG: fg, BG: bg},
			{Runes: []rune("2"), Width: 1, FG: fg, BG: bg},
			{Runes: []rune("3"), Width: 1, FG: fg, BG: bg},
			{Runes: []rune("4"), Width: 1, FG: fg, BG: bg},
			{Runes: []rune("5"), Width: 1, FG: fg, BG: bg},
			{Runes: []rune("6"), Width: 1, FG: fg, BG: bg},
			{Runes: []rune("7"), Width: 1, FG: fg, BG: bg},
			{Runes: []rune("8"), Width: 1, FG: fg, BG: bg},
			{Runes: []rune("9"), Width: 1, FG: fg, BG: bg},
			{Runes: []rune("A"), Width: 1, FG: fg, BG: bg},
			{Runes: []rune("B"), Width: 1, FG: fg, BG: bg},
			{Runes: []rune{}, Width: 1, FG: fg, BG: bg},
		},
		CursorPos:     Pos{Row: 1, Col: 2},
		CursorVisible: true,
		CursorBlink:   false,
		CursorShape:   CursorShapeBarLeft,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %#v, got %#v", want, got)
	}
}

func TestScreenCaptureRGB(t *testing.T) {
	defaultFG := NewColorRGB(0, 0, 0)
	defaultBG := NewColorRGB(224, 224, 224)

	vt := New(4, 4)
	_ = vt.Output().Close()
	vt.Screen().SetDefaultColor(defaultFG, defaultBG)
	_, _ = vt.Input().Write([]byte("\x1B[30m0\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[31m1\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[32m2\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[33m3\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[34m4\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[35m5\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[36m6\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[37m7\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[40m0\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[41m1\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[42m2\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[43m3\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[44m4\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[45m5\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[46m6\x1B[0m"))
	_, _ = vt.Input().Write([]byte("\x1B[47m7\x1B[0m"))
	got := vt.Screen().CaptureRGB()

	want := ScreenShot{
		Stride: 4,
		Cell: []Cell{
			{Runes: []rune("0"), Width: 1, FG: NewColorRGB(0, 0, 0), BG: defaultBG},
			{Runes: []rune("1"), Width: 1, FG: NewColorRGB(224, 0, 0), BG: defaultBG},
			{Runes: []rune("2"), Width: 1, FG: NewColorRGB(0, 224, 0), BG: defaultBG},
			{Runes: []rune("3"), Width: 1, FG: NewColorRGB(224, 224, 0), BG: defaultBG},
			{Runes: []rune("4"), Width: 1, FG: NewColorRGB(0, 0, 224), BG: defaultBG},
			{Runes: []rune("5"), Width: 1, FG: NewColorRGB(224, 0, 224), BG: defaultBG},
			{Runes: []rune("6"), Width: 1, FG: NewColorRGB(0, 224, 224), BG: defaultBG},
			{Runes: []rune("7"), Width: 1, FG: NewColorRGB(224, 224, 224), BG: defaultBG},
			{Runes: []rune("0"), Width: 1, FG: defaultFG, BG: NewColorRGB(0, 0, 0)},
			{Runes: []rune("1"), Width: 1, FG: defaultFG, BG: NewColorRGB(224, 0, 0)},
			{Runes: []rune("2"), Width: 1, FG: defaultFG, BG: NewColorRGB(0, 224, 0)},
			{Runes: []rune("3"), Width: 1, FG: defaultFG, BG: NewColorRGB(224, 224, 0)},
			{Runes: []rune("4"), Width: 1, FG: defaultFG, BG: NewColorRGB(0, 0, 224)},
			{Runes: []rune("5"), Width: 1, FG: defaultFG, BG: NewColorRGB(224, 0, 224)},
			{Runes: []rune("6"), Width: 1, FG: defaultFG, BG: NewColorRGB(0, 224, 224)},
			{Runes: []rune("7"), Width: 1, FG: defaultFG, BG: NewColorRGB(224, 224, 224)},
		},
		CursorPos:     Pos{Row: 3, Col: 3},
		CursorVisible: true,
		CursorBlink:   true,
		CursorShape:   CursorShapeBlock,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %#v, got %#v", want, got)
	}
}

func TestScreenCell(t *testing.T) {
	s := &strings.Builder{}

	// Runes
	_, _ = s.WriteString("A")
	_, _ = s.WriteString("\u0061\u0302")
	_, _ = s.WriteString("\u0061\u0302\u0302\u0302\u0302\u0302")
	_, _ = s.WriteString("\r\n")

	// Width
	_, _ = s.WriteString("あ")
	_, _ = s.WriteString("\r\n")

	// CellAttrs
	_, _ = s.WriteString("\x1B[0mA")          // Normal
	_, _ = s.WriteString("\x1B[1mA\x1B[0m")   // Bold
	_, _ = s.WriteString("\x1B[4mA\x1B[0m")   // UnderlineSingle
	_, _ = s.WriteString("\x1B[21mA\x1B[0m")  // UnderlineDouble
	_, _ = s.WriteString("\x1B[4:3mA\x1B[0m") // UnderlineCurly
	_, _ = s.WriteString("\x1B[3mA\x1B[0m")   // Italic
	_, _ = s.WriteString("\x1B[5mA\x1B[0m")   // Blink
	_, _ = s.WriteString("\x1B[7mA\x1B[0m")   // Reverse
	_, _ = s.WriteString("\x1B[8mA\x1B[0m")   // Conceal
	_, _ = s.WriteString("\x1B[9mA\x1B[0m")   // Strike
	_, _ = s.WriteString("\x1B[14mA\x1B[0m")  // Font 4
	_, _ = s.WriteString("\x1B[19mA\x1B[0m")  // Font 9
	_, _ = s.WriteString("\x1B[73mA\x1B[0m")  // Small, BaselineRaise
	_, _ = s.WriteString("\x1B[74mA\x1B[0m")  // Small, BaselineLower
	_, _ = s.WriteString("\r\n")
	_, _ = s.WriteString("\x1B#6A\r\n") // DWL
	_, _ = s.WriteString("\x1B#3A\r\n") // DWL, DHLTop
	_, _ = s.WriteString("\x1B#4A\r\n") // DWL, DHLBottom

	// FG
	_, _ = s.WriteString("\x1B[34mA\x1B[0m")            // Idx 4
	_, _ = s.WriteString("\x1B[37mA\x1B[0m")            // Idx 7
	_, _ = s.WriteString("\x1B[38;5;127mA\x1B[0m")      // Idx 127
	_, _ = s.WriteString("\x1B[38;5;255mA\x1B[0m")      // Idx 255
	_, _ = s.WriteString("\x1B[38;2;10;20;30mA\x1B[0m") // RGB
	_, _ = s.WriteString("\r\n")

	// BG
	_, _ = s.WriteString("\x1B[44mA\x1B[0m")            // Idx 4
	_, _ = s.WriteString("\x1B[47mA\x1B[0m")            // Idx 7
	_, _ = s.WriteString("\x1B[48;5;127mA\x1B[0m")      // Idx 127
	_, _ = s.WriteString("\x1B[48;5;255mA\x1B[0m")      // Idx 255
	_, _ = s.WriteString("\x1B[48;2;10;20;30mA\x1B[0m") // RGB
	_, _ = s.WriteString("\r\n")

	vt := New(30, 120)
	_ = vt.Output().Close()
	vt.SetUTF8(true)
	vt.Screen().SetDefaultColor(NewColorIndexed(7), NewColorIndexed(0))
	_, _ = vt.Input().Write([]byte(s.String()))

	defaultFG := Color{
		Type: ColorIndexed | ColorDefaultFG,
		Idx:  7,
	}
	defaultBG := Color{
		Type: ColorIndexed | ColorDefaultBG,
		Idx:  0,
	}

	tt := []struct {
		inPos    Pos
		wantCell Cell
		wantOK   bool
	}{
		{
			inPos: Pos{Row: 0, Col: 0},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 0, Col: 1},
			wantCell: Cell{
				Runes: []rune("\u0061\u0302"),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 0, Col: 2},
			wantCell: Cell{
				Runes: []rune("\u0061\u0302\u0302\u0302\u0302\u0302"),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 0, Col: 3},
			wantCell: Cell{
				Runes: []rune(""),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 1, Col: 0},
			wantCell: Cell{
				Runes: []rune("あ"),
				Width: 2,
				Attrs: CellAttrs{},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 2, Col: 0},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 2, Col: 1},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{Bold: true},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 2, Col: 2},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{Underline: UnderlineSingle},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 2, Col: 3},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{Underline: UnderlineDouble},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 2, Col: 4},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{Underline: UnderlineCurly},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 2, Col: 5},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{Italic: true},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 2, Col: 6},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{Blink: true},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 2, Col: 7},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{Reverse: true},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 2, Col: 8},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{Conceal: true},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 2, Col: 9},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{Strike: true},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 2, Col: 10},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{Font: 4},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 2, Col: 11},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{Font: 9},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 2, Col: 12},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{Small: true, Baseline: BaselineRaise},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 2, Col: 13},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{Small: true, Baseline: BaselineLower},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 3, Col: 0},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{DWL: true},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 4, Col: 0},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{DWL: true, DHL: DHLTop},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 5, Col: 0},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{DWL: true, DHL: DHLBottom},
				FG:    defaultFG,
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 6, Col: 0},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    NewColorIndexed(4),
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 6, Col: 1},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    NewColorIndexed(7),
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 6, Col: 2},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    NewColorIndexed(127),
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 6, Col: 3},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    NewColorIndexed(255),
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 6, Col: 4},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    NewColorRGB(10, 20, 30),
				BG:    defaultBG,
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 7, Col: 0},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    defaultFG,
				BG:    NewColorIndexed(4),
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 7, Col: 1},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    defaultFG,
				BG:    NewColorIndexed(7),
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 7, Col: 2},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    defaultFG,
				BG:    NewColorIndexed(127),
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 7, Col: 3},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    defaultFG,
				BG:    NewColorIndexed(255),
			},
			wantOK: true,
		},
		{
			inPos: Pos{Row: 7, Col: 4},
			wantCell: Cell{
				Runes: []rune("A"),
				Width: 1,
				Attrs: CellAttrs{},
				FG:    defaultFG,
				BG:    NewColorRGB(10, 20, 30),
			},
			wantOK: true,
		},
		{
			inPos:    Pos{Row: -1, Col: -1},
			wantCell: Cell{},
			wantOK:   false,
		},
	}

	for _, tc := range tt {
		tc := tc
		name := fmt.Sprintf("R%dC%d", tc.inPos.Row, tc.inPos.Col)
		t.Run(name, func(t *testing.T) {
			gotCell, gotOK := vt.Screen().Cell(tc.inPos)

			if !reflect.DeepEqual(gotCell, tc.wantCell) {
				t.Errorf("cell: expected %#v, got %#v", tc.wantCell, gotCell)
			}
			if gotOK != tc.wantOK {
				t.Errorf("ok: expected %t, got %t", tc.wantOK, gotOK)
			}
		})
	}
}

func TestScreenCursorPos(t *testing.T) {
	vt := New(30, 120)
	_ = vt.Output().Close()
	in := vt.Input()
	scr := vt.Screen()

	wantPos := Pos{}
	gotPos := scr.CursorPos()
	if !reflect.DeepEqual(gotPos, wantPos) {
		t.Errorf("init: expected %#v, got %#v", wantPos, gotPos)
	}

	_, _ = in.Write([]byte("\x1B[10;20H"))
	wantPos = Pos{Row: 9, Col: 19}
	gotPos = scr.CursorPos()
	if !reflect.DeepEqual(gotPos, wantPos) {
		t.Errorf("set: expected %#v, got %#v", wantPos, gotPos)
	}
}

func TestScreenCursorVisible(t *testing.T) {
	vt := New(30, 120)
	_ = vt.Output().Close()
	in := vt.Input()
	scr := vt.Screen()

	want := true
	got := scr.CursorVisible()
	if got != want {
		t.Errorf("init: expected %t, got %t", want, got)
	}

	_, _ = in.Write([]byte("\x1B[?25l"))
	want = false
	got = scr.CursorVisible()
	if got != want {
		t.Errorf("set: expected %t, got %t", want, got)
	}
}

func TestScreenCursorBlink(t *testing.T) {
	vt := New(30, 120)
	_ = vt.Output().Close()
	in := vt.Input()
	scr := vt.Screen()

	want := true
	got := scr.CursorBlink()
	if got != want {
		t.Errorf("init: expected %t, got %t", want, got)
	}

	_, _ = in.Write([]byte("\x1B[?12l"))
	want = false
	got = scr.CursorBlink()
	if got != want {
		t.Errorf("set: expected %t, got %t", want, got)
	}
}

func TestScreenCursorShape(t *testing.T) {
	vt := New(30, 120)
	_ = vt.Output().Close()
	in := vt.Input()
	scr := vt.Screen()

	want := CursorShapeBlock
	got := scr.CursorShape()
	if got != want {
		t.Errorf("init: expected %d, got %d", want, got)
	}

	_, _ = in.Write([]byte("\x1B[3 q"))
	want = CursorShapeUnderline
	got = scr.CursorShape()
	if got != want {
		t.Errorf("set1: expected %d, got %d", want, got)
	}

	_, _ = in.Write([]byte("\x1B[5 q"))
	want = CursorShapeBarLeft
	got = scr.CursorShape()
	if got != want {
		t.Errorf("set2: expected %d, got %d", want, got)
	}
}

func TestScreenConvertColorToRGB(t *testing.T) {
	tt := []struct {
		name string
		in   Color
		want Color
	}{
		{
			name: "Idx0",
			in:   NewColorIndexed(0),
			want: NewColorRGB(0, 0, 0),
		},
		{
			name: "Idx1",
			in:   NewColorIndexed(1),
			want: NewColorRGB(224, 0, 0),
		},
		{
			name: "Idx2",
			in:   NewColorIndexed(2),
			want: NewColorRGB(0, 224, 0),
		},
		{
			name: "Idx3",
			in:   NewColorIndexed(3),
			want: NewColorRGB(224, 224, 0),
		},
		{
			name: "Idx4",
			in:   NewColorIndexed(4),
			want: NewColorRGB(0, 0, 224),
		},
		{
			name: "Idx5",
			in:   NewColorIndexed(5),
			want: NewColorRGB(224, 0, 224),
		},
		{
			name: "Idx6",
			in:   NewColorIndexed(6),
			want: NewColorRGB(0, 224, 224),
		},
		{
			name: "Idx7",
			in:   NewColorIndexed(7),
			want: NewColorRGB(224, 224, 224),
		},
	}

	vt := New(30, 120)
	_ = vt.Output().Close()
	scr := vt.Screen()

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := scr.ConvertColorToRGB(tc.in)
			if !got.Equal(tc.want) {
				t.Errorf("expected %#v, got %#v", tc.want, got)
			}
		})
	}
}

func TestScreenShotSize(t *testing.T) {
	tt := []struct {
		name             string
		inSS             ScreenShot
		wantRow, wantCol int
	}{
		{
			name:    "Zero",
			inSS:    ScreenShot{},
			wantRow: 0,
			wantCol: 0,
		},
		{
			name:    "ZeroStride",
			inSS:    ScreenShot{Cell: make([]Cell, 3)},
			wantRow: 0,
			wantCol: 0,
		},
		{
			name:    "ZeroLen",
			inSS:    ScreenShot{Stride: 2, Cell: make([]Cell, 0)},
			wantRow: 0,
			wantCol: 0,
		},
		{
			name:    "InvalidStride",
			inSS:    ScreenShot{Stride: 2, Cell: make([]Cell, 3)},
			wantRow: 0,
			wantCol: 0,
		},
		{
			name:    "Normal",
			inSS:    ScreenShot{Stride: 2, Cell: make([]Cell, 6)},
			wantRow: 3,
			wantCol: 2,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotRow, gotCol := tc.inSS.Size()
			if gotRow != tc.wantRow || gotCol != tc.wantCol {
				t.Errorf("expected (%d, %d), got (%d, %d)", tc.wantRow, tc.wantCol, gotRow, gotCol)
			}
		})
	}
}

func TestScreenShotAt(t *testing.T) {
	tt := []struct {
		name     string
		inSS     ScreenShot
		inPos    Pos
		wantCell Cell
	}{
		{
			name: "ErrorNegativeRow",
			inSS: ScreenShot{
				Stride: 1,
				Cell:   []Cell{{Width: 1}},
			},
			inPos:    Pos{Row: -1, Col: 0},
			wantCell: Cell{},
		},
		{
			name: "ErrorNegativeCol",
			inSS: ScreenShot{
				Stride: 1,
				Cell:   []Cell{{Width: 1}},
			},
			inPos:    Pos{Row: 0, Col: -1},
			wantCell: Cell{},
		},
		{
			name: "ErrorOverflowRow",
			inSS: ScreenShot{
				Stride: 2,
				Cell:   []Cell{{Width: 1}, {Width: 1}},
			},
			inPos:    Pos{Row: 1, Col: 0},
			wantCell: Cell{},
		},
		{
			name: "ErrorOverflowCol",
			inSS: ScreenShot{
				Stride: 1,
				Cell:   []Cell{{Width: 1}, {Width: 1}},
			},
			inPos:    Pos{Row: 0, Col: 1},
			wantCell: Cell{},
		},
		{
			name: "ErrorStride",
			inSS: ScreenShot{
				Stride: 0,
				Cell:   []Cell{{Width: 1}},
			},
			inPos:    Pos{Row: 0, Col: 0},
			wantCell: Cell{},
		},
		{
			name: "NormalZero",
			inSS: ScreenShot{
				Stride: 1,
				Cell:   []Cell{{Width: 1}},
			},
			inPos:    Pos{Row: 0, Col: 0},
			wantCell: Cell{Width: 1},
		},
		{
			name: "NormalRow",
			inSS: ScreenShot{
				Stride: 2,
				Cell:   []Cell{{Width: 1}, {Width: 1}},
			},
			inPos:    Pos{Row: 0, Col: 0},
			wantCell: Cell{Width: 1},
		},
		{
			name: "NormalCol",
			inSS: ScreenShot{
				Stride: 1,
				Cell:   []Cell{{Width: 1}, {Width: 1}},
			},
			inPos:    Pos{Row: 0, Col: 0},
			wantCell: Cell{Width: 1},
		},
		{
			name: "Normal",
			inSS: ScreenShot{
				Stride: 3,
				Cell: []Cell{
					{Width: 0}, {Width: 0}, {Width: 0},
					{Width: 0}, {Width: 0}, {Width: 1},
				},
			},
			inPos:    Pos{Row: 1, Col: 2},
			wantCell: Cell{Width: 1},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotCell := tc.inSS.At(tc.inPos)
			if !reflect.DeepEqual(gotCell, tc.wantCell) {
				t.Errorf("expected %#v, got %#v", tc.wantCell, gotCell)
			}
		})
	}
}

package main

import (
	"strconv"

	"github.com/gcrtnst/sw-term-server/internal/vterm"
)

func EncodeScreenShot(ss vterm.ScreenShot) [][]string {
	script := [][]string{{"#sw-term/screen"}}
	var cmd []string

	rows, cols := ss.Size()
	cmd = []string{
		"screen",
		strconv.Itoa(rows),
		strconv.Itoa(cols),
	}
	script = append(script, cmd)

	cmd = []string{
		"cursor",
		encodeBool(ss.CursorVisible),
		encodeBool(ss.CursorBlink),
		strconv.FormatUint(uint64(ss.CursorShape), 10),
		strconv.Itoa(int(ss.CursorPos.Row)),
		strconv.Itoa(int(ss.CursorPos.Col)),
	}
	script = append(script, cmd)

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			pos := vterm.Pos{Row: row, Col: col}
			cell := ss.At(pos)

			cmd := []string{
				"cell",
				strconv.Itoa(row),
				strconv.Itoa(col),
			}

			for _, sub := range encodeCell(cell) {
				cmd := append([]string{}, cmd...)
				cmd = append(cmd, sub...)
				script = append(script, cmd)
			}
		}
	}

	return script
}

func encodeCell(cell vterm.Cell) [][]string {
	script := [][]string{}

	if len(cell.Runes) > 0 && cell.Width > 0 {
		cmd := []string{
			"char",
			strconv.Itoa(cell.Width),
			string(cell.Runes),
		}
		script = append(script, cmd)
	}

	for _, sub := range encodeCellAttrs(cell.Attrs) {
		cmd := []string{"attr"}
		cmd = append(cmd, sub...)
		script = append(script, cmd)
	}

	for _, sub := range encodeCellColor(cell.FG, cell.BG) {
		cmd := []string{"color"}
		cmd = append(cmd, sub...)
		script = append(script, cmd)
	}

	return script
}

func encodeCellAttrs(attrs vterm.CellAttrs) [][]string {
	script := [][]string{}

	if attrs.Bold {
		cmd := []string{"bold"}
		script = append(script, cmd)
	}

	if attrs.Underline != vterm.UnderlineOff {
		cmd := []string{
			"underline",
			strconv.FormatUint(uint64(attrs.Underline), 10),
		}
		script = append(script, cmd)
	}

	if attrs.Italic {
		cmd := []string{"italic"}
		script = append(script, cmd)
	}

	if attrs.Blink {
		cmd := []string{"blink"}
		script = append(script, cmd)
	}

	if attrs.Reverse {
		cmd := []string{"reverse"}
		script = append(script, cmd)
	}

	if attrs.Conceal {
		cmd := []string{"conceal"}
		script = append(script, cmd)
	}

	if attrs.Strike {
		cmd := []string{"strike"}
		script = append(script, cmd)
	}

	if attrs.Font != 0 {
		cmd := []string{
			"font",
			strconv.Itoa(attrs.Font),
		}
		script = append(script, cmd)
	}

	if attrs.DWL {
		cmd := []string{"dwl"}
		script = append(script, cmd)
	}

	if attrs.DHL != vterm.DHLOff {
		cmd := []string{
			"dhl",
			strconv.FormatUint(uint64(attrs.DHL), 10),
		}
		script = append(script, cmd)
	}

	if attrs.Small {
		cmd := []string{"small"}
		script = append(script, cmd)
	}

	if attrs.Baseline != vterm.BaselineNormal {
		cmd := []string{
			"baseline",
			strconv.FormatUint(uint64(attrs.Baseline), 10),
		}
		script = append(script, cmd)
	}

	return script
}

func encodeCellColor(fg, bg vterm.Color) [][]string {
	script := [][]string{}

	if !fg.IsDefaultFG() {
		cmd := []string{"fg"}
		sub := encodeColor(fg)
		cmd = append(cmd, sub...)
		script = append(script, cmd)
	}

	if !bg.IsDefaultBG() {
		cmd := []string{"bg"}
		sub := encodeColor(bg)
		cmd = append(cmd, sub...)
		script = append(script, cmd)
	}

	return script
}

func encodeColor(col vterm.Color) []string {
	if col.IsRGB() {
		return []string{
			"rgb",
			strconv.FormatUint(uint64(col.Red), 10),
			strconv.FormatUint(uint64(col.Green), 10),
			strconv.FormatUint(uint64(col.Blue), 10),
		}
	}

	if col.IsIndexed() {
		return []string{
			"idx",
			strconv.FormatUint(uint64(col.Idx), 10),
		}
	}

	panic("invalid color type")
}

func encodeBool(b bool) string {
	if !b {
		return "0"
	}
	return "1"
}

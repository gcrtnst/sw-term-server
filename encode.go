package main

import (
	"bytes"
	"encoding/binary"

	"github.com/gcrtnst/sw-term-server/internal/vterm"
)

func EncodeScreenShot(ss vterm.ScreenShot) []byte {
	buf := new(bytes.Buffer)
	encodeScreenShot(buf, ss)
	return buf.Bytes()
}

func encodeScreenShot(buf *bytes.Buffer, ss vterm.ScreenShot) {
	encodeBool(buf, ss.CursorVisible)
	encodeBool(buf, ss.CursorBlink)
	_ = buf.WriteByte(byte(ss.CursorShape))
	encodeInt(buf, ss.CursorPos.Row)
	encodeInt(buf, ss.CursorPos.Col)

	rows, cols := ss.Size()
	encodeInt(buf, cols)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			pos := vterm.Pos{Row: row, Col: col}
			cell := ss.At(pos)

			encodeBool(buf, cell.Attrs.Bold)
			_ = buf.WriteByte(byte(cell.Attrs.Underline))
			encodeBool(buf, cell.Attrs.Italic)
			encodeBool(buf, cell.Attrs.Blink)
			encodeBool(buf, cell.Attrs.Reverse)
			encodeBool(buf, cell.Attrs.Conceal)
			encodeBool(buf, cell.Attrs.Strike)
			_ = buf.WriteByte(byte(cell.Attrs.Font))
			encodeBool(buf, cell.Attrs.DWL)
			_ = buf.WriteByte(byte(cell.Attrs.DHL))
			encodeBool(buf, cell.Attrs.Small)
			_ = buf.WriteByte(byte(cell.Attrs.Baseline))

			encodeColor(buf, cell.FG)
			encodeColor(buf, cell.BG)

			_ = buf.WriteByte(byte(cell.Width))
			encodeString(buf, string(cell.Runes))
		}
	}
}

func encodeColor(buf *bytes.Buffer, col vterm.Color) {
	_ = buf.WriteByte(byte(col.Type))

	var b [3]byte
	switch {
	case col.IsIndexed():
		b[0] = col.Idx
	case col.IsRGB():
		b[0] = col.Red
		b[1] = col.Green
		b[2] = col.Blue
	}
	_, _ = buf.Write(b[:])
}

func encodeString(buf *bytes.Buffer, s string) {
	encodeInt(buf, len(s))
	_, _ = buf.WriteString(s)
}

func encodeBool(buf *bytes.Buffer, b bool) {
	var c byte
	if b {
		c = 0x01
	}

	_ = buf.WriteByte(c)
}

func encodeInt(buf *bytes.Buffer, n int) {
	x := uint64(n)

	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], x)
	_, _ = buf.Write(b[:])
}

func EscapeZero(b []byte) []byte {
	buf := new(bytes.Buffer)
	for _, c := range b {
		switch c {
		case 0x00:
			_, _ = buf.WriteString(`\0`)
		case '\\':
			_, _ = buf.WriteString(`\\`)
		default:
			_ = buf.WriteByte(c)
		}
	}
	return buf.Bytes()
}

package vterm

import "testing"

func TestColorEqual(t *testing.T) {
	tt := []struct {
		name     string
		inA, inB Color
		want     bool
	}{
		{
			name: "EqRGB",
			inA:  Color{Type: ColorRGB, Red: 1, Green: 2, Blue: 3, Idx: 4},
			inB:  Color{Type: ColorRGB, Red: 1, Green: 2, Blue: 3, Idx: 5},
			want: true,
		},
		{
			name: "EqIndexed",
			inA:  Color{Type: ColorIndexed, Red: 1, Green: 2, Blue: 3, Idx: 4},
			inB:  Color{Type: ColorIndexed, Red: 5, Green: 6, Blue: 7, Idx: 4},
			want: true,
		},
		{
			name: "NeqType",
			inA:  Color{Type: ColorRGB, Red: 1, Green: 2, Blue: 3, Idx: 4},
			inB:  Color{Type: ColorIndexed, Red: 1, Green: 2, Blue: 3, Idx: 4},
			want: false,
		},
		{
			name: "NeqRGBRed",
			inA:  Color{Type: ColorRGB, Red: 1, Green: 2, Blue: 3, Idx: 4},
			inB:  Color{Type: ColorRGB, Red: 5, Green: 2, Blue: 3, Idx: 4},
			want: false,
		},
		{
			name: "NeqRGBGreen",
			inA:  Color{Type: ColorRGB, Red: 1, Green: 2, Blue: 3, Idx: 4},
			inB:  Color{Type: ColorRGB, Red: 1, Green: 5, Blue: 3, Idx: 4},
			want: false,
		},
		{
			name: "NeqRGBBlue",
			inA:  Color{Type: ColorRGB, Red: 1, Green: 2, Blue: 3, Idx: 4},
			inB:  Color{Type: ColorRGB, Red: 1, Green: 2, Blue: 5, Idx: 4},
			want: false,
		},
		{
			name: "NeqIndexed",
			inA:  Color{Type: ColorIndexed, Red: 1, Green: 2, Blue: 3, Idx: 4},
			inB:  Color{Type: ColorIndexed, Red: 1, Green: 2, Blue: 3, Idx: 5},
			want: false,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := tc.inA.Equal(tc.inB)
			if got != tc.want {
				t.Errorf("expected %t, got %t", tc.want, got)
			}
		})
	}
}

package vterm

import "testing"

func TestNew(t *testing.T) {
	tt := []struct {
		name     string
		inRows   int
		inCols   int
		wantRows int
		wantCols int
	}{
		{
			name:     "Normal",
			inRows:   30,
			inCols:   120,
			wantRows: 30,
			wantCols: 120,
		},
		{
			name:     "MinRows",
			inRows:   0,
			inCols:   120,
			wantRows: 1,
			wantCols: 120,
		},
		{
			name:     "MinCols",
			inRows:   30,
			inCols:   0,
			wantRows: 30,
			wantCols: 1,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			vt := New(tc.inRows, tc.inCols)
			gotRows, gotCols := vt.GetSize()

			if gotRows != tc.wantRows {
				t.Errorf("rows: expected %d, got %d", tc.wantRows, gotRows)
			}
			if gotCols != tc.wantCols {
				t.Errorf("cols: expected %d, got %d", tc.wantCols, gotCols)
			}
		})
	}
}

func TestVTermSize(t *testing.T) {
	tt := []struct {
		name     string
		inRows   int
		inCols   int
		wantRows int
		wantCols int
	}{
		{
			name:     "Normal",
			inRows:   30,
			inCols:   120,
			wantRows: 30,
			wantCols: 120,
		},
		{
			name:     "MinRows",
			inRows:   0,
			inCols:   120,
			wantRows: 1,
			wantCols: 120,
		},
		{
			name:     "MinCols",
			inRows:   30,
			inCols:   0,
			wantRows: 30,
			wantCols: 1,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			vt := New(60, 240)
			vt.SetSize(tc.inRows, tc.inCols)
			gotRows, gotCols := vt.GetSize()

			if gotRows != tc.wantRows {
				t.Errorf("rows: expected %d, got %d", tc.wantRows, gotRows)
			}
			if gotCols != tc.wantCols {
				t.Errorf("cols: expected %d, got %d", tc.wantCols, gotCols)
			}
		})
	}
}

func TestVTermUTF8(t *testing.T) {
	vt := New(30, 120)

	vt.SetUTF8(true)
	if !vt.GetUTF8() {
		t.FailNow()
	}

	vt.SetUTF8(false)
	if vt.GetUTF8() {
		t.FailNow()
	}

	vt.SetUTF8(true)
	if !vt.GetUTF8() {
		t.FailNow()
	}
}

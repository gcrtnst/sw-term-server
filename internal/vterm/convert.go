package vterm

// #include <limits.h>
// #include <stdbool.h>
// #include <stddef.h>
// #include <stdint.h>
import "C"
import "math"

const (
	MaxInt = C.INT_MAX
	MinInt = C.INT_MIN
)

func go2cBool(x bool) C.int {
	if x {
		return C.true
	}
	return C.false
}

func go2cInt(x int) (C.int, bool) {
	if math.MinInt <= C.INT_MIN {
		cmin := C.int(C.INT_MIN)
		gmin := int(cmin)
		if x <= gmin {
			return C.INT_MIN, false
		}
	}
	if math.MaxInt >= C.INT_MAX {
		cmax := C.int(C.INT_MAX)
		gmax := int(cmax)
		if x >= gmax {
			return C.INT_MAX, false
		}
	}
	return C.int(x), true
}

func c2goInt(x C.int) (int, bool) {
	if C.INT_MIN <= math.MinInt {
		gmin := math.MinInt
		cmin := C.int(gmin)
		if x <= cmin {
			return math.MinInt, false
		}
	}
	if C.INT_MAX >= math.MaxInt {
		gmax := math.MaxInt
		cmax := C.int(gmax)
		if x >= cmax {
			return math.MaxInt, false
		}
	}
	return int(x), true
}

func go2cSize(x int) (C.size_t, bool) {
	if x <= 0 {
		return 0, false
	}
	if math.MaxInt >= C.SIZE_MAX {
		cmax := C.size_t(C.SIZE_MAX)
		gmax := int(cmax)
		if x >= gmax {
			return C.SIZE_MAX, false
		}
	}
	return C.size_t(x), true
}

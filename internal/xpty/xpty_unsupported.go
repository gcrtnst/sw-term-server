//go:build !linux && !windows

package xpty

func open() (Terminal, error) {
	return nil, ErrUnsupported
}

//go:build !linux

package xpty

func open() (Terminal, error) {
	return nil, ErrUnsupported
}

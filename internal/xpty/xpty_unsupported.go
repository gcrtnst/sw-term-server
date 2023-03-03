//go:build !linux

package xpty

func open() (Master, Slave, error) {
	return nil, nil, ErrUnsupported
}

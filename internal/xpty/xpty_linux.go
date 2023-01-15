//go:build linux

package xpty

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"syscall"

	"golang.org/x/sys/unix"
)

func open() (Terminal, error) {
	ptm, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	name, err := ptsname(ptm)
	if err != nil {
		_ = ptm.Close()
		return nil, err
	}

	err = unlockpt(ptm)
	if err != nil {
		_ = ptm.Close()
		return nil, err
	}

	pts, err := os.OpenFile(name, os.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		_ = ptm.Close()
		return nil, err
	}

	t := &terminal{
		ptm: ptm,
		pts: pts,
	}
	return t, nil
}

type terminal struct {
	ptm, pts *os.File
}

func (t *terminal) Read(p []byte) (int, error) {
	return t.ptm.Read(p)
}

func (t *terminal) Write(p []byte) (int, error) {
	return t.ptm.Write(p)
}

func (t *terminal) Close() error {
	errPTM := t.ptm.Close()
	errPTS := t.pts.Close()

	if errPTM != nil {
		return errPTM
	}
	if errPTS != nil {
		return errPTS
	}
	return nil
}

func (t *terminal) Session(size Size) (Session, error) {
	s := &session{t: t}

	err := s.SetSize(size)
	if err != nil {
		return nil, err
	}

	return s, nil
}

type session struct {
	t *terminal
}

func (s *session) StartProcess(cmd Cmd) (*os.Process, error) {
	return os.StartProcess(cmd.Path, cmd.Args[:], &os.ProcAttr{
		Files: []*os.File{s.t.pts, s.t.pts, s.t.pts},
		Sys: &syscall.SysProcAttr{
			Setsid:  true,
			Setctty: true,
			Ctty:    0,
		},
	})
}

func (s *session) GetSize() (Size, error) {
	raw, errRaw := s.t.ptm.SyscallConn()
	if errRaw != nil {
		return Size{}, errRaw
	}

	var ws *unix.Winsize
	var errIoctl error
	errCtrl := raw.Control(func(fd uintptr) {
		ws, errIoctl = unix.IoctlGetWinsize(int(fd), unix.TIOCGWINSZ)
	})
	if errCtrl != nil {
		return Size{}, errCtrl
	}
	if errIoctl != nil {
		return Size{}, fmt.Errorf("ioctl TIOCGWINSZ: %w", errIoctl)
	}

	size := Size{
		Row: int(ws.Row),
		Col: int(ws.Col),
	}
	return size, nil
}

func (s *session) SetSize(size Size) error {
	ws, errWs := castWinsize(size)
	if errWs != nil {
		return errWs
	}

	raw, errRaw := s.t.ptm.SyscallConn()
	if errRaw != nil {
		return errRaw
	}

	var errIoctl error
	errCtrl := raw.Control(func(fd uintptr) {
		errIoctl = unix.IoctlSetWinsize(int(fd), unix.TIOCSWINSZ, &ws)
	})
	if errCtrl != nil {
		return errCtrl
	}
	if errIoctl != nil {
		return fmt.Errorf("ioctl TIOCSWINSZ: %w", errIoctl)
	}
	return nil
}

func (s *session) Close() error {
	return nil
}

func ptsname(f *os.File) (string, error) {
	raw, errRaw := f.SyscallConn()
	if errRaw != nil {
		return "", errRaw
	}

	var ptn uint32
	var errIoctl error
	errCtrl := raw.Control(func(fd uintptr) {
		ptn, errIoctl = unix.IoctlGetUint32(int(fd), unix.TIOCGPTN)
	})
	if errCtrl != nil {
		return "", errCtrl
	}
	if errIoctl != nil {
		return "", fmt.Errorf("ioctl TIOCGPTN: %w", errIoctl)
	}

	name := "/dev/pts/" + strconv.FormatUint(uint64(ptn), 10)
	return name, nil
}

func unlockpt(f *os.File) error {
	raw, errRaw := f.SyscallConn()
	if errRaw != nil {
		return errRaw
	}

	var errIoctl error
	errCtrl := raw.Control(func(fd uintptr) {
		errIoctl = unix.IoctlSetPointerInt(int(fd), unix.TIOCSPTLCK, 0)
	})
	if errCtrl != nil {
		return errCtrl
	}
	if errIoctl != nil {
		return fmt.Errorf("ioctl TIOCSPTLCK: %w", errIoctl)
	}
	return nil
}

func castWinsize(size Size) (unix.Winsize, error) {
	if size.Row <= 0 || math.MaxUint16 < size.Row || size.Col <= 0 || math.MaxUint16 < size.Col {
		return unix.Winsize{}, &SizeError{Size: size}
	}

	ws := unix.Winsize{Row: uint16(size.Row), Col: uint16(size.Col)}
	return ws, nil
}

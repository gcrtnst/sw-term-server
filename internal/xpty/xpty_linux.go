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

type terminal struct {
	m *os.File
	s *os.File
}

func (t *terminal) Read(p []byte) (int, error) {
	return t.m.Read(p)
}

func (t *terminal) Write(p []byte) (int, error) {
	return t.m.Write(p)
}

func (t *terminal) Close() error {
	errs := t.s.Close()
	errm := t.m.Close()

	if errm != nil {
		return errm
	}
	return errs
}

func (t *terminal) Start(cmd *Cmd) (*os.Process, error) {
	return os.StartProcess(cmd.Path, cmd.Args[:], &os.ProcAttr{
		Files: []*os.File{t.s, t.s, t.s},
		Sys: &syscall.SysProcAttr{
			Setsid:  true,
			Setctty: true,
			Ctty:    0,
		},
	})
}

func (t *terminal) GetSize() (int, int, error) {
	raw, err := t.m.SyscallConn()
	if err != nil {
		return 0, 0, err
	}

	var ws *unix.Winsize
	cerr := raw.Control(func(fd uintptr) {
		ws, err = unix.IoctlGetWinsize(int(fd), unix.TIOCGWINSZ)
	})
	if cerr != nil {
		return 0, 0, cerr
	}
	if err != nil {
		return 0, 0, fmt.Errorf("ioctl TIOCGWINSZ: %w", err)
	}

	return int(ws.Row), int(ws.Col), nil
}

func (t *terminal) SetSize(row, col int) error {
	if row < 0 || math.MaxUint16 < row || col < 0 || math.MaxUint16 < col {
		return &SizeError{Row: row, Col: col}
	}

	raw, err := t.m.SyscallConn()
	if err != nil {
		return err
	}

	ws := &unix.Winsize{
		Row: uint16(row),
		Col: uint16(col),
	}
	cerr := raw.Control(func(fd uintptr) {
		err = unix.IoctlSetWinsize(int(fd), unix.TIOCSWINSZ, ws)
	})
	if cerr != nil {
		return cerr
	}
	if err != nil {
		return fmt.Errorf("ioctl TIOCSWINSZ: %w", err)
	}

	return nil
}

func open() (t Terminal, err error) {
	ptm, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = ptm.Close()
		}
	}()

	ptsn, err := ptsname(ptm)
	if err != nil {
		return nil, err
	}

	err = unlockpt(ptm)
	if err != nil {
		return nil, err
	}

	pts, err := os.OpenFile(ptsn, os.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		return nil, err
	}

	t = &terminal{m: ptm, s: pts}
	return t, nil
}

func ptsname(ptm *os.File) (string, error) {
	raw, err := ptm.SyscallConn()
	if err != nil {
		return "", err
	}

	var ptyno uint32
	cerr := raw.Control(func(fd uintptr) {
		ptyno, err = unix.IoctlGetUint32(int(fd), unix.TIOCGPTN)
	})
	if cerr != nil {
		return "", cerr
	}
	if err != nil {
		return "", fmt.Errorf("ioctl TIOCGPTN: %w", err)
	}

	return "/dev/pts/" + strconv.FormatUint(uint64(ptyno), 10), nil
}

func unlockpt(ptm *os.File) error {
	raw, err := ptm.SyscallConn()
	if err != nil {
		return err
	}

	cerr := raw.Control(func(fd uintptr) {
		err = unix.IoctlSetPointerInt(int(fd), unix.TIOCSPTLCK, 0)
	})
	if cerr != nil {
		return cerr
	}
	if err != nil {
		return fmt.Errorf("ioctl TIOCSPTLCK: %w", err)
	}

	return nil
}

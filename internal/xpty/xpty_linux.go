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

type master struct {
	f *os.File
}

func (m *master) Read(p []byte) (int, error) {
	return m.f.Read(p)
}

func (m *master) Write(p []byte) (int, error) {
	return m.f.Write(p)
}

func (m *master) Close() error {
	return m.f.Close()
}

func (m *master) GetSize() (int, int, error) {
	raw, err := m.f.SyscallConn()
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

func (m *master) SetSize(row, col int) error {
	if row < 0 || math.MaxUint16 < row || col < 0 || math.MaxUint16 < col {
		return &SizeError{Row: row, Col: col}
	}

	raw, err := m.f.SyscallConn()
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

type slave struct {
	f *os.File
}

func (s *slave) Read(p []byte) (int, error) {
	return s.f.Read(p)
}

func (s *slave) Write(p []byte) (int, error) {
	return s.f.Write(p)
}

func (s *slave) Close() error {
	return s.f.Close()
}

func (s *slave) Start(cmd *Cmd) (*os.Process, error) {
	return os.StartProcess(cmd.Path, cmd.Args[:], &os.ProcAttr{
		Files: []*os.File{s.f, s.f, s.f},
		Sys: &syscall.SysProcAttr{
			Setsid:  true,
			Setctty: true,
			Ctty:    0,
		},
	})
}

func open() (m Master, s Slave, err error) {
	ptm, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			_ = ptm.Close()
		}
	}()

	ptsn, err := ptsname(ptm)
	if err != nil {
		return nil, nil, err
	}

	err = unlockpt(ptm)
	if err != nil {
		return nil, nil, err
	}

	pts, err := os.OpenFile(ptsn, os.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}

	m = &master{f: ptm}
	s = &slave{f: pts}
	return m, s, nil
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

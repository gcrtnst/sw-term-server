//go:build windows

package xpty

import (
	"errors"
	"fmt"
	"math"
	"os"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

func open() (Terminal, error) {
	var irh, iwh windows.Handle
	errI := windows.CreatePipe(&irh, &iwh, nil, 0)
	if errI != nil {
		return nil, fmt.Errorf("create conpty pipe: %w", errI)
	}

	var orh, owh windows.Handle
	errO := windows.CreatePipe(&orh, &owh, nil, 0)
	if errO != nil {
		_ = windows.CloseHandle(irh)
		_ = windows.CloseHandle(iwh)
		return nil, fmt.Errorf("create conpty pipe: %w", errO)
	}

	iwf := os.NewFile(uintptr(iwh), "")
	if iwf == nil {
		_ = windows.CloseHandle(irh)
		_ = windows.CloseHandle(orh)
		_ = windows.CloseHandle(owh)
		panic(errors.New("failed to create conpty pipe"))
	}

	orf := os.NewFile(uintptr(orh), "")
	if orf == nil {
		_ = windows.CloseHandle(irh)
		_ = iwf.Close()
		_ = windows.CloseHandle(owh)
		panic(errors.New("failed to create conpty pipe"))
	}

	t := &terminal{
		iw: iwf,
		or: orf,
		ok: true,
		ir: irh,
		ow: owh,
	}
	return t, nil
}

type terminal struct {
	iw *os.File
	or *os.File

	mu sync.RWMutex
	ok bool
	ir windows.Handle
	ow windows.Handle
}

func (t *terminal) Read(p []byte) (int, error) {
	return t.or.Read(p)
}

func (t *terminal) Write(p []byte) (int, error) {
	return t.iw.Write(p)
}

func (t *terminal) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.ok {
		return os.ErrClosed
	}
	t.ok = false

	errIW := t.iw.Close()
	errOR := t.or.Close()
	errIR := windows.CloseHandle(t.ir)
	errOW := windows.CloseHandle(t.ow)

	if errIW != nil {
		return errIW
	}
	if errOR != nil {
		return errOR
	}
	if errIR != nil {
		return fmt.Errorf("close conpty pipe: %w", errIR)
	}
	if errOW != nil {
		return fmt.Errorf("close conpty pipe: %w", errOW)
	}
	return nil
}

func (t *terminal) Session(size Size) (Session, error) {
	wsz, err := castWindowsCoord(size)
	if err != nil {
		return nil, err
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	if !t.ok {
		return nil, os.ErrClosed
	}

	var pc windows.Handle
	err = createPseudoConsole(wsz, t.ir, t.ow, 0, &pc)
	if err != nil {
		return nil, err
	}

	s := &session{
		ok: true,
		pc: pc,
		sz: size,
	}
	return s, nil
}

type session struct {
	mu sync.RWMutex
	ok bool
	pc windows.Handle
	sz Size
}

func (s *session) StartProcess(cmd Cmd) (*os.Process, error) {
	const PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE = 0x20016

	pathw, err := windows.UTF16PtrFromString(cmd.Path)
	if err != nil {
		return nil, fmt.Errorf("encode command path to UTF-16: %w", err)
	}

	argss := windows.ComposeCommandLine(cmd.Args)
	argsw, err := windows.UTF16PtrFromString(argss)
	if err != nil {
		return nil, fmt.Errorf("encode command arguments to UTF-16: %w", err)
	}

	al, err := windows.NewProcThreadAttributeList(1)
	if err != nil {
		return nil, fmt.Errorf("NewProcThreadAttributeList: %w", err)
	}
	defer al.Delete()

	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.ok {
		return nil, os.ErrClosed
	}

	err = al.Update(PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE, unsafeExternPointer(uintptr(s.pc)), unsafe.Sizeof(s.pc))
	if err != nil {
		return nil, fmt.Errorf("UpdateProcThreadAttribute: %w", err)
	}

	si := &windows.StartupInfoEx{}
	si.Cb = uint32(unsafe.Sizeof(*si))
	si.ProcThreadAttributeList = al.List()

	pi := &windows.ProcessInformation{}
	flags := uint32(windows.CREATE_DEFAULT_ERROR_MODE | windows.CREATE_UNICODE_ENVIRONMENT | windows.EXTENDED_STARTUPINFO_PRESENT)
	err = windows.CreateProcess(pathw, argsw, nil, nil, false, flags, nil, nil, &si.StartupInfo, pi)
	if err != nil {
		return nil, fmt.Errorf("CreateProcess: %w", err)
	}
	defer func() {
		err := windows.CloseHandle(pi.Process)
		if err != nil {
			panic(err)
		}
	}()
	defer func() {
		err := windows.CloseHandle(pi.Thread)
		if err != nil {
			panic(err)
		}
	}()

	return os.FindProcess(int(pi.ProcessId))
}

func (s *session) GetSize() (Size, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.ok {
		return Size{}, os.ErrClosed
	}
	return s.sz, nil
}

func (s *session) SetSize(size Size) error {
	wsz, err := castWindowsCoord(size)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.ok {
		return os.ErrClosed
	}

	err = resizePseudoConsole(s.pc, wsz)
	if err != nil {
		return err
	}

	s.sz = size
	return nil
}

func (s *session) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.ok {
		return os.ErrClosed
	}
	s.ok = false

	return closePseudoConsole(s.pc)
}

var (
	dllKernel32             = windows.NewLazySystemDLL("kernel32.dll")
	procCreatePseudoConsole = dllKernel32.NewProc("CreatePseudoConsole")
	procResizePseudoConsole = dllKernel32.NewProc("ResizePseudoConsole")
	procClosePseudoConsole  = dllKernel32.NewProc("ClosePseudoConsole")
)

func createPseudoConsole(size windows.Coord, hInput windows.Handle, hOutput windows.Handle, dwFlags uint32, phPC *windows.Handle) error {
	err := dllKernel32.Load()
	if err != nil {
		return ErrUnsupported
	}

	err = procCreatePseudoConsole.Find()
	if err != nil {
		return ErrUnsupported
	}

	result, _, _ := procCreatePseudoConsole.Call(uintptr(*(*uint32)(unsafe.Pointer(&size))), uintptr(hInput), uintptr(hOutput), uintptr(dwFlags), uintptr(unsafe.Pointer(phPC)))
	if result != uintptr(windows.S_OK) {
		err = hresult{n: int32(result)}
		return fmt.Errorf("CreatePseudoConsole: %w", err)
	}
	return nil
}

func resizePseudoConsole(hPC windows.Handle, size windows.Coord) error {
	err := dllKernel32.Load()
	if err != nil {
		return ErrUnsupported
	}

	err = procResizePseudoConsole.Find()
	if err != nil {
		return ErrUnsupported
	}

	result, _, _ := procResizePseudoConsole.Call(uintptr(hPC), uintptr(*(*uint32)(unsafe.Pointer(&size))))
	if result != uintptr(windows.S_OK) {
		err = hresult{n: int32(result)}
		return fmt.Errorf("ResizePseudoConsole: %w", err)
	}
	return nil
}

func closePseudoConsole(hPC windows.Handle) error {
	err := dllKernel32.Load()
	if err != nil {
		return ErrUnsupported
	}

	err = procClosePseudoConsole.Find()
	if err != nil {
		return ErrUnsupported
	}

	_, _, _ = procClosePseudoConsole.Call(uintptr(hPC))
	return nil
}

type hresult struct {
	n int32
}

func (h hresult) Error() string {
	return fmt.Sprintf("HRESULT(0x%08x)", uint32(h.n))
}

func castWindowsCoord(size Size) (windows.Coord, error) {
	if size.Row <= 0 || math.MaxInt16 < size.Row || size.Col <= 0 || math.MaxInt16 < size.Col {
		return windows.Coord{}, &SizeError{Size: size}
	}

	wsz := windows.Coord{X: int16(size.Col), Y: int16(size.Row)}
	return wsz, nil
}

func unsafeExternPointer(addr uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&addr))
}

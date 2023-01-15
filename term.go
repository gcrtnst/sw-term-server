package main

import (
	"errors"
	"io"
	"os"
	"sync"

	"github.com/gcrtnst/sw-term-server/internal/vterm"
	"github.com/gcrtnst/sw-term-server/internal/xpty"
)

type Term struct {
	pt xpty.Terminal
	ps xpty.Session
	vt *vterm.VTerm

	oc sync.Once
	di <-chan struct{}
	do <-chan struct{}
}

func NewTerm(cfg TermConfig) (*Term, error) {
	pt, err := cfg.Open()
	if err != nil {
		return nil, err
	}

	vt := vterm.New(cfg.Row, cfg.Col)
	vt.Screen().SetAltScreen(true)
	vt.Screen().SetReflow(false)

	// FIXME: Enabling UTF-8 support in the current libvterm (v0.3.2) may cause
	// crashes when displaying wide characters at the screen edge.
	//
	// The bug has been fixed in the libvterm bundled with Vim,
	// but not in the original version.
	// https://github.com/vim/vim/commit/8b89614e69b9b2330539d0482e44f4724053e780
	vt.SetUTF8(true)

	fg := vterm.NewColorIndexed(7)
	bg := vterm.NewColorIndexed(0)
	vt.Screen().SetDefaultColor(fg, bg)

	di := make(chan struct{})
	vi := vt.Input()
	go func() {
		_, err := io.Copy(vi, pt)
		if err != nil && !errors.Is(err, os.ErrClosed) {
			panic(err)
		}
		close(di)
	}()

	do := make(chan struct{})
	vo := vt.Output()
	go func() {
		_, err := io.Copy(pt, vo)
		if err != nil && !errors.Is(err, os.ErrClosed) {
			panic(err)
		}
		close(do)
	}()

	ps, err := pt.Session(xpty.Size{Row: cfg.Row, Col: cfg.Col})
	if err != nil {
		if err := pt.Close(); err != nil {
			panic(err)
		}
		if err := vo.Close(); err != nil {
			panic(err)
		}
		<-di
		<-do
		return nil, err
	}

	proc, err := ps.StartProcess(cfg.Cmd)
	if err != nil {
		if err := ps.Close(); err != nil {
			panic(err)
		}
		if err := pt.Close(); err != nil {
			panic(err)
		}
		if err := vo.Close(); err != nil {
			panic(err)
		}
		<-di
		<-do
		return nil, err
	}

	err = proc.Release()
	if err != nil {
		if err := ps.Close(); err != nil {
			panic(err)
		}
		if err := pt.Close(); err != nil {
			panic(err)
		}
		if err := vo.Close(); err != nil {
			panic(err)
		}
		<-di
		<-do
		panic(err)
	}

	t := &Term{
		pt: pt,
		ps: ps,
		vt: vt,
		di: di,
		do: do,
	}
	return t, nil
}

func (t *Term) Keyboard(key Key, mod vterm.Modifier) bool {
	if vk, ok := key.VTermKey(); ok {
		t.vt.KeyboardKey(vk, mod)
		return true
	}
	if r, ok := key.Rune(); ok {
		t.vt.KeyboardRune(r, mod)
		return true
	}
	return false
}

func (t *Term) Capture() vterm.ScreenShot {
	return t.vt.Screen().Capture()
}

func (t *Term) Close() error {
	t.oc.Do(t.close)
	return nil
}

func (t *Term) close() {
	var err error

	err = t.ps.Close()
	if err != nil {
		panic(err)
	}

	err = t.pt.Close()
	if err != nil {
		panic(err)
	}

	err = t.vt.Output().Close()
	if err != nil {
		panic(err)
	}

	<-t.di
	<-t.do
}

type TermConfig struct {
	Open     func() (xpty.Terminal, error)
	Row, Col int
	Cmd      xpty.Cmd
}

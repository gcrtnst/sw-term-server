package main

import (
	"errors"
	"sync"

	"github.com/gcrtnst/sw-term-server/internal/vterm"
)

var ErrInvalidKey = errors.New("invalid key")

type TermSlot struct {
	mu   sync.Mutex
	cfg  TermConfig
	term *Term
}

func NewTermSlot(cfg TermConfig) *TermSlot {
	return &TermSlot{cfg: cfg}
}

func (s *TermSlot) Keyboard(key Key, mod vterm.Modifier) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.start()
	if err != nil {
		return err
	}

	ok := s.term.Keyboard(key, mod)
	if !ok {
		return ErrInvalidKey
	}

	return nil
}

func (s *TermSlot) Capture() (vterm.ScreenShot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.start()
	if err != nil {
		return vterm.ScreenShot{}, err
	}

	ss := s.term.Capture()
	return ss, nil
}

func (s *TermSlot) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.term == nil {
		return
	}

	err := s.term.Close()
	if err != nil {
		panic(err)
	}

	s.term = nil
}

func (s *TermSlot) Close() error {
	s.Stop()
	return nil
}

func (s *TermSlot) start() error {
	if s.term != nil {
		return nil
	}

	term, err := NewTerm(s.cfg)
	if err != nil {
		return err
	}

	s.term = term
	return nil
}

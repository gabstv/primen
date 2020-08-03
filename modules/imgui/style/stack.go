package style

import "sync"

// Stack controls the positional variables stack like width, height
type Stack struct {
	l       sync.RWMutex
	last    stackState
	history []stackState
}

type stackState struct {
	MaxWidth  *float32
	MaxHeight *float32
	Width     *float32
	Height    *float32
}

func (s *Stack) MaxWidth() (val float32, ok bool) {
	s.l.RLock()
	defer s.l.RUnlock()
	if v := s.last.MaxWidth; v != nil {
		return *v, true
	}
	return 0, false
}

func (s *Stack) MaxWidthD(defaultv float32) float32 {
	s.l.RLock()
	defer s.l.RUnlock()
	if v := s.last.MaxWidth; v != nil {
		return *v
	}
	return defaultv
}

func (s *Stack) Width() (val float32, ok bool) {
	s.l.RLock()
	defer s.l.RUnlock()
	if v := s.last.Width; v != nil {
		return *v, true
	}
	return 0, false
}

func (s *Stack) PushMaxWidth(v float32) {
	s.l.Lock()
	defer s.l.Unlock()
	s.history = append(s.history, s.last)
	vcopy := v
	s.last.MaxWidth = &vcopy
}

func (s *Stack) PushWidth(v float32) {
	s.l.Lock()
	defer s.l.Unlock()
	s.history = append(s.history, s.last)
	vcopy := v
	s.last.Width = &vcopy
}

func (s *Stack) Pop() {
	s.l.Lock()
	defer s.l.Unlock()
	if len(s.history) < 1 {
		return
	}
	s.last = s.history[len(s.history)-1]
	s.history = s.history[:len(s.history)-1]
}

package main

import (
	"sync"

	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

type StateChange uint8

func (sc StateChange) Includes(other ...StateChange) bool {
	for _, x := range other {
		if sc&x == x {
			return true
		}
	}
	return false
}

const (
	NoChange StateChange = 1 << iota
	SelectedChanged
	HueChanged
	SaturationChanged
	ValueChanged
)

const AllChanged = SelectedChanged | HueChanged | SaturationChanged | ValueChanged

type subState struct {
	rgb        [3]int32
	hue        float64
	saturation float64
	value      float64
}

func (ss *subState) Selected() tcell.Color {
	return tcell.NewRGBColor(ss.rgb[0], ss.rgb[1], ss.rgb[2])
}

type State struct {
	pos     int
	sstates [7]*subState
	lock    sync.RWMutex
	pending StateChange
}

func NewDefaultState() *State {
	hue := 20.0
	s := NewState()
	for n := range s.sstates {
		r, g, b := noire.NewHSV(hue, 100, 100).RGB()
		s.sstates[n] = &subState{
			hue:        hue,
			saturation: 100,
			value:      100,
		}
		s.sstates[n].rgb[0] = int32(r)
		s.sstates[n].rgb[1] = int32(g)
		s.sstates[n].rgb[2] = int32(b)
		hue += 30
	}
	return s
}

func NewState() *State { return &State{pending: AllChanged} }

func (s *State) SubColors() []tcell.Color {
	a := make([]tcell.Color, len(s.sstates))
	for n := range s.sstates {
		a[n] = s.sstates[n].Selected()
	}
	return a
}

func (s *State) Pos() int              { return s.pos }
func (s *State) Len() int              { return len(s.sstates) }
func (s *State) Hue() float64          { return s.sstates[s.pos].hue }
func (s *State) Value() float64        { return s.sstates[s.pos].value }
func (s *State) Saturation() float64   { return s.sstates[s.pos].saturation }
func (s *State) Selected() tcell.Color { return s.sstates[s.pos].Selected() }

// BaseColor returns the current color at full saturation and brightness
func (s *State) BaseColor() *noire.Color { return noire.NewHSV(s.Hue(), 100, 100) }

// increment substate
func (s *State) Next() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.pos+1 >= len(s.sstates) {
		s.pos = 0
	} else {
		s.pos++
	}
	s.pending = AllChanged
}

// decrement substate
func (s *State) Prev() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.pos-1 < 0 {
		s.pos = len(s.sstates) - 1
	} else {
		s.pos--
	}
	s.pending = AllChanged
}

// Return StateChange since previous flush
func (s *State) Flush() StateChange {
	s.lock.Lock()
	defer s.lock.Unlock()
	a := s.pending
	s.pending = NoChange
	return a
}

func (s *State) SetSelected(r, g, b int32) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sstates[s.pos].rgb[0] = r
	s.sstates[s.pos].rgb[1] = g
	s.sstates[s.pos].rgb[2] = b
	s.pending = s.pending | SelectedChanged
}

func (s *State) SetHue(n float64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sstates[s.pos].hue = n
	s.pending = s.pending | HueChanged
}

func (s *State) SetSaturation(n float64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sstates[s.pos].saturation = n
	s.pending = s.pending | SaturationChanged
}

func (s *State) SetValue(n float64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sstates[s.pos].value = n
	s.pending = s.pending | ValueChanged
}

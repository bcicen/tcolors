package state

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bcicen/tcolors/logging"
	"github.com/gdamore/tcell"
	"github.com/olekukonko/tablewriter"
	"github.com/teacat/noire"
)

var (
	log          = logging.Init()
	malformedErr = fmt.Errorf("malformed state file")
)

const (
	subStateCount    = 7
	subStateByteSize = 56
	stateByteSize    = subStateCount*subStateByteSize + 8
)

type State struct {
	name    string
	path    string
	pos     int
	sstates [subStateCount]*subState // must be odd number for centering to work properly
	lock    sync.RWMutex
	pending Change
}

// Load attempts to read a stored state from disk. If no stored state exists or
// is readable, a default State and error will be returned.
func Load(path string) (*State, error) {
	s := NewDefault()
	s.path = path
	if err := s.load(); err != nil {
		return s, fmt.Errorf("failed to load state: %s", err)
	}
	return s, nil
}

// NewDefault returns a State initialized with default colors
func NewDefault() *State {
	hue := 20.0
	s := New()
	for n := range s.sstates {
		s.sstates[n] = &subState{noire.NewHSV(hue, 100, 100)}
		hue += 30
	}
	return s
}

func New() *State { return &State{pending: AllChanged} }

func (s *State) Save() error {
	if err := s.save(); err != nil {
		return fmt.Errorf("failed to save palette state: %s", err)
	}
	return nil
}

func (s *State) Name() string {
	if s.name == "" {
		return strings.ReplaceAll(filepath.Base(s.path), ".toml", "")
	}
	return s.name
}

func (s *State) SubColors() []tcell.Color {
	a := make([]tcell.Color, len(s.sstates))
	for n := range s.sstates {
		a[n] = s.sstates[n].TColor()
	}
	return a
}

func (s *State) Pos() int              { return s.pos }
func (s *State) Len() int              { return len(s.sstates) }
func (s *State) Selected() tcell.Color { return s.sstates[s.pos].TColor() }

func (s *State) Hue() float64 {
	hue, _, _ := s.sstates[s.pos].HSV()
	return hue
}

func (s *State) Saturation() float64 {
	_, sat, _ := s.sstates[s.pos].HSV()
	return sat
}

func (s *State) Value() float64 {
	_, _, val := s.sstates[s.pos].HSV()
	return val
}

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

// Return state Change since previous flush
func (s *State) Flush() Change {
	s.lock.Lock()
	defer s.lock.Unlock()
	a := s.pending
	s.pending = NoChange
	return a
}

func (s *State) SetHue(n float64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sstates[s.pos].SetHue(n)
	s.pending = s.pending | HueChanged
}

func (s *State) SetSaturation(n float64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sstates[s.pos].SetSaturation(n)
	s.pending = s.pending | SaturationChanged
}

func (s *State) SetValue(n float64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sstates[s.pos].SetValue(n)
	s.pending = s.pending | ValueChanged
}

// TableString returns an ascii table formatted representation of the current State
func (s *State) TableString() string {
	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"#", "Hex", "HSV", "RGB"})

	for n, ss := range s.sstates {
		table.Append([]string{
			fmt.Sprintf("%d", n),
			ss.HexString(),
			ss.HSVString(),
			ss.RGBString(),
		})
	}

	table.Render()
	return buf.String()
}

func (s *State) HexString() string {
	var txt []string
	for _, ss := range s.sstates {
		txt = append(txt, ss.HexString())
	}
	return strings.Join(txt, ", ")
}

func (s *State) HSVString() string {
	var txt []string
	for _, ss := range s.sstates {
		txt = append(txt, ss.HSVString())
	}
	return strings.Join(txt, ", ")
}

func (s *State) RGBString() string {
	var txt []string
	for _, ss := range s.sstates {
		txt = append(txt, ss.RGBString())
	}
	return strings.Join(txt, ", ")
}

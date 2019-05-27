package state

import (
	"bytes"
	"fmt"
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
	if err := s.load(path); err != nil {
		return s, fmt.Errorf("failed to load state: %s", err)
	}
	return s, nil
}

// NewDefault returns a State initialized with default colors
func NewDefault() *State {
	hue := 20.0
	s := New()
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

func New() *State { return &State{pending: AllChanged} }

func (s *State) Save(path string) error {
	if err := s.save(path); err != nil {
		return fmt.Errorf("failed to save palette state: %s", err)
	}
	return nil
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
func (s *State) Hue() float64          { return s.sstates[s.pos].hue }
func (s *State) Name() string          { return s.name }
func (s *State) Value() float64        { return s.sstates[s.pos].value }
func (s *State) Saturation() float64   { return s.sstates[s.pos].saturation }
func (s *State) Selected() tcell.Color { return s.sstates[s.pos].TColor() }

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

func (s *State) SetRGB(r, g, b int32) {
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

// TableString returns an ascii table formatted representation of the current State
func (s *State) TableString() string {
	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"#", "Hex", "HSV", "RGB"})

	for n, ss := range s.sstates {
		hex := fmt.Sprintf("%06x", ss.TColor().Hex())
		hsv := fmt.Sprintf("%03.0f %03.0f %03.0f", ss.hue, ss.saturation, ss.value)
		rgb := fmt.Sprintf("%03d %03d %03d", ss.rgb[0], ss.rgb[1], ss.rgb[2])
		table.Append([]string{fmt.Sprintf("%d", n), hex, hsv, rgb})
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

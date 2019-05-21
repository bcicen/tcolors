package state

import (
	"fmt"
	"os"
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
	pos     int
	sstates [subStateCount]*subState // must be odd number for centering to work properly
	lock    sync.RWMutex
	pending Change
}

// Load attempts to read a stored state from disk. If no stored state exists or
// is readable, a default State and error will be returned.
func Load() (*State, error) {
	s := NewDefault()
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

func (s *State) load() error {
	var buf [stateByteSize]byte
	var offset int

	path, err := statePath()
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	n, err := f.Read(buf[:])
	if err != nil {
		return err
	}
	if n != stateByteSize {
		return malformedErr
	}

	sscount := int(int32FromBytes(buf[offset : offset+4]))
	offset += 4
	s.pos = int(int32FromBytes(buf[offset : offset+4]))
	offset += 4

	for i := 0; i < sscount; i++ {
		if i > subStateCount {
			break
		}
		s.sstates[i].load(buf[offset : offset+subStateByteSize])
		offset += subStateByteSize
		log.Debugf("loaded substate [%d] from %s", i, path)
	}

	log.Infof("loaded state [%s]", path)
	return nil
}

func (s *State) Save() error {
	path, err := statePath()
	if err != nil {
		return err
	}
	log.Infof("saving state [%s]", path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(s.bytes())
	return err
}

// Bytes returns a byte-serialized representation
// of the current state
func (s *State) bytes() []byte {
	var buf [stateByteSize]byte
	var offset int
	offset += writeInt32Bytes(buf[offset:], int32(len(s.sstates)))
	offset += writeInt32Bytes(buf[offset:], int32(s.pos))
	for _, ss := range s.sstates {
		offset += copy(buf[offset:], ss.bytes())
	}
	return buf[:]
}

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

// OutputTable prints a table-formatted representation of the current State
func (s *State) OutputTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Hex", "HSV", "RGB"})

	for n, ss := range s.sstates {
		hex := fmt.Sprintf("%06x", ss.Selected().Hex())
		hsv := fmt.Sprintf("%03.0f %03.0f %03.0f", ss.hue, ss.saturation, ss.value)
		rgb := fmt.Sprintf("%03d %03d %03d", ss.rgb[0], ss.rgb[1], ss.rgb[2])
		table.Append([]string{fmt.Sprintf("%d", n), hex, hsv, rgb})
	}

	table.Render()
}

func (s *State) OutputHex() string {
	var txt []string
	for _, ss := range s.sstates {
		txt = append(txt, fmt.Sprintf("%06x", ss.Selected().Hex()))
	}
	return strings.Join(txt, ", ")
}

func (s *State) OutputHSV() string {
	var txt []string
	for _, ss := range s.sstates {
		txt = append(txt, fmt.Sprintf("%03.0f %03.0f %03.0f", ss.hue, ss.saturation, ss.value))
	}
	return strings.Join(txt, ", ")
}

func (s *State) OutputRGB() string {
	var txt []string
	for _, ss := range s.sstates {
		txt = append(txt, fmt.Sprintf("%03d %03d %03d", ss.rgb[0], ss.rgb[1], ss.rgb[2]))
	}
	return strings.Join(txt, ", ")
}

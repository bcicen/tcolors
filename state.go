package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"sync"

	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

var malformedErr = fmt.Errorf("malformed state file")

const subStateCount = 7
const subStateByteSize = 56
const stateByteSize = subStateCount*subStateByteSize + 8

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

func (ss *subState) load(b []byte) error {
	var offset int

	if len(b) != subStateByteSize {
		return malformedErr
	}

	ss.rgb[0] = int32FromBytes(b[offset : offset+4])
	offset += 4
	ss.rgb[1] = int32FromBytes(b[offset : offset+4])
	offset += 4
	ss.rgb[2] = int32FromBytes(b[offset : offset+4])
	offset += 4

	ss.hue = float64FromBytes(b[offset : offset+8])
	offset += 8
	ss.saturation = float64FromBytes(b[offset : offset+8])
	offset += 8
	ss.value = float64FromBytes(b[offset : offset+8])
	offset += 8

	return nil
}

func (ss *subState) bytes() []byte {
	var buf [subStateByteSize]byte
	var offset int
	for _, n := range ss.rgb {
		offset += writeInt32Bytes(buf[offset:], n)
	}
	offset += writeFloat64Bytes(buf[offset:], ss.hue)
	offset += writeFloat64Bytes(buf[offset:], ss.saturation)
	offset += writeFloat64Bytes(buf[offset:], ss.value)
	return buf[:]
}

func (ss *subState) Selected() tcell.Color {
	return tcell.NewRGBColor(ss.rgb[0], ss.rgb[1], ss.rgb[2])
}

type State struct {
	pos     int
	sstates [subStateCount]*subState // must be odd number for centering to work properly
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

func (s *State) Load() error {
	var buf [stateByteSize]byte
	var offset int

	path := fmt.Sprintf("%s/state.bin", configPath())
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
	}

	log.Infof("loaded state from %s", path)
	return nil
}

func (s *State) Save() {
	path := fmt.Sprintf("%s/state.bin", configPath())
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Write(s.bytes())
	if err != nil {
		panic(err)
	}
	log.Infof("saved state to %s", path)
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

func writeInt32Bytes(buf []byte, n int32) int {
	binary.BigEndian.PutUint32(buf, uint32(n))
	return 4
}

func int32FromBytes(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

func float64FromBytes(buf []byte) float64 {
	bits := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(bits)
}

func writeFloat64Bytes(buf []byte, float float64) int {
	bits := math.Float64bits(float)
	binary.LittleEndian.PutUint64(buf[:], bits)
	return 8
}

func configPath() string {
	return "."
}

package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

const (
	satMin   = 0.5
	satMax   = 100.0
	satIncr  = 0.5
	satCount = int(satMax/satIncr) + 1
)

type SaturationBar struct {
	items  [satCount]tcell.Color // navigation colors
	scale  [satCount]float64
	pos    int
	offset int
	width  int
	pst    tcell.Style // pointer style
	state  *State
}

func NewSaturationBar(s *State) *SaturationBar {
	bar := &SaturationBar{state: s}

	i := satMin
	for n, _ := range bar.scale {
		bar.scale[n] = i
		i += 0.5
	}

	return bar
}

// Draw redraws bar at given coordinates and screen, returning the number
// of rows occupied
func (bar *SaturationBar) Draw(x, y int, s tcell.Screen) int {
	var st tcell.Style

	n := bar.offset
	col := 0

	// border bars
	s.SetCell(x-1, y+1, bar.pst, '│')
	s.SetCell(x-1, y+2, bar.pst, '│')
	s.SetCell(bar.width+x+1, y+1, bar.pst, '│')
	s.SetCell(bar.width+x+1, y+2, bar.pst, '│')

	for col <= bar.width && n < len(bar.items) {
		st = st.Background(bar.items[n])
		s.SetCell(col+x, y, blkSt, '█')
		s.SetCell(col+x, y+1, st, ' ')
		s.SetCell(col+x, y+2, st, ' ')

		col++
		n++
	}

	ix := (bar.pos - bar.offset) + x
	s.SetCell(ix, y, bar.pst, '▾')

	s.SetCell(bar.width/2, y+3, bar.pst, []rune(fmt.Sprintf("%5.1f  ", bar.Value()-0.5))...)

	return 4
}

func (bar *SaturationBar) Value() float64 { return bar.scale[bar.pos] }
func (bar *SaturationBar) SetPos(n float64) {
	var idx int
	for idx < len(bar.scale)-1 {
		if bar.scale[idx+1] > n {
			break
		}
		idx++
	}

	switch {
	case idx > bar.pos:
		bar.up(idx - bar.pos)
	case idx < bar.pos:
		bar.down(bar.pos - idx)
	}
}

func (bar *SaturationBar) Resize(w int) {
	bar.width = w
	bar.up(0)
	bar.down(0)
}

// State change handler
func (bar *SaturationBar) Handle(change StateChange) {
	var nc *noire.Color

	if change.Includes(HueChanged, ValueChanged) {
		nc = bar.state.BaseColor()

		for n, val := range bar.scale {
			nc = noire.NewHSV(bar.state.Hue(), val, bar.state.Value())
			bar.items[n] = toTColor(nc)
		}
	}

	if change.Includes(SaturationChanged) {
		bar.SetPos(bar.state.Saturation())
	}
}

func (bar *SaturationBar) setState() {
	bar.state.SetSaturation(bar.scale[bar.pos])
}

func (bar *SaturationBar) Width() int { return bar.width }

func (bar *SaturationBar) Up(step int) {
	bar.up(step)
	bar.setState()
}

func (bar *SaturationBar) Down(step int) {
	bar.down(step)
	bar.setState()
}

func (bar *SaturationBar) up(step int) {
	max := len(bar.items) - 1
	maxOffset := max - bar.width
	switch {
	case step <= 0:
	case bar.pos == max:
		return
	case bar.pos+step > max:
		bar.pos = max
	default:
		log.Debugf("pos=%d", bar.pos)
		bar.pos += step
		log.Debugf("pos=%d", bar.pos)
	}

	if (bar.pos - bar.offset) > bar.width-scrollAhead {
		bar.offset = (bar.pos - bar.width) + scrollAhead
	}
	if bar.offset >= maxOffset {
		bar.offset = maxOffset
	}
}

func (bar *SaturationBar) down(step int) {
	log.Debugf("DOWN=%d", step)
	switch {
	case step <= 0:
	case bar.pos == 0:
		return
	case bar.pos-step < 0:
		bar.pos = 0
	default:
		bar.pos -= step
	}

	if bar.pos-bar.offset < scrollAhead {
		bar.offset = bar.pos - scrollAhead
	}
	if bar.offset < 0 {
		bar.offset = 0
	}
}

func (bar *SaturationBar) SetPointerStyle(st tcell.Style) { bar.pst = st }

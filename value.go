package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

const scrollAhead = 3

type ValueBar struct {
	items  []tcell.Color // navigation colors
	scale  []float64
	pos    int
	offset int
	width  int
	pst    tcell.Style // pointer style
	state  *State
}

func NewValueBar(s *State) *ValueBar {
	bar := &ValueBar{state: s}
	for i := 0.0; i < 100.1; i += 0.5 {
		bar.scale = append(bar.scale, i)
	}
	bar.items = make([]tcell.Color, len(bar.scale))
	return bar
}

// Draw redraws bar at given coordinates and screen, returning the number
// of rows occupied
func (bar *ValueBar) Draw(x, y int, s tcell.Screen) int {
	var st tcell.Style

	n := bar.offset
	col := 0
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

	s.SetCell(bar.width/2, y+3, bar.pst, []rune(fmt.Sprintf("%5.1f  ", bar.Value()))...)

	return 4
}

func (bar *ValueBar) Value() float64 { return bar.scale[bar.pos] }
func (bar *ValueBar) SetValue(n float64) {
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

func (bar *ValueBar) Resize(w int) {
	bar.width = w
	bar.up(0)
	bar.down(0)
}

// State change handler
func (bar *ValueBar) Handle(change StateChange) {
	var nc *noire.Color

	if change.Includes(HueChanged, SaturationChanged) {
		for n, val := range bar.scale {
			nc = noire.NewHSV(bar.state.Hue(), bar.state.Saturation(), val)
			bar.items[n] = toTColor(nc)
		}
	}

	if change.Includes(ValueChanged) {
		bar.SetValue(bar.state.Value())
	}

}

func (bar *ValueBar) Width() int { return bar.width }

func (bar *ValueBar) Up(step int) {
	bar.up(step)
	bar.setState()
}

func (bar *ValueBar) Down(step int) {
	bar.down(step)
	bar.setState()
}

func (bar *ValueBar) SetPointerStyle(st tcell.Style) { bar.pst = st }

func (bar *ValueBar) up(step int) {
	max := len(bar.items) - 1
	maxOffset := max - bar.width
	switch {
	case step <= 0:
	case bar.pos == max:
		return
	case bar.pos+step > max:
		bar.pos = max
	default:
		bar.pos += step
	}

	if (bar.pos - bar.offset) > bar.width-scrollAhead {
		bar.offset = (bar.pos - bar.width) + scrollAhead
	}
	if bar.offset >= maxOffset {
		bar.offset = maxOffset
	}
}

func (bar *ValueBar) down(step int) {
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

func (bar *ValueBar) setState() {
	bar.state.SetValue(bar.scale[bar.pos])
}

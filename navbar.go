package main

import (
	"math"

	"github.com/bcicen/tcolors/state"
	"github.com/gdamore/tcell"
)

const scrollAhead = 3

type NavBar struct {
	items  []tcell.Color // navigation colors
	label  string
	pos    int
	offset int
	width  int
	height int
	pst    tcell.Style // pointer style
	state  *state.State
}

func NewNavBar(s *state.State, length int) *NavBar {
	return &NavBar{
		items: make([]tcell.Color, length),
		state: s,
	}
}

// Draw redraws bar at given coordinates and screen, returning the number
// of rows occupied
func (bar *NavBar) Draw(x, y int, s tcell.Screen) int {
	var st tcell.Style

	n := bar.offset
	col := 0

	// border bars
	for i := 1; i <= bar.height; i++ {
		s.SetCell(x-1, y+i, bar.pst, '│')
		s.SetCell(bar.width+x, y+i, bar.pst, '│')
	}

	for col < bar.width && n < len(bar.items) {
		st = st.Background(bar.items[n])
		s.SetCell(col+x, y, blkSt, '█')
		for i := 1; i <= bar.height; i++ {
			s.SetCell(col+x, y+i, st, ' ')
		}

		col++
		n++
	}

	ix := (bar.pos - bar.offset) + x
	s.SetCell(ix, y, bar.pst, '▾')

	labelX := x + ((bar.width - 4) / 2)
	for n, ch := range []rune(bar.label) {
		s.SetCell(labelX+n, y+4, bar.pst, ch)
	}

	return bar.height + 1
}

func (bar *NavBar) SetLabel(s string) { bar.label = s }

func (bar *NavBar) SetPos(idx int) {
	switch {
	case idx > bar.pos:
		bar.up(idx - bar.pos)
	case idx < bar.pos:
		bar.down(bar.pos - idx)
	}
}

func (bar *NavBar) Resize(w, h int) {
	bar.height = barHeight(h)
	bar.width = w
	bar.up(0)
	bar.down(0)
}

// NavBar implements Section
func (bar *NavBar) Up(int)                         {}
func (bar *NavBar) Down(int)                       {}
func (bar *NavBar) Handle(state.Change)            {}
func (bar *NavBar) SetPointerStyle(st tcell.Style) { bar.pst = st }

func (bar *NavBar) up(step int) {
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

func (bar *NavBar) down(step int) {
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

func roundFloat(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

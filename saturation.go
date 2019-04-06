package main

import (
	"sync"

	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

type SaturationBar struct {
	items  []tcell.Color // navigation colors
	scale  []float64
	pos    int
	offset int
	width  int
	pst    tcell.Style // pointer style
	lock   sync.RWMutex
}

func NewSaturationBar(width int) *SaturationBar {
	bar := &SaturationBar{width: width}
	for i := -1.0; i < 1.01; i += 0.005 {
		bar.scale = append(bar.scale, i)
	}
	bar.items = make([]tcell.Color, len(bar.scale))
	return bar
}

// Draw redraws bar at given coordinates and screen, returning the number
// of rows occupied
func (bar *SaturationBar) Draw(x, y int, s tcell.Screen) int {
	var st tcell.Style

	n := bar.offset
	col := 0
	for col <= bar.width && n < len(bar.items) {
		st = st.Background(bar.items[n])
		s.SetCell(col+x, y+1, st, ' ')
		s.SetCell(col+x, y+2, st, ' ')

		if n == bar.pos {
			s.SetCell(col+x, y, bar.pst, '▼')
		} else {
			s.SetCell(col+x, y, blkSt, '▼')
		}
		col++
		n++
	}

	return 4
}

func (bar *SaturationBar) Value() float64 { return bar.scale[bar.pos] }
func (bar *SaturationBar) SetPos(n int)   { bar.pos = n }

func (bar *SaturationBar) Resize(w int) {
	bar.width = w
	bar.Up(1)
	bar.Down(1)
}

func (bar *SaturationBar) Update(base *noire.Color) {
	bar.lock.Lock()
	defer bar.lock.Unlock()

	for n, val := range bar.scale {
		bar.items[n] = toTColor(applySaturation(val, base))
	}
}

func (bar *SaturationBar) Up(step int) {
	bar.lock.Lock()
	defer bar.lock.Unlock()

	max := len(bar.items) - 1
	switch {
	case bar.pos == max:
		return
	case bar.pos+step > max:
		bar.pos = max
	default:
		bar.pos += step
	}

	if (bar.pos-bar.offset) > bar.width-scrollAhead && bar.pos < len(bar.items)-scrollAhead {
		bar.offset = (bar.pos - bar.width) + scrollAhead
	}
}

func (bar *SaturationBar) Down(step int) {
	bar.lock.Lock()
	defer bar.lock.Unlock()

	switch {
	case bar.pos == 0:
		return
	case bar.pos-step < 0:
		bar.pos = 0
	default:
		bar.pos -= step
	}

	if bar.pos-bar.offset < scrollAhead && bar.pos > scrollAhead {
		bar.offset = bar.pos - scrollAhead
	}
}

func (bar *SaturationBar) SetPointerStyle(st tcell.Style) { bar.pst = st }

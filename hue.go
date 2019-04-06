package main

import (
	"sync"

	"github.com/gdamore/tcell"
)

type HueBar struct {
	items  []tcell.Color // navigation colors
	mItems []int         // minimap sample indices
	pos    int
	width  int
	pst    tcell.Style // pointer style
	lock   sync.RWMutex
}

func NewHueBar(width int) *HueBar {
	return &HueBar{width: width}
}

func (bar *HueBar) SetPos(n int)          { bar.pos = n }
func (bar *HueBar) Selected() tcell.Color { return bar.items[bar.pos] }
func (bar *HueBar) center() int           { return (bar.width / 2) }

// Draw redraws bar at given coordinates and screen, returning the number
// of rows occupied
func (bar *HueBar) Draw(x, y int, s tcell.Screen) int {
	var st tcell.Style
	center := bar.width / 2

	s.SetCell(center+x, y, bar.pst, '▼')

	for col, color := range bar.Items() {
		st = st.Background(color)
		s.SetCell(col+x, y+1, st, ' ')
		s.SetCell(col+x, y+2, st, ' ')
	}

	boxPad := bar.width / 30
	if boxPad < 2 {
		boxPad = 2
	}

	s.SetCell(center+x-boxPad, y+3, indicatorSt, '┌')
	s.SetCell(center+x+boxPad, y+3, indicatorSt, '┐')

	for col, color := range bar.MiniMap() {
		st = st.Background(color)
		s.SetCell(col+x, y+4, st, ' ')
	}

	s.SetCell(center+x-boxPad, y+5, indicatorSt, '└')
	s.SetCell(center+x+boxPad, y+5, indicatorSt, '┘')

	return 6
}

func (bar *HueBar) miniStep() int {
	n := len(bar.items) / bar.width
	if n > 13 {
		n = 13
	}
	return n
}

func (bar *HueBar) Resize(w int) {
	bar.width = w
	bar.Update(bar.items)
}

func (bar *HueBar) Update(a []tcell.Color) {
	bar.lock.Lock()
	defer bar.lock.Unlock()
	bar.items = append(bar.items[:0], a...)

	// build minimap indices
	bar.mItems = bar.mItems[0:]
	miniStep := bar.miniStep()
	for n := 0; n < len(bar.items); n += miniStep {
		bar.mItems = append(bar.mItems, n)
	}
}

func (bar *HueBar) Items() []tcell.Color {
	l := bar.pos - bar.center()
	r := bar.pos + bar.center() + 1
	ilen := len(bar.items)
	if l < 0 {
		return append(bar.items[ilen+l:], bar.items[0:r]...)
	}
	if r > ilen {
		return append(bar.items[l:ilen-1], bar.items[0:r-ilen]...)
	}
	return bar.items[l:r]
}

func (bar *HueBar) MiniMap() []tcell.Color {
	var mpos int
	for mpos < len(bar.mItems)-1 {
		if bar.mItems[mpos+1] >= bar.pos {
			break
		}
		mpos++
	}

	l := mpos - bar.center()
	r := mpos + bar.center() + 1
	ilen := len(bar.mItems)

	var a []int
	switch {
	case l < 0:
		a = append(bar.mItems[ilen+l:], bar.mItems[0:r]...)
	case r > ilen:
		a = append(bar.mItems[l:], bar.mItems[0:r-ilen]...)
	default:
		a = bar.mItems[l:r]
	}

	var colors []tcell.Color
	for _, idx := range a {
		colors = append(colors, bar.items[idx])
	}

	return colors
}

func (bar *HueBar) Up(step int) {
	bar.lock.Lock()
	defer bar.lock.Unlock()
	bar.pos += step
	if bar.pos >= len(bar.items)-1 {
		bar.pos -= len(bar.items) - 1
	}
}

func (bar *HueBar) Down(step int) {
	bar.lock.Lock()
	defer bar.lock.Unlock()
	bar.pos -= step
	if bar.pos < 0 {
		bar.pos += len(bar.items) - 1
	}
}

func (bar *HueBar) SetPointerStyle(st tcell.Style) { bar.pst = st }

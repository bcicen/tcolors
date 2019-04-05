package main

import (
	"sync"

	"github.com/gdamore/tcell"
)

type HueNavBar struct {
	items  []tcell.Color // navigation colors
	mItems []int         // minimap sample indices
	pos    int
	width  int
	lock   sync.RWMutex
}

func NewHueNavBar(width int) *HueNavBar {
	return &HueNavBar{width: width}
}

// Draw redraws bar at given coordinates and screen, returning the number
// of rows occupied
func (nb *HueNavBar) Draw(x, y int, s tcell.Screen) int {
	var st tcell.Style
	center := nb.width / 2

	//s.SetCell(center+x, y, indicatorSt, '⬇')
	s.SetCell(center+x-1, y, indicatorSt, '↿')
	s.SetCell(center+x+1, y, indicatorSt, '↾')
	//s.SetCell(center+x, y, indicatorSt, '﹀')

	for col, color := range nb.Items() {
		st = st.Background(color)
		s.SetCell(col+x, y+1, st, ' ')
		s.SetCell(col+x, y+2, st, ' ')
	}

	boxPad := nb.width / 30
	if boxPad < 2 {
		boxPad = 2
	}

	s.SetCell(center+x-boxPad, y+3, indicatorSt, '┌')
	s.SetCell(center+x+boxPad, y+3, indicatorSt, '┐')

	for col, color := range nb.MiniMap() {
		st = st.Background(color)
		s.SetCell(col+x, y+4, st, ' ')
	}

	s.SetCell(center+x-boxPad, y+5, indicatorSt, '└')
	s.SetCell(center+x+boxPad, y+5, indicatorSt, '┘')

	return 6
}

func (nb *HueNavBar) SetPos(n int)          { nb.pos = n }
func (nb *HueNavBar) Selected() tcell.Color { return nb.items[nb.pos] }

func (nb *HueNavBar) center() int { return (nb.width / 2) + 1 }

func (nb *HueNavBar) miniStep() int {
	n := len(nb.items) / nb.width
	if n > 13 {
		n = 13
	}
	return n
}

func (nb *HueNavBar) Resize(w int) {
	nb.width = w
	nb.Update(nb.items)
}

func (nb *HueNavBar) Update(a []tcell.Color) {
	nb.lock.Lock()
	defer nb.lock.Unlock()
	nb.items = a

	// build minimap indices
	nb.mItems = nb.mItems[0:]
	miniStep := nb.miniStep()
	for n := 0; n < len(nb.items); n += miniStep {
		nb.mItems = append(nb.mItems, n)
	}
}

func (nb *HueNavBar) Items() []tcell.Color {
	l := nb.pos - nb.center()
	r := nb.pos + nb.center()
	ilen := len(nb.items)
	if l < 0 {
		return append(nb.items[ilen+l:], nb.items[0:r]...)
	}
	if r > ilen {
		return append(nb.items[l:ilen-1], nb.items[0:r-ilen]...)
	}
	return nb.items[l:r]
}

func (nb *HueNavBar) MiniMap() []tcell.Color {
	var mpos int
	for mpos < len(nb.mItems)-1 {
		if nb.mItems[mpos+1] >= nb.pos {
			break
		}
		mpos++
	}

	l := mpos - nb.center()
	r := mpos + nb.center()
	ilen := len(nb.mItems)

	var a []int
	switch {
	case l < 0:
		a = append(nb.mItems[ilen+l:], nb.mItems[0:r]...)
	case r > ilen:
		a = append(nb.mItems[l:], nb.mItems[0:r-ilen]...)
	default:
		a = nb.mItems[l:r]
	}

	var colors []tcell.Color
	for _, idx := range a {
		colors = append(colors, nb.items[idx])
	}

	return colors
}

func (nb *HueNavBar) Up(step int) {
	nb.lock.Lock()
	defer nb.lock.Unlock()
	nb.pos += step
	if nb.pos >= len(nb.items)-1 {
		nb.pos -= len(nb.items) - 1
	}
}

func (nb *HueNavBar) Down(step int) {
	nb.lock.Lock()
	defer nb.lock.Unlock()
	nb.pos -= step
	if nb.pos < 0 {
		nb.pos += len(nb.items) - 1
	}
}

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
	loop   bool
	lock   sync.RWMutex
}

func NewHueNavBar() *HueNavBar {
	return &HueNavBar{}
}

func (nb *HueNavBar) SetPos(n int)          { nb.pos = n }
func (nb *HueNavBar) SetWidth(w int)        { nb.width = w }
func (nb *HueNavBar) Selected() tcell.Color { return nb.items[nb.pos] }

func (nb *HueNavBar) miniStep() int { return len(nb.items) / nb.width }
func (nb *HueNavBar) center() int   { return (nb.width / 2) + 1 }

func (nb *HueNavBar) Update(a []tcell.Color) {
	nb.lock.Lock()
	defer nb.lock.Unlock()
	nb.items = a

	// build minimap indices
	nb.mItems = nb.mItems[0:]
	for n := 0; n < len(nb.items); n += nb.miniStep() {
		nb.mItems = append(nb.mItems, n)
	}
}

func (nb *HueNavBar) Append(c tcell.Color) { nb.items = append(nb.items, c) }

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
	mpos := nb.pos / nb.miniStep()

	l := mpos - nb.center()
	r := mpos + nb.center()
	ilen := len(nb.mItems)

	var a []int
	switch {
	case l < 0:
		a = append(nb.mItems[ilen+l:], nb.mItems[0:ilen+l]...)
	case r > ilen:
		a = append(nb.mItems[l:], nb.mItems[0:r-ilen]...)
	default:
		a = nb.mItems[:]
	}

	var colors []tcell.Color
	for _, idx := range a {
		colors = append(colors, nb.items[idx])
	}

	return colors
}

package main

import (
	"sync"

	"github.com/gdamore/tcell"
)

type HueNavBar struct {
	items []tcell.Color // navigation colors
	pos   int
	width int
	loop  bool
	lock  sync.RWMutex
}

func NewHueNavBar() *HueNavBar {
	return &HueNavBar{}
}

func (nb *HueNavBar) SetPos(n int)          { nb.pos = n }
func (nb *HueNavBar) SetWidth(w int)        { nb.width = w }
func (nb *HueNavBar) Clear()                { nb.items = nb.items[:0] }
func (nb *HueNavBar) Append(c tcell.Color)  { nb.items = append(nb.items, c) }
func (nb *HueNavBar) Selected() tcell.Color { return nb.items[nb.pos] }

func (nb *HueNavBar) center() int { return (nb.width / 2) + 1 }

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

// return minimap sample indices
func (nb *HueNavBar) mItems() (a []int, mpos int) {
	mpos = -1
	step := len(nb.items) / nb.width

	for n := 0; n < len(nb.items); n += step {
		a = append(a, n)
		if mpos == -1 && n >= nb.pos {
			mpos = n
		}
	}
	if mpos == -1 {
		mpos = len(a) - 1
	}
	return
}

func (nb *HueNavBar) MiniMap() []tcell.Color {
	items, pos := nb.mItems()

	l := pos - nb.center()
	r := pos + nb.center()
	ilen := len(items)

	var a []int
	switch {
	case l < 0:
		a = append(items[ilen+l:], items[0:ilen+l]...)
	case r > ilen:
		a = append(items[l:], items[0:r-ilen]...)
	default:
		a = items[:]
	}

	var colors []tcell.Color
	for _, idx := range a {
		colors = append(colors, nb.items[idx])
	}

	return colors
}

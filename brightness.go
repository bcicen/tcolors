package main

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

const scrollAhead = 3

type BrightnessBar struct {
	items  []tcell.Color // navigation colors
	scale  []float64
	pos    int
	offset int
	width  int
	pst    tcell.Style // pointer style
	lock   sync.RWMutex
}

func NewBrightnessBar(width int) *BrightnessBar {
	bar := &BrightnessBar{width: width}
	for i := -0.50; i < 1.005; i += 0.005 {
		bar.scale = append(bar.scale, i)
	}
	bar.items = make([]tcell.Color, len(bar.scale))
	return bar
}

// Draw redraws bar at given coordinates and screen, returning the number
// of rows occupied
func (bar *BrightnessBar) Draw(x, y int, s tcell.Screen) int {
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

	s.SetCell(bar.width/2, y+3, bar.pst, []rune(fmt.Sprintf("%+3.2f", bar.Value()))...)

	return 4
}

func (bar *BrightnessBar) Value() float64 { return bar.scale[bar.pos] }
func (bar *BrightnessBar) SetValue(n float64) {
	var idx int
	for idx < len(bar.scale)-1 {
		if bar.scale[idx+1] > n {
			break
		}
		idx++
	}

	switch {
	case idx > bar.pos:
		bar.Up(idx - bar.pos)
	case idx < bar.pos:
		bar.Down(bar.pos - idx)
	}
}

func (bar *BrightnessBar) Resize(w int) {
	bar.width = w
	bar.Up(1)
	bar.Down(1)
}

func (bar *BrightnessBar) Update(base *noire.Color) {
	bar.lock.Lock()
	defer bar.lock.Unlock()

	for n, val := range bar.scale {
		bar.items[n] = toTColor(applyBrightness(val, base))
	}
}

func (bar *BrightnessBar) Up(step int) {
	bar.lock.Lock()
	defer bar.lock.Unlock()

	max := len(bar.items) - 1
	maxOffset := max - bar.width
	switch {
	case bar.pos == max:
		return
	case bar.pos+step > max:
		bar.pos = max
	default:
		bar.pos += step
	}

	if (bar.pos - bar.offset) > bar.width-scrollAhead {
		bar.offset = (bar.pos - bar.width) + scrollAhead
		if bar.offset >= maxOffset {
			bar.offset = maxOffset
		}
	}
}

func (bar *BrightnessBar) Down(step int) {
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

	if bar.pos-bar.offset < scrollAhead {
		bar.offset = bar.pos - scrollAhead
		if bar.offset < 0 {
			bar.offset = 0
		}
	}
}

func (bar *BrightnessBar) SetPointerStyle(st tcell.Style) { bar.pst = st }

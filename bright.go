package main

import (
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
	lock   sync.RWMutex
}

func NewBrightnessBar(width int) *BrightnessBar {
	bar := &BrightnessBar{width: width}
	for i := -0.49; i < 1.01; i += 0.01 {
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
		s.SetCell(col+x, y+1, st, ' ')
		s.SetCell(col+x, y+2, st, ' ')

		if n == bar.pos {
			s.SetCell(col+x, y, indicatorSt, '▼')
		} else {
			s.SetCell(col+x, y, indicatorSt, ' ')
		}
		col++
		n++
	}

	return 3
}

func (bar *BrightnessBar) Value() float64 { return bar.scale[bar.pos] }

func (bar *BrightnessBar) SetPos(n int)          { bar.pos = n }
func (bar *BrightnessBar) Selected() tcell.Color { return bar.items[bar.pos] }

func (bar *BrightnessBar) center() int { return (bar.width / 2) + 1 }

func (bar *BrightnessBar) miniStep() int {
	n := len(bar.items) / bar.width
	if n > 13 {
		n = 13
	}
	return n
}

func (bar *BrightnessBar) Resize(w int) {
	bar.width = w
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
	switch {
	case bar.pos == max:
		return
	case bar.pos+step > max:
		bar.pos = max
	default:
		bar.pos += step
	}

	if (bar.pos-bar.offset) > bar.width-scrollAhead && bar.pos <= len(bar.items)-scrollAhead {
		bar.offset = (bar.pos - bar.width) + scrollAhead
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

	if bar.pos-bar.offset < scrollAhead && bar.pos >= scrollAhead {
		bar.offset = bar.pos - scrollAhead
	}
}

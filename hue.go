package main

import (
	"github.com/bcicen/tcolors/state"
	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

const (
	hueMax   = 359.0
	hueIncr  = 0.5
	hueCount = int(hueMax / hueIncr)
)

type HueBar struct {
	items  [hueCount]tcell.Color // navigation colors
	mItems []int                 // minimap sample indices
	pos    int
	width  int
	height int
	pst    tcell.Style // pointer style
	state  *state.State
}

func NewHueBar(s *state.State) *HueBar { return &HueBar{state: s} }

func (bar *HueBar) SetPos(n float64) { bar.pos = int(n / hueIncr) }
func (bar *HueBar) center() int      { return (bar.width / 2) }

// Draw redraws bar at given coordinates and screen, returning the number
// of rows occupied
func (bar *HueBar) Draw(x, y int, s tcell.Screen) int {
	center := bar.width / 2
	st := tcell.StyleDefault.
		Foreground(tcell.ColorBlack)

	s.SetCell(center+x, y, bar.pst, '▾')

	// border bars
	for i := 1; i <= bar.height; i++ {
		s.SetCell(x-1, y+i, bar.pst, '│')
		s.SetCell(bar.width+x, y+i, bar.pst, '│')
	}

	// draw main bar
	idx := bar.pos - bar.center()
	if idx < 0 {
		idx += len(bar.items)
	}
	for col := 0; col < bar.width; col++ {
		if idx >= len(bar.items) {
			idx = 0
		}
		st = st.Background(bar.items[idx])
		for i := 1; i < bar.height-1; i++ {
			s.SetCell(col+x, y+i, st, ' ')
		}
		s.SetCell(col+x, y+bar.height-1, st, '▁')
		idx++
	}

	// draw minimap
	midx := bar.miniPos() - bar.center() - 1
	if midx < 0 {
		midx += len(bar.mItems)
	}
	for col := 0; col < bar.width; col++ {
		if midx >= len(bar.mItems) {
			midx = 0
		}
		idx = bar.mItems[midx]
		st = st.Background(bar.items[idx])
		s.SetCell(col+x, y+bar.height, st, ' ')
		midx++
	}

	y += bar.height + 1

	// draw map->minimap borders
	spacing := 1
	if bar.height > 2 {
		spacing = 2
		miniBoxW := bar.miniBox()
		s.SetCell(center+x-miniBoxW, y, bar.pst, '└')
		if bar.width%2 != 0 && miniBoxW%2 == 0 {
			miniBoxW++
		}
		s.SetCell(center+x+miniBoxW, y, bar.pst, '┘')
	}

	return bar.height + spacing
}

func (bar *HueBar) Resize(w, h int) {
	bar.height = barHeight(h) + 1
	bar.width = w
	bar.buildMini()
}

func (bar *HueBar) Handle(change state.Change) {
	if change.Includes(state.SaturationChanged, state.ValueChanged) {
		var hue float64
		for n := 0; n < hueCount; n++ {
			nc := noire.NewHSV(hue, bar.state.Saturation(), bar.state.Value())
			bar.items[n] = toTColor(nc)
			hue += hueIncr
		}
	}

	if change.Includes(state.SelectedChanged, state.HueChanged) {
		bar.SetPos(bar.state.Hue())
	}
}

// return width of minimap displayed in main map
func (bar *HueBar) miniBox() int {
	w := bar.width / 30
	if w < 2 {
		w = 2
	}
	return w
}

// return minimap step increment
func (bar *HueBar) miniStep() int {
	return roundFloat(float64(len(bar.items)) / float64(bar.width))
}

// return minimap index position (center) of active color
func (bar *HueBar) miniPos() int {
	var mpos int
	for mpos < len(bar.mItems)-1 {
		if bar.mItems[mpos+1] >= bar.pos {
			break
		}
		mpos++
	}
	return mpos
}

// build minimap indices
func (bar *HueBar) buildMini() {
	bar.mItems = []int{}
	miniStep := bar.miniStep()
	for n := 0; n < len(bar.items); n += miniStep {
		bar.mItems = append(bar.mItems, n)
	}
}

func (bar *HueBar) Up(step int) {
	n := int(bar.state.Hue()) + step
	if n > hueMax {
		n -= hueMax - 1
	}
	bar.state.SetHue(float64(n))
}

func (bar *HueBar) Down(step int) {
	n := int(bar.state.Hue()) - step
	if n < 0 {
		n += hueMax - 1
	}
	bar.state.SetHue(float64(n))
}

func (bar *HueBar) SetPointerStyle(st tcell.Style) { bar.pst = st }

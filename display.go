package main

import (
	"sync"

	"github.com/bcicen/tcolors/state"
	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

const (
	paddingX = 2
	maxWidth = 105
)

type ChangeHandler interface {
	Handle(state.Change)
}

type Section interface {
	ChangeHandler
	Up(int)
	Down(int)
	Draw(int, int, tcell.Screen) int
	Width() int
	Resize(int) // resize section to given width
	SetPointerStyle(tcell.Style)
}

type Display struct {
	rgb      []int32
	sections []Section
	sectionN int
	width    int
	state    *state.State
	lock     sync.RWMutex
}

func NewDisplay() *Display {
	state := state.NewDefault()
	err := state.Load()
	if err != nil {
		panic(err)
	}
	d := &Display{
		state: state,
	}
	d.sections = []Section{
		NewPaletteBox(state),
		NewHueBar(state),
		NewSaturationBar(state),
		NewValueBar(state),
	}
	return d
}

// Draw redraws display at given coordinates, returning the number
// of rows occupied
func (d *Display) Draw(x, y int, s tcell.Screen) int {
	d.lock.Lock()
	defer d.lock.Unlock()
	for n, sec := range d.sections {
		if n == d.sectionN {
			sec.SetPointerStyle(hiIndicatorSt)
		} else {
			sec.SetPointerStyle(indicatorSt)
		}
		y += sec.Draw(x, y, s)
	}
	return y
}

func (d *Display) Resize(w int) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.width = w - ((paddingX * 2) + 1)
	if d.width > maxWidth {
		d.width = maxWidth
	}
	for _, sec := range d.sections {
		sec.Resize(d.width)
	}
}

func toTColor(c *noire.Color) tcell.Color {
	r, g, b := c.RGB()
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}

func (d *Display) build() {
	change := d.state.Flush()
	log.Debugf("handling change: %08b", change)
	log.Debugf("state: [h=%0.3f s=%0.3f v=%0.3f]", d.state.Hue(), d.state.Saturation(), d.state.Value())
	for _, sec := range d.sections {
		sec.Handle(change)
	}
}

func (d *Display) SetColor(c tcell.Color) {
	r, g, b := c.RGB()
	nc := noire.NewRGB(float64(r), float64(g), float64(b))
	h, s, v := nc.HSV()

	d.state.SetSaturation(s)
	d.state.SetValue(v)
	d.state.SetHue(h)
	d.build()
}

func (d *Display) SectionUp() (ok bool) {
	if d.sectionN == 0 {
		return false
	}
	d.sectionN -= 1
	return true
}

func (d *Display) SectionDown() (ok bool) {
	if d.sectionN == len(d.sections)-1 {
		return false
	}
	d.sectionN += 1
	return true
}

func (d *Display) ValueUp(step int) (ok bool) {
	d.sections[d.sectionN].Up(step)
	d.build()
	return true
}

func (d *Display) ValueDown(step int) (ok bool) {
	d.sections[d.sectionN].Down(step)
	d.build()
	return true
}

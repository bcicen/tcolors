package main

import (
	"sync"

	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

const (
	paddingX     = 2
	paddingY     = 1
	step         = 0.005 // default step for bar scale
	maxWidth     = 200
	defaultGlyph = ' '
)

type ChangeHandler interface {
	Handle(StateChange)
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
	rgb        []int32
	HueNav     *HueBar
	SatNav     *SaturationBar
	ValueNav   *ValueBar
	PaletteNav *PaletteBox
	sections   []Section
	sectionN   int
	hues       []*noire.Color // modified hues
	xHues      []*noire.Color // base hues
	width      int
	height     int
	state      *State
	lock       sync.RWMutex
}

func NewDisplay() *Display {
	state := NewDefaultState()
	d := &Display{
		state:      state,
		HueNav:     NewHueBar(state),
		SatNav:     NewSaturationBar(state),
		ValueNav:   NewValueBar(state),
		PaletteNav: NewPaletteBox(state),
	}
	d.sections = []Section{
		d.PaletteNav,
		d.HueNav,
		d.SatNav,
		d.ValueNav,
	}
	d.mkhues()
	return d
}

// Draw redraws display at given coordinates, returning the number
// of rows occupied
func (d *Display) Draw(x, y int, s tcell.Screen) int {
	for n, sec := range d.sections {
		if n == d.sectionN {
			sec.SetPointerStyle(hiIndicatorSt)
		} else {
			sec.SetPointerStyle(indicatorSt)
		}
		if n == 0 {
			offset := (d.width - sec.Width()) / 2
			y += sec.Draw(x+offset, y, s)
		} else {
			y += sec.Draw(x, y, s)
		}
	}
	return y
}

func (d *Display) Reset() {
	d.HueNav.SetPos(100)
	d.SatNav.SetPos(100)
	d.ValueNav.SetPos(100)
	d.build()
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

func (d *Display) mkhues() {
	for i := 0.0; i < 359; i += 0.5 {
		d.xHues = append(d.xHues, noire.NewHSV(float64(i), 100, 100))
	}
}

func toTColor(c *noire.Color) tcell.Color {
	r, g, b := c.RGB()
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}

func (d *Display) build() {
	change := d.state.Flush()
	log.Infof("handling change: %08b", change)
	log.Infof("state: [h=%0.3f s=%0.3f v=%0.3f]", d.state.Hue(), d.state.Saturation(), d.state.Value())
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

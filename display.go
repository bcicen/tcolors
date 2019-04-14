package main

import (
	"sync"

	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

const (
	padding      = 2
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
	Resize(int)
	SetPointerStyle(tcell.Style)
}

type Display struct {
	rgb        []int32
	HueNav     *HueBar
	SatNav     *SaturationBar
	BrightNav  *BrightnessBar
	PaletteNav *PaletteBox
	sections   []Section
	sectionN   int
	hues       []*noire.Color // modified hues
	xHues      []*noire.Color // base hues
	bigStep    bool           // navigation step basis
	width      int
	state      *State
	lock       sync.RWMutex
}

func NewDisplay(width int) *Display {
	state := NewDefaultState()
	d := &Display{
		state:      state,
		HueNav:     NewHueBar(state),
		SatNav:     NewSaturationBar(state),
		BrightNav:  NewBrightnessBar(state),
		PaletteNav: NewPaletteBox(state),
	}
	d.sections = []Section{
		d.PaletteNav,
		d.HueNav,
		d.SatNav,
		d.BrightNav,
	}
	d.mkhues()
	d.Resize(width)
	d.Reset()
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
		y += sec.Draw(x, y, s)
	}
	return y
}

func (d *Display) Reset() {
	d.HueNav.SetValue(100)
	d.SatNav.SetValue(100)
	d.BrightNav.SetValue(100)
	d.build()
}

func (d *Display) Resize(w int) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.width = w - ((padding * 2) + 1)
	if d.width > maxWidth {
		d.width = maxWidth
	}
	for _, sec := range d.sections {
		sec.Resize(d.width)
	}
}

func (d *Display) Saturation() float64   { return d.SatNav.Value() }
func (d *Display) Brightness() float64   { return d.BrightNav.Value() }
func (d *Display) Selected() tcell.Color { return d.HueNav.Selected() }

func (d *Display) mkhues() {
	for i := 0.0; i < 359; i += 0.5 {
		d.xHues = append(d.xHues, noire.NewHSV(float64(i), 100, 100))
	}
}

func applySaturation(s float64, c *noire.Color) *noire.Color {
	h := c.Hue()
	l := c.Lightness()
	return noire.NewHSV(h, s, l)
}

func applyValue(l float64, c *noire.Color) *noire.Color {
	h := c.Hue()
	s := c.Saturation()
	return noire.NewHSV(h, s, l)
}

func toTColor(c *noire.Color) tcell.Color {
	r, g, b := c.RGB()
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}

func (d *Display) build() {
	change := d.state.Flush()
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

func (d *Display) ValueUp() (ok bool) {
	d.sections[d.sectionN].Up(d.stepSize())
	d.build()
	return true
}

func (d *Display) ValueDown() (ok bool) {
	d.sections[d.sectionN].Down(d.stepSize())
	d.build()
	return true
}

func (d *Display) ToggleStep() {
	d.bigStep = d.bigStep != true
}

func (d *Display) stepSize() int {
	if d.bigStep {
		return 10
	}
	return 2
}

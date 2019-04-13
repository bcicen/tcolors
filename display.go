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

type Section interface {
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
	xHues      []*noire.Color // base hues
	bigStep    bool           // navigation step basis
	width      int
	lock       sync.RWMutex
}

func NewDisplay(width int) *Display {
	d := &Display{
		HueNav:     NewHueBar(0),
		SatNav:     NewSaturationBar(0),
		BrightNav:  NewBrightnessBar(0),
		PaletteNav: NewPaletteBox(0),
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
	d.HueNav.SetPos(0)
	d.SatNav.SetValue(1)
	d.BrightNav.SetValue(0)
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
	var (
		incr float64 = 1
		r    float64 = 255
		g    float64 = 0
		b    float64 = 0
	)
	d.xHues = append(d.xHues, noire.NewRGB(r, g, b))

	for b < 256 {
		b += incr
		d.xHues = append(d.xHues, noire.NewRGB(r, g, b))
	}
	for r > 0 {
		r -= incr
		d.xHues = append(d.xHues, noire.NewRGB(r, g, b))
	}
	for g < 256 {
		g += incr
		d.xHues = append(d.xHues, noire.NewRGB(r, g, b))
	}
	for b > 0 {
		b -= incr
		d.xHues = append(d.xHues, noire.NewRGB(r, g, b))
	}
	for r < 256 {
		r += incr
		d.xHues = append(d.xHues, noire.NewRGB(r, g, b))
	}
	for g > 0 {
		g -= incr
		d.xHues = append(d.xHues, noire.NewRGB(r, g, b))
	}
}

func applySaturation(level float64, c *noire.Color) *noire.Color {
	switch {
	case level > step:
		return c.Saturate(level)
	case level < -step:
		return c.Desaturate(level * -1)
	}
	return c
}

func applyBrightness(level float64, c *noire.Color) *noire.Color {
	switch {
	case level > step:
		return c.Brighten(level)
	case level < -step:
		return c.Darken(level * -1)
	}
	return c
}

func toTColor(c *noire.Color) tcell.Color {
	r, g, b := c.RGB()
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}

func (d *Display) build() {
	var n int
	var c *noire.Color
	buf := make([]tcell.Color, len(d.xHues))

	for n = range d.xHues {
		c = applySaturation(d.Saturation(), d.xHues[n])
		c = applyBrightness(d.Brightness(), c)
		buf[n] = toTColor(c)
	}
	d.HueNav.Update(buf[0:n])

	d.SatNav.Update(d.xHues[d.HueNav.pos])
	d.BrightNav.Update(d.xHues[d.HueNav.pos])
	d.PaletteNav.Update(d.HueNav.Selected())
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

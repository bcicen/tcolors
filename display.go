package main

import (
	"sync"

	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

const (
	padding            = 2
	navIncr      uint8 = 5
	navMax       uint8 = 200
	navMin       uint8 = 0
	maxWidth           = 1200
	defaultGlyph       = ' '
)

type Display struct {
	rgb        []int32
	HueNav     *HueBar
	BrightNav  *BrightnessBar
	xHues      []*noire.Color // base hues
	saturation uint8          // 0 to 200
	center     int
	screen     tcell.Screen
	lock       sync.RWMutex
}

func NewDisplay(s tcell.Screen) *Display {
	d := &Display{
		screen:    s,
		HueNav:    NewHueBar(0),
		BrightNav: NewBrightnessBar(0),
	}
	d.mkhues()
	d.Reset()
	return d
}

func (d *Display) Reset() {
	d.saturation = 100
	d.Resize()
	d.HueNav.SetPos(0)
	d.BrightNav.SetPos(60)
	d.build()
}

func (d *Display) Resize() {
	d.lock.Lock()
	defer d.lock.Unlock()
	w, _ := d.screen.Size()
	w = w - ((padding * 2) + 1)
	d.HueNav.Resize(w)
	d.BrightNav.Resize(w)
}

func (d *Display) Saturation() float64   { return (float64(d.saturation) / 100) - 1 }
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
	if level < 0 {
		return c.Desaturate(level * -1)
	}
	return c.Saturate(level)
}

func applyBrightness(level float64, c *noire.Color) *noire.Color {
	if level < 0 {
		return c.Darken(level * -1)
	}
	return c.Brighten(level)
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

	d.BrightNav.Update(d.xHues[d.HueNav.pos])
}

func (d *Display) SaturationUp() (ok bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.saturation == navMax {
		return false
	}
	d.saturation += navIncr
	d.build()
	return true
}

func (d *Display) SaturationDown() (ok bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.saturation == navMin {
		return false
	}
	d.saturation -= navIncr
	d.build()
	return true
}

func (d *Display) HueUp(step int) (ok bool) {
	d.HueNav.Up(step)
	return true
}

func (d *Display) HueDown(step int) (ok bool) {
	d.HueNav.Down(step)
	return true
}

func (d *Display) BrightnessUp(step int) (ok bool) {
	d.BrightNav.Up(step)
	d.build()
	return true
}

func (d *Display) BrightnessDown(step int) (ok bool) {
	d.BrightNav.Down(step)
	d.build()
	return true
}

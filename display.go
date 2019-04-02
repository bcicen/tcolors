package main

import (
	"math"
	"sync"

	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

const (
	padding       = 2
	navIncr uint8 = 5
	navMax  uint8 = 200
	navMin  uint8 = 0
)

type Display struct {
	rgb        []int32
	hues       []tcell.Color  // hue nav
	mHues      []int          // minimap hue indices
	xHues      []*noire.Color // base hues
	saturation uint8          // 0 to 200
	brightness uint8          // 0 to 200
	pos        int
	center     int
	width      int
	screen     tcell.Screen
	lock       sync.RWMutex
}

func NewDisplay(s tcell.Screen) *Display {
	d := &Display{
		screen: s,
	}
	d.mkhues()
	d.Reset()
	return d
}

func (d *Display) Reset() {
	d.saturation = 100
	d.brightness = 100
	d.Resize()
	d.pos = d.center
}

func (d *Display) Saturation() float64   { return (float64(d.saturation) / 100) - 1 }
func (d *Display) Brightness() float64   { return (float64(d.brightness) / 100) - 1 }
func (d *Display) Selected() tcell.Color { return d.hues[d.pos] }

func (d *Display) MiniHues() []tcell.Color {
	var n int
	for n < len(d.mHues)-1 {
		if d.mHues[n+1] >= d.pos {
			break
		}
		n++
	}

	l := n - (d.center + 1)
	r := n + (d.center + 1)
	hlen := len(d.mHues)

	var a []int
	switch {
	case l < 0:
		a = append(d.mHues[hlen+l:], d.mHues[0:hlen+l]...)
	case r > hlen:
		a = append(d.mHues[l:], d.mHues[0:r-hlen]...)
	default:
		a = d.mHues[:]
	}

	var colors []tcell.Color
	for _, idx := range a {
		colors = append(colors, d.hues[idx])
	}

	return colors
}

func (d *Display) Hues() []tcell.Color {
	l := d.pos - (d.center + 1)
	r := d.pos + (d.center + 1)
	hlen := len(d.hues)
	if l < 0 {
		return append(d.hues[hlen+l:], d.hues[0:r]...)
	}
	if r > hlen {
		return append(d.hues[l:hlen-1], d.hues[0:r-hlen]...)
	}
	return d.hues[l:r]
}

func (d *Display) Resize() {
	w, _ := d.screen.Size()
	d.width = w - (padding * 3)
	d.center = d.width / 2
	d.build()
}

func (d *Display) mkhues() {
	var (
		incr float64 = 1
		r    float64 = 255
		g    float64 = 0
		b    float64 = 0
	)
	d.xHues = append(d.xHues, noire.NewRGB(r, g, b))

	for b < 255 {
		b += incr
		d.xHues = append(d.xHues, noire.NewRGB(r, g, b))
	}
	for r > 0 {
		r -= incr
		d.xHues = append(d.xHues, noire.NewRGB(r, g, b))
	}
	for g < 255 {
		g += incr
		d.xHues = append(d.xHues, noire.NewRGB(r, g, b))
	}
	for b > 0 {
		b -= incr
		d.xHues = append(d.xHues, noire.NewRGB(r, g, b))
	}
	for r < 255 {
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

func (d *Display) build() {
	var tc tcell.Color
	d.hues = d.hues[:0]
	d.mHues = d.mHues[:0]
	miniStep := len(d.xHues) / d.width

	for idx, c := range d.xHues {
		c = applySaturation(d.Saturation(), c)
		c = applyBrightness(d.Brightness(), c)
		r, g, b := c.RGB()
		tc = tcell.NewRGBColor(int32(r), int32(g), int32(b))
		d.hues = append(d.hues, tc)
		if idx == 0 || idx%miniStep == 0 {
			d.mHues = append(d.mHues, idx)
		}
	}
}

func (d *Display) Pos() int  { return d.pos }
func (d *Display) HPos() int { return int(math.Mod(float64(d.pos), 51.0)) }

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
	d.lock.Lock()
	defer d.lock.Unlock()
	d.pos += step
	if d.pos >= len(d.hues)-1 {
		d.pos -= len(d.hues) - 1
	}
	d.build()
	return true
}

func (d *Display) HueDown(step int) (ok bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.pos -= step
	if d.pos < 0 {
		d.pos += len(d.hues) - 1
	}
	d.build()
	return true
}

func (d *Display) BrightnessUp() (ok bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.brightness == navMax {
		return false
	}
	d.brightness += navIncr
	d.build()
	return true
}

func (d *Display) BrightnessDown() (ok bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.brightness == navMin {
		return false
	}
	d.brightness -= navIncr
	d.build()
	return true
}

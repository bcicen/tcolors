package main

import (
	"sync"
	"time"

	"github.com/bcicen/tcolors/state"
	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

const (
	paddingX   = 2
	maxWidth   = 105
	littleStep = 1
	bigStep    = 10
)

type ChangeHandler interface {
	Handle(state.Change)
}

type Section interface {
	ChangeHandler
	Up(int)
	Down(int)
	Draw(int, int, tcell.Screen) int
	Resize(int) // resize section to given width
	SetPointerStyle(tcell.Style)
}

type Display struct {
	rgb       []int32
	sections  []Section
	sectionN  int
	width     int
	errMsg    *ErrorMsg
	stepBasis int
	state     *state.State
	quit      chan struct{}
	lock      sync.RWMutex
}

func NewDisplay(s tcell.Screen) *Display {
	tstate, err := state.Load()

	d := &Display{
		state:  tstate,
		errMsg: NewErrorMsg(),
		quit:   make(chan struct{}),
		sections: []Section{
			NewPaletteBox(tstate),
			NewHueBar(tstate),
			NewSaturationBar(tstate),
			NewValueBar(tstate),
		},
	}

	if err != nil {
		d.errMsg.Set(err.Error())
	}

	w, _ := s.Size()
	d.Resize(w)
	d.build()

	go d.eventHandler(s)
	return d
}

func (d *Display) Done() error {
	for {
		select {
		case <-d.quit:
			return d.state.Save()
		case <-time.After(time.Millisecond * 50):
		}
	}
}

func (d *Display) Draw(s tcell.Screen) {
	d.lock.Lock()
	defer d.lock.Unlock()

	now := time.Now()

	w, h := s.Size()
	if w == 0 || h == 0 {
		return
	}

	x, y := paddingX, 1
	if d.width == maxWidth {
		x = (w - maxWidth) / 2 // center display
	}

	if d.stepBasis == bigStep {
		s.SetCell(1, 0, tcell.StyleDefault, '⏩')
	} else {
		s.SetCell(1, 0, tcell.StyleDefault, ' ')
	}

	for n, sec := range d.sections {
		if n == d.sectionN {
			sec.SetPointerStyle(hiIndicatorSt)
		} else {
			sec.SetPointerStyle(indicatorSt)
		}
		y += sec.Draw(x, y, s)
	}
	d.errMsg.Draw(x, s)

	s.Show()
	log.Debugf("draw [%3.3fms]", time.Since(now).Seconds()*1000)
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
	d.errMsg.Resize(d.width)
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

func (d *Display) ValueUp() (ok bool) {
	d.sections[d.sectionN].Up(d.stepBasis)
	d.build()
	return true
}

func (d *Display) ValueDown() (ok bool) {
	d.sections[d.sectionN].Down(d.stepBasis)
	d.build()
	return true
}

func (d *Display) eventHandler(s tcell.Screen) {
	for {
		redraw := false
		d.stepBasis = littleStep

		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Modifiers()&tcell.ModShift == tcell.ModShift {
				d.stepBasis = bigStep
			}
			switch ev.Key() {
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'h':
					redraw = d.ValueDown()
				case 'k':
					redraw = d.SectionUp()
				case 'j':
					redraw = d.SectionDown()
				case 'l':
					redraw = d.ValueUp()
				case 'q':
					close(d.quit)
					return
				default:
					log.Debugf("ignoring key [%s]", string(ev.Rune()))
				}
			case tcell.KeyRight:
				redraw = d.ValueUp()
			case tcell.KeyLeft:
				redraw = d.ValueDown()
			case tcell.KeyUp:
				redraw = d.SectionUp()
			case tcell.KeyDown:
				redraw = d.SectionDown()
			case tcell.KeyEscape, tcell.KeyCtrlC:
				close(d.quit)
				return
			case tcell.KeyCtrlL:
				s.Sync()
			}

		case *tcell.EventResize:
			w, h := s.Size()
			log.Debugf("handling resize: w=%04d h=%04d", w, h)
			d.Resize(w)
			s.Clear()
			s.Sync()
			redraw = true
		}

		if redraw {
			d.Draw(s)
		}
	}
}

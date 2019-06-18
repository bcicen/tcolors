package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/bcicen/tcolors/state"
	"github.com/bcicen/tcolors/widgets"
	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

const (
	paddingX   = 2
	minWidth   = 26
	minHeight  = 22
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
	Resize(int, int) // resize section to given width and height
	SetPointerStyle(tcell.Style)
}

type Display struct {
	rgb       []int32
	sections  []Section
	sectionN  int
	xPos      int
	width     int
	errMsg    *widgets.ErrorMsg
	stepBasis int
	state     *state.State
	quit      chan struct{}
	lock      sync.RWMutex
}

func NewDisplay(s tcell.Screen, tstate *state.State) *Display {
	d := &Display{
		state:  tstate,
		errMsg: widgets.NewErrorMsg(),
		quit:   make(chan struct{}),
		sections: []Section{
			widgets.NewPaletteBox(tstate),
			widgets.NewHueBar(tstate),
			widgets.NewSaturationBar(tstate),
			widgets.NewValueBar(tstate),
		},
	}

	w, h := s.Size()
	d.Resize(w, h)
	d.build()

	//if d.state.IsNew() {
	//msg := fmt.Sprintf("creating new palette file: %s", d.state.Path())
	//d.errMsg.Set(msg)
	//}

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

func (d *Display) drawSizeErr(s tcell.Screen) {
	w, h := s.Size()
	st := widgets.ErrSt
	s.SetCell(1, 0, st, []rune("screen too small!")...)
	s.SetCell(1, 1, st, []rune(fmt.Sprintf("[cur] %dx%d", w, h))...)
	s.SetCell(1, 2, st, []rune(fmt.Sprintf("[min] %dx%d", minWidth, minHeight))...)
	s.Show()
}

func (d *Display) Draw(s tcell.Screen) {
	d.lock.Lock()
	defer d.lock.Unlock()

	timer := log.NewTimer("draw")
	defer timer.End()

	if d.width < 0 {
		d.drawSizeErr(s)
		return
	}

	x, y := d.xPos, 0

	// draw header
	if d.stepBasis == bigStep {
		s.SetCell(x, y, widgets.IndicatorSt, '⏩')
	} else {
		s.SetCell(x, y, widgets.IndicatorSt, '⏵')
	}

	sname := d.state.Name()
	s.SetCell((x+d.width)-len(sname), y, widgets.IndicatorSt, []rune(sname)...)
	y += 2

	// draw sections
	for n, sec := range d.sections {
		if n == d.sectionN {
			sec.SetPointerStyle(widgets.HiIndicatorSt)
		} else {
			sec.SetPointerStyle(widgets.IndicatorSt)
		}
		y += sec.Draw(x, y, s)
	}
	d.errMsg.Draw(x, s)

	s.Show()
}

func (d *Display) Resize(w, h int) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if w < minWidth || h < minHeight {
		d.width = -1
		return
	}

	d.width = w - ((paddingX * 2) + 1)
	if d.width > maxWidth {
		d.width = maxWidth
	}
	// ensure total width aligns well with palette count
	d.xPos = (w - d.width) / 2 // center display

	d.width = (d.width / d.state.Len()) * d.state.Len()
	for _, sec := range d.sections {
		sec.Resize(d.width, h)
	}
	d.errMsg.Resize(d.width)
}

func (d *Display) build() {
	change := d.state.Flush()
	log.Debugf("handling change: %08b", change)
	log.Debugf("state: [h=%0.3f s=%0.3f v=%0.3f]", d.state.Hue(), d.state.Saturation(), d.state.Value())
	timer := log.NewTimer("handle")
	defer timer.End()

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
		resize := false
		d.stepBasis = littleStep

		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Modifiers()&tcell.ModShift == tcell.ModShift && d.sectionN != 0 {
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
				case 'a':
					d.state.Add()
					resize = true
				case 'x':
					d.state.Remove()
					resize = true
				case 'q':
					close(d.quit)
					return
				default:
					log.Debugf("ignoring rune key [%s]", string(ev.Rune()))
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
			case tcell.KeyInsert:
				d.state.Add()
				resize = true
			case tcell.KeyDelete:
				d.state.Remove()
				resize = true
			default:
				log.Debugf("ignoring event key [%s]", ev.Name())
			}

		case *tcell.EventResize:
			resize = true
		}

		if resize {
			w, h := s.Size()
			log.Debugf("handling resize: w=%04d h=%04d", w, h)
			d.Resize(w, h)
			s.Clear()
			s.Sync()
			redraw = true
		}

		if redraw {
			d.Draw(s)
		}
	}
}

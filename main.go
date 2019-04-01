package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

var (
	boxW = 120
	boxH = 80
	disp *Display
)

func draw(s tcell.Screen) {
	w, h := s.Size()

	if w == 0 || h == 0 {
		return
	}

	lh := h / 2
	lw := w / 2
	lx := w / 4
	ly := h / 4
	st := tcell.StyleDefault
	gl := ' '

	st = st.Background(disp.Selected())

	for row := 0; row < lh; row++ {
		for col := 0; col < lw; col++ {
			s.SetCell(lx+col, ly+row, st, gl)
		}
	}

	s.SetCell(1, 1, tcell.StyleDefault, []rune(fmt.Sprintf("%03d %3.3f", disp.brightness, disp.Brightness()))...)
	s.SetCell(1, 2, tcell.StyleDefault, []rune(fmt.Sprintf("%03d %3.3f", disp.saturation, disp.Saturation()))...)
	r, g, b := disp.Selected().RGB()
	s.SetCell(1, 3, tcell.StyleDefault, []rune(fmt.Sprintf("%03d %03d %03d", r, g, b))...)
	s.SetCell(1, 4, tcell.StyleDefault, []rune(fmt.Sprintf("%04d", disp.pos))...)
	s.SetCell(1, 5, tcell.StyleDefault, []rune(fmt.Sprintf("%04d", disp.center))...)

	for col, color := range disp.Hues() {
		st = st.Background(color)
		s.SetCell(col+padding, ly+lh+1, st, gl)
		s.SetCell(col+padding, ly+lh+2, st, gl)
		if col == disp.center {
			s.SetCell(col+padding, ly+lh, tcell.StyleDefault, '↓')
			s.SetCell(col+padding, ly+lh+3, tcell.StyleDefault, '↑')
		} else {
			s.SetCell(col+padding, ly+lh, tcell.StyleDefault, ' ')
			s.SetCell(col+padding, ly+lh+3, tcell.StyleDefault, ' ')
		}
	}

	for col, color := range disp.MiniHues() {
		st = st.Background(color)
		s.SetCell(col+padding, ly+lh+4, st, gl)
	}

	s.Show()
}

func main() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorDefault))
	s.Clear()
	disp = NewDisplay(s)

	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyRune:
					switch ev.Rune() {
					case 'r':
						disp.Reset()
						draw(s)
					case 'l':
						if ok := disp.SaturationUp(); ok {
							draw(s)
						}
					case 'h':
						if ok := disp.SaturationDown(); ok {
							draw(s)
						}
					}
				case tcell.KeyRight:
					if ok := disp.HueUp(10); ok {
						draw(s)
					}
				case tcell.KeyLeft:
					if ok := disp.HueDown(10); ok {
						draw(s)
					}
				case tcell.KeyUp:
					if ok := disp.BrightnessUp(); ok {
						draw(s)
					}
				case tcell.KeyDown:
					if ok := disp.BrightnessDown(); ok {
						draw(s)
					}
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	draw(s)

loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Millisecond * 50):
		}
	}

	w, h := s.Size()
	s.Fini()
	fmt.Printf("w=%d h=%d hues=%d\n", w, h, len(disp.hues))
}

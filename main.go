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
	navHeight := 1
	w, h := s.Size()

	if w == 0 || h == 0 {
		return
	}

	if h > 50 {
		navHeight = 2
	}

	lh := h / 4
	lw := w / 2
	lx := w / 4
	ly := 1
	st := tcell.StyleDefault
	gl := ' '
	indicatorSt := tcell.StyleDefault.
		Foreground(tcell.NewRGBColor(120, 120, 120)).
		Background(tcell.ColorBlack)

	st = st.Background(disp.Selected())

	for row := 0; row < lh; row++ {
		for col := 0; col < lw; col++ {
			s.SetCell(lx+col, ly, st, gl)
		}
		ly++
	}

	r, g, b := disp.Selected().RGB()
	s.SetCell((w-11)/2, ly, tcell.StyleDefault, []rune(fmt.Sprintf("%03d %03d %03d", r, g, b))...)
	ly += 2

	s.SetCell(disp.center+padding, ly, indicatorSt, '﹀')

	for col, color := range disp.HueNav.Items() {
		st = st.Background(color)
		s.SetCell(col+padding, ly+1, st, gl)
		if navHeight == 2 {
			s.SetCell(col+padding, ly+2, st, gl)
		}
	}

	ly += navHeight
	//s.SetCell(disp.center+padding, ly, indicatorSt, '︿')

	ly++

	s.SetCell(disp.center+padding-3, ly, indicatorSt, '┌')
	s.SetCell(disp.center+padding+4, ly, indicatorSt, '┐')

	mmap := disp.HueNav.MiniMap()
	for col, color := range mmap {
		st = st.Background(color)
		s.SetCell(col+padding, ly+1, st, gl)
		//if navHeight
		//s.SetCell(col+padding, ly+2, st, gl)
		//}
	}
	ly += 2
	//ly += navHeight + 1

	s.SetCell(disp.center+padding-3, ly, indicatorSt, '└')
	s.SetCell(disp.center+padding+4, ly, indicatorSt, '┘')

	s.SetCell(1, h-6, tcell.StyleDefault, []rune(fmt.Sprintf("%03d %3.3f", disp.brightness, disp.Brightness()))...)
	s.SetCell(1, h-5, tcell.StyleDefault, []rune(fmt.Sprintf("%03d %3.3f", disp.saturation, disp.Saturation()))...)
	s.SetCell(1, h-3, tcell.StyleDefault, []rune(fmt.Sprintf("%04d [w=%04d]", disp.HueNav.pos, disp.HueNav.width))...)
	s.SetCell(1, h-2, tcell.StyleDefault, []rune(fmt.Sprintf("%04d", len(mmap)))...)

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
		Background(tcell.ColorBlack))
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
	fmt.Printf("w=%d h=%d hues=%d\n", w, h, len(disp.HueNav.items))
}

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

var (
	boxW  = 120
	boxH  = 80
	disp  *Display
	blkSt = tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorBlack)
	indicatorSt = tcell.StyleDefault.
			Foreground(tcell.NewRGBColor(110, 110, 110)).
			Background(tcell.ColorBlack)
	hiIndicatorSt = tcell.StyleDefault.
			Foreground(tcell.NewRGBColor(255, 255, 255)).
			Background(tcell.ColorBlack)
)

func draw(s tcell.Screen) {
	w, h := s.Size()

	if w == 0 || h == 0 {
		return
	}

	lh := h / 4
	lw := w / 2
	lx := w / 4
	ly := 1
	st := tcell.StyleDefault
	gl := ' '

	if disp.bigStep {
		s.SetCell(1, 0, st, '⏩')
	} else {
		s.SetCell(1, 0, st, ' ')
	}

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

	ly += disp.Draw(padding, ly, s)

	s.SetCell(1, h-3, tcell.StyleDefault, []rune(fmt.Sprintf("%04d [w=%04d]", disp.HueNav.pos, disp.HueNav.width))...)
	s.SetCell(1, h-2, tcell.StyleDefault, []rune(fmt.Sprintf("%04d [off=%04d] [i=%04d] [w=%04d]", disp.BrightNav.pos, disp.BrightNav.offset, len(disp.BrightNav.items), disp.BrightNav.width))...)

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
	w, _ := s.Size()
	disp = NewDisplay(w)

	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyRune:
					switch ev.Rune() {
					case 's':
						disp.ToggleStep()
						draw(s)
					case 'r':
						disp.Reset()
						draw(s)
					case 'l':
						if ok := disp.ValueDown(); ok {
							draw(s)
						}
					case 'h':
						if ok := disp.ValueUp(); ok {
							draw(s)
						}
					case 'q':
						close(quit)
						return
					}
				case tcell.KeyRight:
					if ok := disp.ValueUp(); ok {
						draw(s)
					}
				case tcell.KeyLeft:
					if ok := disp.ValueDown(); ok {
						draw(s)
					}
				case tcell.KeyUp:
					if ok := disp.SectionUp(); ok {
						draw(s)
					}
				case tcell.KeyDown:
					if ok := disp.SectionDown(); ok {
						draw(s)
					}
				case tcell.KeyEscape, tcell.KeyCtrlC:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				}
			case *tcell.EventResize:
				w, _ := s.Size()
				disp.Resize(w)
				s.Clear()
				draw(s)
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
	fmt.Printf("w=%d h=%d hues=%d bscale=%v\n", w, h, len(disp.HueNav.items), len(disp.BrightNav.scale))
	//for n, x := range disp.BrightNav.scale {
	//fmt.Printf("[%d] %+0.3f\n", n, x)
	//}
	//fmt.Printf("%v\n", disp.BrightNav.scale)
}

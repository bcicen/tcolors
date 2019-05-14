package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bcicen/tcolors/logging"
	"github.com/gdamore/tcell"
)

const (
	littleStep = 1
	bigStep    = 10
)

var (
	log   = logging.Init()
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
	errSt = tcell.StyleDefault.
		Foreground(tcell.NewRGBColor(255, 000, 043)).
		Background(tcell.ColorBlack)
	stepBasis int
)

func draw(s tcell.Screen) {
	w, h := s.Size()

	if w == 0 || h == 0 {
		return
	}

	ly := 1
	st := tcell.StyleDefault

	if stepBasis == bigStep {
		s.SetCell(1, 0, st, '⏩')
	} else {
		s.SetCell(1, 0, st, ' ')
	}

	x := paddingX
	if disp.width == maxWidth {
		x = (w - maxWidth) / 2
	}
	ly += disp.Draw(x, ly, s)

	s.Show()
}

func main() {
	defer log.Exit()
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
	disp = NewDisplay()
	disp.Resize(w)
	disp.build()

	quit := make(chan struct{})

	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				stepBasis = littleStep
				if ev.Modifiers()&tcell.ModShift == tcell.ModShift {
					stepBasis = bigStep
				}
				switch ev.Key() {
				case tcell.KeyRune:
					switch ev.Rune() {
					case 'h':
						if ok := disp.ValueDown(stepBasis); ok {
							draw(s)
						}
					case 'k':
						if ok := disp.SectionUp(); ok {
							draw(s)
						}
					case 'j':
						if ok := disp.SectionDown(); ok {
							draw(s)
						}
					case 'l':
						if ok := disp.ValueUp(stepBasis); ok {
							draw(s)
						}
					case 'q':
						close(quit)
						return
					default:
						log.Debugf("ignoring key [%s]", string(ev.Rune()))
					}
				case tcell.KeyRight:
					if ok := disp.ValueUp(stepBasis); ok {
						draw(s)
					}
				case tcell.KeyLeft:
					if ok := disp.ValueDown(stepBasis); ok {
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
				w, h := s.Size()
				log.Debugf("handling resize: w=%04d h=%04d", w, h)
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

	s.Clear()
	s.Fini()
	if err := disp.Close(); err != nil {
		fmt.Println(err)
	}
}

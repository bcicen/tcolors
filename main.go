package main

import (
	"fmt"
	"os"

	"github.com/bcicen/tcolors/logging"
	"github.com/gdamore/tcell"
)

var (
	log   = logging.Init()
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
)

func main() {
	defer log.Exit()

	// initialize screen
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

	// initialize Display
	disp := NewDisplay(s)

	err := disp.Done()
	s.Clear()
	s.Fini()
	if err != nil {
		fmt.Println(err)
	}
}

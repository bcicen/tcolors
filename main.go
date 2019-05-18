package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bcicen/tcolors/logging"
	"github.com/bcicen/tcolors/state"
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

	var (
		printFlag = flag.Bool("p", false, "output current palette contents in common representations")
	)

	flag.Parse()
	if *printFlag {
		tstate := state.NewDefault()
		errExit(tstate.Load())
		fmt.Printf("%s\n\n", tstate.OutputHex())
		fmt.Printf("%s\n\n", tstate.OutputHSV())
		fmt.Printf("%s\n\n", tstate.OutputRGB())
		os.Exit(1)
	}

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

func errExit(err error) {
	if err != nil {
		fmt.Printf("[error]: %s\n", err.Error())
		os.Exit(1)
	}
}

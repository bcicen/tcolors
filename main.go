package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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
		printFlag  = flag.Bool("p", false, "output current palette contents")
		outputFlag = flag.String("o", "all", "color format to output (hex, rgb, hsv, all)")
		fileFlag   = flag.String("f", state.DefaultPalettePath, "specify palette file")
	)

	flag.Parse()
	tstate, err := state.Load(*fileFlag)
	errExit(err)

	if *printFlag {
		cfmt := strings.ToLower(strings.Trim(*outputFlag, " "))
		switch cfmt {
		case "all":
			fmt.Printf("%s\n", tstate.TableString())
		case "hex":
			fmt.Printf("%s\n", tstate.HexString())
		case "hsv":
			fmt.Printf("%s\n", tstate.HSVString())
		case "rgb":
			fmt.Printf("%s\n", tstate.RGBString())
		default:
			errExit(fmt.Errorf("unknown format \"%s\"", cfmt))
		}
		os.Exit(0)
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
	disp := NewDisplay(s, tstate)

	err = disp.Done()
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

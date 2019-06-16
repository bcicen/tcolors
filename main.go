package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bcicen/tcolors/logging"
	"github.com/bcicen/tcolors/state"
	"github.com/fatih/color"
	"github.com/gdamore/tcell"
)

var (
	version = "unknown"
	build   = "dev"
	log     = logging.Init()
	red     = color.New(color.FgRed).SprintFunc()
)

func main() {
	defer log.Exit()

	var (
		printFlag        = flag.Bool("p", false, "output palette contents")
		outputFlag       = flag.String("o", "all", "color format to output (hex, rgb, hsv, all)")
		outputOnExitFlag = flag.Bool("output-on-exit", false, "output palette file contents on exit")
		fileFlag         = flag.String("f", state.DefaultPalettePath, "specify palette file")
		versionFlag      = flag.Bool("v", false, "print version info")
	)

	flag.Parse()

	if *versionFlag {
		fmt.Printf("tcolors v%s %s\n", version, build)
		os.Exit(0)
	}

	tstate, err := state.Load(*fileFlag)
	errExit(err)

	if *printFlag {
		printPalette(tstate, *outputFlag)
		os.Exit(0)
	}

	if *outputOnExitFlag {
		defer printPalette(tstate, *outputFlag)
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

func printPalette(tstate *state.State, cfmt string) {
	cfmt = strings.ToLower(strings.Trim(cfmt, " "))
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
}

func errExit(err error) {
	if err != nil {
		fmt.Printf("%s %s\n", red("err"), err.Error())
		os.Exit(1)
	}
}

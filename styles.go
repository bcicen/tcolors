package main

import (
	"github.com/gdamore/tcell"
)

var (
	defaultSt = tcell.StyleDefault
	blkSt     = defaultSt.
			Background(tcell.ColorBlack).
			Foreground(tcell.ColorBlack)
	indicatorSt = defaultSt.
			Foreground(tcell.NewRGBColor(110, 110, 110)).
			Background(tcell.ColorBlack)
	hiIndicatorSt = defaultSt.
			Foreground(tcell.NewRGBColor(255, 255, 255)).
			Background(tcell.ColorBlack)
	errSt = defaultSt.
		Foreground(tcell.NewRGBColor(255, 000, 043)).
		Background(tcell.ColorBlack)
)

// return bar height for given screen height
func barHeight(h int) int {
	switch {
	case h >= 30:
		return 2
	case h >= 24:
		return 2
	default:
		return 1
	}
}

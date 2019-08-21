package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

var (
	DefaultSt = tcell.StyleDefault
	BlkSt     = DefaultSt.
			Background(tcell.ColorBlack).
			Foreground(tcell.ColorBlack)
	IndicatorSt = DefaultSt.
			Foreground(tcell.NewRGBColor(50, 50, 50)).
			Background(tcell.ColorBlack)
	HiIndicatorSt = DefaultSt.
			Foreground(tcell.NewRGBColor(255, 255, 255)).
			Background(tcell.ColorBlack)
	TextBoxSt = DefaultSt.
			Foreground(tcell.NewRGBColor(160, 160, 160)).
			Background(tcell.ColorBlack)
	ErrSt = DefaultSt.
		Foreground(tcell.NewRGBColor(255, 000, 043)).
		Background(tcell.ColorBlack)
)

// return bar height for given screen height
func barHeight(h int) int {
	switch {
	case h >= 29:
		return 2
	//case h >= 24:
	//return 2
	default:
		return 1
	}
}

func toTColor(c *noire.Color) tcell.Color {
	r, g, b := c.RGB()
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}

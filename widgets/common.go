package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
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

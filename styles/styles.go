package styles

import "github.com/gdamore/tcell"

var (
	Default     = tcell.StyleDefault.Background(tcell.ColorBlack)
	Blank       = Default.Foreground(tcell.ColorBlack)
	Indicator   = Default.Foreground(tcell.NewRGBColor(50, 50, 50))
	IndicatorHi = Default.Foreground(tcell.NewRGBColor(255, 255, 255))
	TextBox     = Default.Foreground(tcell.NewRGBColor(160, 160, 160))
	Error       = Default.Foreground(tcell.NewRGBColor(255, 000, 043))
)

func Load(bg tcell.Color) {
	Default = tcell.StyleDefault.Background(bg)
	Blank = Blank.Foreground(bg)
	Indicator = Indicator.Background(bg)
	IndicatorHi = IndicatorHi.Background(bg)
	TextBox = TextBox.Background(bg)
	Error = Error.Background(bg)
}

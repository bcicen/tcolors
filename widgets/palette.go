package widgets

import (
	"fmt"

	"github.com/bcicen/tcolors/state"
	"github.com/bcicen/tcolors/styles"
	"github.com/gdamore/tcell"
)

const (
	padPalette     = true
	palettePadding = 2
)

type PaletteBox struct {
	width     int
	boxWidth  int
	boxHeight int
	xStretch  int
	pst       tcell.Style // pointer style
	state     *state.State
}

func NewPaletteBox(s *state.State) *PaletteBox {
	pb := &PaletteBox{state: s}
	return pb
}

// Draw redraws p at given coordinates and screen, returning the number
// of rows occupied
func (pb *PaletteBox) Draw(x, y int, s tcell.Screen) int {
	activePaletteHeight := int(float64(pb.boxHeight)*2.5) - 1

	pos := pb.state.Pos()
	items := pb.state.SubColors()
	selected := items[pos] // selected termbox color

	// distribute stretch evenly across boxes
	// where appropriate to facilitate centering
	centerIdx := pb.state.Len() / 2
	boxWidths := make([]int, pb.state.Len())
	boxWidths[centerIdx] = pb.xStretch

	for boxWidths[centerIdx]/3 >= 1 {
		boxWidths[centerIdx] -= 2
		boxWidths[centerIdx-1] += 1
	}

	nextIdx := centerIdx - 1
	for nextIdx >= 0 {
		for boxWidths[nextIdx] >= 2 {
			boxWidths[nextIdx] -= 1
			boxWidths[nextIdx-1] += 1
		}
		nextIdx--
	}
	// mirror first half of array
	for n := len(boxWidths) - 1; n > centerIdx; n-- {
		boxWidths[n] = boxWidths[len(boxWidths)-1-n]
	}
	// apply default boxwidth
	for n := range boxWidths {
		boxWidths[n] += pb.boxWidth
	}

	// text box header
	textBox := []rune(pb.text())
	textBoxX := x + (pb.width-len(textBox))/2

	for col := 0; col < pb.width; col++ {
		s.SetCell(x+col, y, styles.TextBox, '▁')
	}
	y++

	for col := 0; col < pb.width; col++ {
		switch {
		case col == 0:
			s.SetCell(x+col, y, styles.TextBox, '▎')
		case col == pb.width-1:
			s.SetCell(x+col, y, styles.TextBox, '▕')
		case x+col == textBoxX:
			s.SetCell(x+col, y, styles.TextBox, textBox...)
			col += len(textBox) - 1
		}
	}
	y++

	// palette main
	hiSt := styles.IndicatorHi.Background(selected)
	loSt := styles.Indicator.Background(selected)
	topSt := styles.TextBox.Background(selected)
	st := hiSt

	for row := 0; row < activePaletteHeight; row++ {
		for col := 0; col < pb.width; col++ {
			if row == 0 {
				s.SetCell(x+col, y, topSt, '▔')
			} else {
				s.SetCell(x+col, y, st, ' ')
			}
		}
		y++
	}

	lx := x
	for n := range items {
		bw := boxWidths[n]
		if n == pos {
			st = hiSt
		} else {
			st = loSt
		}
		for col := 0; col < bw; col++ {
			s.SetCell(lx+col, y, st, '▁')
		}
		lx += bw
	}
	y++

	lx = x
	cst := styles.Default
	for n, color := range items {
		bw := boxWidths[n]
		cst = cst.Foreground(color)

		switch {
		case padPalette && n == pos:
			st = styles.IndicatorHi
		case n == pos:
			st = styles.IndicatorHi.Background(color)
		case padPalette:
			st = styles.Indicator
		default:
			st = styles.Indicator.Background(color)
		}

		for col := 0; col < bw; col++ {
			for row := 0; row < pb.boxHeight; row++ {
				switch {
				case col == 0:
					s.SetCell(lx, y+row, st, '▎')
				case col == bw-1:
					s.SetCell(lx, y+row, st, '▕')
				case padPalette && row == 0:
					s.SetCell(lx, y+row, cst, '▄')
				case padPalette && row == pb.boxHeight-1:
					s.SetCell(lx, y+row, cst, '▀')
				default:
					s.SetCell(lx, y+row, cst, '█')
				}
			}
			lx++
		}
	}
	y += pb.boxHeight

	lx = x
	for n := range items {
		bw := boxWidths[n]
		if n == pos {
			st = styles.IndicatorHi
		} else {
			st = styles.Indicator
		}
		for col := 0; col < bw; col++ {
			s.SetCell(lx+col, y, st, '▔')
		}
		lx += bw
	}

	return activePaletteHeight + pb.boxHeight + 4
}

func (pb *PaletteBox) text() string {
	const spacer = "  ▎ "

	txt := "▎"
	selected := pb.state.SubColors()[pb.state.Pos()]

	r, g, b := selected.RGB()
	txt = fmt.Sprintf("%03d %03d %03d", r, g, b)

	txt += spacer + "#" + pb.state.Selected().Hex()

	h, s, l := pb.state.Selected().HSL()
	txt += spacer + fmt.Sprintf("%03.0f %03.0f %03.0f", h, s, l)

	return txt
}

func (pb *PaletteBox) Resize(w, h int) {
	pb.boxHeight = barHeight(h) + 1
	pb.boxWidth = w / pb.state.Len()
	pb.width = w
}

func (pb *PaletteBox) Handle(state.Change)            {}
func (pb *PaletteBox) Up(step int)                    { pb.state.Next() }
func (pb *PaletteBox) Down(step int)                  { pb.state.Prev() }
func (pb *PaletteBox) SetPointerStyle(st tcell.Style) { pb.pst = st }

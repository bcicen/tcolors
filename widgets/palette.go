package widgets

import (
	"fmt"

	"github.com/bcicen/tcolors/state"
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
	activePaletteHeight := int(float64(pb.boxHeight) * 2.5)

	pos := pb.state.Pos()
	items := pb.state.SubColors()
	selected := items[pos]

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

	r, g, b := selected.RGB()
	s.SetCell(x+(pb.width-11)/2, y, HiIndicatorSt, []rune(fmt.Sprintf("%03d %03d %03d", r, g, b))...)
	y++

	hiSt := HiIndicatorSt.Background(selected)
	loSt := IndicatorSt.Background(selected)
	st := hiSt

	for row := 0; row < activePaletteHeight; row++ {
		for col := 0; col < pb.width; col++ {
			s.SetCell(x+col, y, st, ' ')
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
	cst := tcell.StyleDefault
	for n, color := range items {
		bw := boxWidths[n]
		cst = cst.Background(tcell.ColorBlack).Foreground(color)

		switch {
		case padPalette && n == pos:
			st = HiIndicatorSt
		case n == pos:
			st = HiIndicatorSt.Background(color)
		case padPalette:
			st = IndicatorSt
		default:
			st = IndicatorSt.Background(color)
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
			st = HiIndicatorSt.Background(tcell.ColorBlack)
		} else {
			st = IndicatorSt.Background(tcell.ColorBlack)
		}
		for col := 0; col < bw; col++ {
			s.SetCell(lx+col, y, st, '▔')
		}
		lx += bw
	}

	return activePaletteHeight + pb.boxHeight + 3
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

package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

var padPalette = false
var DefaultPaletteColor = tcell.NewRGBColor(76, 76, 76)

type PaletteBox struct {
	width    int
	boxWidth int
	pst      tcell.Style // pointer style
	state    *State
}

func NewPaletteBox(s *State) *PaletteBox {
	pb := &PaletteBox{state: s}
	return pb
}

// Draw redraws p at given coordinates and screen, returning the number
// of rows occupied
func (pb *PaletteBox) Draw(x, y int, s tcell.Screen) int {
	pos := pb.state.Pos()
	items := pb.state.SubColors()
	selected := items[pos]

	r, g, b := selected.RGB()
	s.SetCell(x+(pb.width-11)/2, y, hiIndicatorSt, []rune(fmt.Sprintf("%03d %03d %03d", r, g, b))...)
	y++

	hiSt := hiIndicatorSt.Background(selected)
	loSt := indicatorSt.Background(selected)
	st := hiSt

	for row := 0; row < 5; row++ {
		for col := 0; col < pb.width; col++ {
			s.SetCell(x+col, y, st, ' ')
		}
		y++
	}

	lx := x
	for n := range items {
		if n == pos {
			st = hiSt
		} else {
			st = loSt
		}
		for col := 0; col < pb.boxWidth; col++ {
			s.SetCell(lx+col, y, st, '▁')
		}
		lx += pb.boxWidth
	}
	y++

	lx = x
	h := pb.boxWidth / 3
	cst := tcell.StyleDefault
	for n, color := range items {
		cst = cst.Background(tcell.ColorBlack).Foreground(color)

		switch {
		case padPalette && n == pos:
			st = hiIndicatorSt
		case n == pos:
			st = hiIndicatorSt.Background(color)
		case padPalette:
			st = indicatorSt
		default:
			st = indicatorSt.Background(color)
		}

		for col := 0; col < pb.boxWidth; col++ {
			for row := 0; row < h; row++ {
				switch {
				case col == 0:
					s.SetCell(lx, y+row, st, '▎')
				case col == pb.boxWidth-1:
					s.SetCell(lx, y+row, st, '▕')
				case padPalette && row == 0:
					s.SetCell(lx, y+row, cst, '▄')
				case padPalette && row == h-1:
					s.SetCell(lx, y+row, cst, '▀')
				default:
					s.SetCell(lx, y+row, cst, '█')
				}
			}
			lx++
		}
	}
	y += h

	lx = x
	for n := range items {
		if n == pos {
			st = hiIndicatorSt.Background(tcell.ColorBlack)
		} else {
			st = indicatorSt.Background(tcell.ColorBlack)
		}
		for col := 0; col < pb.boxWidth; col++ {
			//s.SetCell(lx+col, y, st, []rune(fmt.Sprintf("%d", n))...)
			s.SetCell(lx+col, y, st, '▔')
		}
		lx += pb.boxWidth
	}

	return 10
}

func (pb *PaletteBox) Handle(change StateChange) {
	if change == NoChange {
		return
	}
	nc := noire.NewHSV(pb.state.Hue(), pb.state.Saturation(), pb.state.Value())
	pb.state.SetSelected(toTColor(nc))
}

func (pb *PaletteBox) Resize(w int) {
	pb.boxWidth = (w / 2) / pb.state.Len()
	pb.width = pb.boxWidth * pb.state.Len()
}

func (pb *PaletteBox) Up(step int)   { pb.state.Next() }
func (pb *PaletteBox) Down(step int) { pb.state.Prev() }

func (pb *PaletteBox) SetPointerStyle(st tcell.Style) { pb.pst = st }
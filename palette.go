package main

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell"
)

var padPalette = false
var DefaultPaletteColor = tcell.NewRGBColor(76, 76, 76)

type PaletteBox struct {
	items    [8]tcell.Color // navigation colors
	pos      int
	width    int
	boxWidth int
	pst      tcell.Style // pointer style
	lock     sync.RWMutex
}

func NewPaletteBox(width int) *PaletteBox {
	pb := &PaletteBox{width: width}
	for n := range pb.items {
		pb.items[n] = DefaultPaletteColor
	}
	return pb
}

func (pb *PaletteBox) SetPos(n int)          { pb.pos = n }
func (pb *PaletteBox) Selected() tcell.Color { return pb.items[pb.pos] }

// Draw redraws p at given coordinates and screen, returning the number
// of rows occupied
func (pb *PaletteBox) Draw(x, y int, s tcell.Screen) int {
	r, g, b := pb.Selected().RGB()
	s.SetCell(x+(pb.width-11)/2, y, hiIndicatorSt, []rune(fmt.Sprintf("%03d %03d %03d", r, g, b))...)

	hiSt := hiIndicatorSt.Background(pb.Selected())
	loSt := indicatorSt.Background(pb.Selected())
	st := hiSt

	for row := 0; row < 5; row++ {
		for col := 0; col < pb.width; col++ {
			s.SetCell(x+col, y, st, ' ')
		}
		y++
	}

	lx := x
	for n := range pb.items {
		if n == pb.pos {
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
	for n, color := range pb.items {
		cst = cst.Background(tcell.ColorBlack).Foreground(color)

		switch {
		case padPalette && n == pb.pos:
			st = hiIndicatorSt
		case n == pb.pos:
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
	for n := range pb.items {
		if n == pb.pos {
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

	return 9
}

func (pb *PaletteBox) Resize(w int) {
	pb.boxWidth = (w / 2) / len(pb.items)
	pb.width = pb.boxWidth * len(pb.items)
}

func (pb *PaletteBox) Update(c tcell.Color) {
	pb.lock.Lock()
	defer pb.lock.Unlock()
	pb.items[pb.pos] = c
}

func (pb *PaletteBox) Up(step int) {
	pb.lock.Lock()
	defer pb.lock.Unlock()
	pb.pos++
	if pb.pos >= len(pb.items) {
		pb.pos = 0
	}
}

func (pb *PaletteBox) Down(step int) {
	pb.lock.Lock()
	defer pb.lock.Unlock()
	pb.pos--
	if pb.pos < 0 {
		pb.pos = len(pb.items) - 1
	}
}

func (pb *PaletteBox) SetPointerStyle(st tcell.Style) { pb.pst = st }

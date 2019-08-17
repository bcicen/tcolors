package widgets

import (
	"github.com/gdamore/tcell"
)

type MenuFn func(tcell.Screen) MenuFn

func HelpMenu(s tcell.Screen) MenuFn {
	drawHelpMenu(s)
	for {
		ev := s.PollEvent()
		switch ev.(type) {
		case *tcell.EventKey:
			return nil
		case *tcell.EventResize:
			return HelpMenu
		}
	}
	return nil
}

func drawHelpMenu(s tcell.Screen) {
	var maxL, maxR, menuW int
	for _, item := range helpMenuItems {
		if len(item.key) > maxL {
			maxL = len(item.key)
		}
		if len(item.desc) > maxR {
			maxR = len(item.desc)
		}
	}
	menuW = maxL + maxR + 4

	w, _ := s.Size()

	x := (w - menuW) / 2
	y := 2

	for n, item := range helpMenuItems {
		s.SetCell(x+1, y+n, DefaultSt, []rune(item.key)...)
		s.SetCell(x+maxL+2, y+n, DefaultSt, '|')
		s.SetCell(x+maxL+4, y+n, DefaultSt, []rune(item.desc)...)
	}

	s.Show()
}

type helpMenuItem struct {
	key  string
	desc string
}

var helpMenuItems = []helpMenuItem{
	{"↑, k", "navigate up"},
	{"↓, j", "navigate down"},
	{"←, h", "decrease selected value"},
	{"→, l", "increase selected value"},
	{"<shift> + ←/→/h/l", "more quickly increase/decrease selected value"},
	{"a, <ins>", "add a new palette color"},
	{"x, <del>", "remove the selected palette color"},
	{"q, <esc>", "exit tcolors"},
	{"?", "show this help menu"},
}

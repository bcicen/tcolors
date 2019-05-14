package main

import (
	"time"

	"github.com/gdamore/tcell"
)

const errDisplayTimeout = 2 * time.Second

type ErrorMsg struct {
	ts    time.Time
	text  string
	width int
}

func NewErrorMsg() *ErrorMsg {
	return &ErrorMsg{}
}

// Draw redraws msg at given coordinates and screen, returning the number
// of rows occupied
func (msg *ErrorMsg) Draw(x int, s tcell.Screen) int {
	_, h := s.Size()
	y := h - 2

	if len(msg.text) > 0 && time.Since(msg.ts) >= errDisplayTimeout {
		// clear error message
		for i := x; i <= msg.width; i++ {
			s.SetCell(i, y, errSt, ' ')
		}
		msg.text = ""
	} else {
		for n, ch := range msg.text {
			if x+n > msg.width {
				break
			}
			s.SetCell(x+n, y, errSt, ch)
		}
	}

	return 0
}

func (msg *ErrorMsg) Set(s string) {
	msg.ts = time.Now()
	msg.text = s
}

func (msg *ErrorMsg) Resize(w int) { msg.width = w }

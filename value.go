package main

import (
	"fmt"
	"github.com/teacat/noire"
)

const (
	valMin   = 0.0
	valMax   = 100.0
	valIncr  = 0.5
	valCount = int(valMax/valIncr) + 1
)

type ValueBar struct {
	*NavBar
	scale [valCount]float64
}

func NewValueBar(s *State) *ValueBar {
	bar := &ValueBar{NavBar: NewNavBar(s, valCount)}

	i := valMin
	for n, _ := range bar.scale {
		bar.scale[n] = i
		i += 0.5
	}
	bar.scale[0] = 0.01

	return bar
}

// State change handler
func (bar *ValueBar) Handle(change StateChange) {
	var nc *noire.Color

	if change.Includes(HueChanged, SaturationChanged) {
		for n, val := range bar.scale {
			nc = noire.NewHSV(bar.state.Hue(), bar.state.Saturation(), val)
			bar.items[n] = toTColor(nc)
		}
	}

	if change.Includes(ValueChanged) {
		bar.SetPos(roundFloat(bar.state.Value() / valIncr))
		bar.SetLabel(fmt.Sprintf("%5.1f ", bar.scale[bar.pos]))
	}

}

func (bar *ValueBar) Up(step int) {
	bar.up(step)
	bar.setState()
}

func (bar *ValueBar) Down(step int) {
	bar.down(step)
	bar.setState()
}

func (bar *ValueBar) setState() {
	bar.state.SetValue(bar.scale[bar.pos])
}

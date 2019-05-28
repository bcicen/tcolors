package state

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

type subState struct {
	rgb        [3]int32
	hue        float64
	saturation float64
	value      float64
}

func (ss *subState) NColor() *noire.Color {
	return noire.NewHSV(ss.hue, ss.saturation, ss.value)
}

func (ss *subState) TColor() tcell.Color {
	r, g, b := ss.NColor().RGB()
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}

func (ss *subState) HexString() string {
	return fmt.Sprintf("%06x", ss.TColor().Hex())
}

func (ss *subState) HSVString() string {
	return fmt.Sprintf("%03.0f %03.0f %03.0f", ss.hue, ss.saturation, ss.value)
}

func (ss *subState) RGBString() string {
	return fmt.Sprintf("%03d %03d %03d", ss.rgb[0], ss.rgb[1], ss.rgb[2])
}

func (ss *subState) SetNColor(nc *noire.Color) {
	r, g, b := nc.RGB()
	h, s, v := nc.HSV()
	ss.rgb[0] = int32(r)
	ss.rgb[1] = int32(g)
	ss.rgb[2] = int32(b)
	ss.hue = h
	ss.saturation = s
	ss.value = v
}

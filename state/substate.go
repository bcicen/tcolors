package state

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

type subState struct {
	*noire.Color
}

func newDefaultSubState() *subState {
	return &subState{noire.NewRGB(255, 255, 255)}
}

func (ss *subState) NColor() *noire.Color {
	return ss.Color
}

func (ss *subState) TColor() tcell.Color {
	r, g, b := ss.RGB()
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}

// PColor returns the current substate in a paletteColor suitable for
// saving to config state file
func (ss *subState) PColor() (pc paletteColor) {
	r, g, b := ss.RGB()
	h, s, v := ss.HSV()
	pc.RGB = []int{int(r), int(g), int(b)}
	pc.HEX = ss.HexString()
	pc.HSV = []float64{h, s, v}
	return pc
}

func (ss *subState) HSVIdx(n uint8) float64 {
	h, s, v := ss.HSV()
	switch n {
	case 0:
		return h
	case 1:
		return s
	case 2:
		return v
	default:
		panic(fmt.Sprintf("bad index selector: %d", n))
	}
}

// SetNColor replaces the color for the current subState
func (ss *subState) SetNColor(nc *noire.Color) {
	ss.Color = nc
}

// SetHue sets the HSV hue for the current subState
func (ss *subState) SetHue(n float64) {
	_, s, v := ss.HSV()
	ss.Color = noire.NewHSV(n, s, v)
}

// SetSaturation sets the HSV saturation for the current subState
func (ss *subState) SetSaturation(n float64) {
	h, _, v := ss.HSV()
	ss.Color = noire.NewHSV(h, n, v)
}

// SetValue sets the HSV value for the current subState
func (ss *subState) SetValue(n float64) {
	h, s, _ := ss.HSV()
	ss.Color = noire.NewHSV(h, s, n)
}

func (ss *subState) HexString() string {
	return ss.Hex()
}

func (ss *subState) HSVString() string {
	h, s, v := ss.HSV()
	return fmt.Sprintf("%03.0f %03.0f %03.0f", h, s, v)
}

func (ss *subState) RGBString() string {
	r, g, b := ss.RGB()
	return fmt.Sprintf("%03.0f %03.0f %03.0f", r, g, b)
}

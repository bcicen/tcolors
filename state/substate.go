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
	return tcell.NewRGBColor(ss.rgb[0], ss.rgb[1], ss.rgb[2])
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

func (ss *subState) load(b []byte) error {
	var offset int

	if len(b) != subStateByteSize {
		return malformedErr
	}

	ss.rgb[0] = int32FromBytes(b[offset : offset+4])
	offset += 4
	ss.rgb[1] = int32FromBytes(b[offset : offset+4])
	offset += 4
	ss.rgb[2] = int32FromBytes(b[offset : offset+4])
	offset += 4

	ss.hue = float64FromBytes(b[offset : offset+8])
	offset += 8
	ss.saturation = float64FromBytes(b[offset : offset+8])
	offset += 8
	ss.value = float64FromBytes(b[offset : offset+8])
	offset += 8

	return nil
}

func (ss *subState) bytes() []byte {
	var buf [subStateByteSize]byte
	var offset int
	for _, n := range ss.rgb {
		offset += writeInt32Bytes(buf[offset:], n)
	}
	offset += writeFloat64Bytes(buf[offset:], ss.hue)
	offset += writeFloat64Bytes(buf[offset:], ss.saturation)
	offset += writeFloat64Bytes(buf[offset:], ss.value)
	return buf[:]
}

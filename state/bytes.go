package state

import (
	"encoding/binary"
	"math"
)

func writeInt32Bytes(buf []byte, n int32) int {
	binary.BigEndian.PutUint32(buf, uint32(n))
	return 4
}

func int32FromBytes(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

func float64FromBytes(buf []byte) float64 {
	bits := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(bits)
}

func writeFloat64Bytes(buf []byte, float float64) int {
	bits := math.Float64bits(float)
	binary.LittleEndian.PutUint64(buf[:], bits)
	return 8
}

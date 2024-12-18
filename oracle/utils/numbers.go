package utils

import (
	"encoding/binary"
	"github.com/holiman/uint256"
	"math"
)

func ApplyBPS(amount *uint256.Int, bps int) *uint256.Int {
	result := uint256.NewInt(0)

	result.Mul(amount, uint256.NewInt(uint64(bps)))
	result.Div(result, uint256.NewInt(10_000))

	return result
}

func ApplyReverseBPS(amount *uint256.Int, bps int) *uint256.Int {
	return ApplyBPS(amount, 10_000-bps)
}

func RoundDown(val float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Floor(val*factor) / factor
}

func MoveDecimals(amount *uint256.Int, baseDecimal int, targetDecimal int) *uint256.Int {
	result := uint256.NewInt(0)

	result.Mul(amount, uint256.NewInt(uint64(math.Pow10(targetDecimal))))
	result.Div(result, uint256.NewInt(uint64(math.Pow10(baseDecimal))))

	return result
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func ToLittleEndianBytes(number any) []byte {
	switch v := number.(type) {
	case uint64:
		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, v)
		return bytes
	case uint32:
		bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, v)
		return bytes
	case uint16:
		bytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(bytes, v)
		return bytes
	default:
		panic("Unsupported numeric type")
	}
}

package utils

import (
	"github.com/holiman/uint256"
	"gringotts/models"
	"math"
)

func ApplyBPS(amount *uint256.Int, bps models.BPS) *uint256.Int {
	result := uint256.NewInt(0)

	result.Mul(amount, uint256.NewInt(uint64(bps)))
	result.Div(result, uint256.NewInt(10_000))

	return result
}

func ApplyReverseBPS(amount *uint256.Int, bps models.BPS) *uint256.Int {
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

func Max(a uint16, b uint16) uint16 {
	if a > b {
		return a
	}
	return b
}

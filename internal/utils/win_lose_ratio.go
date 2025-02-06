package utils

import (
	"fmt"
	"math"
)

func CalculateWinLoseRatio(win, loses int) string {
	percent := (float64(win) / float64(win+loses)) * 100
	rounded := math.Round(percent*100) / 100
	return fmt.Sprintf("%.2f%%", rounded)
}

package colormetric

import (
	"math"

	"colorsci/internal/colorconv"
	"colorsci/internal/colorparse"
)

// HarmonyScore evaluates how harmonious a set of colors is based on their
// hue distribution. The score ranges from 0 (dissonant) to 1 (perfectly
// harmonious). A single color or empty set returns 1.
func HarmonyScore(colors []colorparse.Color) float64 {
	n := len(colors)
	if n <= 1 {
		return 1.0
	}

	hues := make([]float64, n)
	for i, c := range colors {
		h, _, _ := colorconv.RGBToHSL(c.R, c.G, c.B)
		hues[i] = h
	}

	// 计算色相对之间与理想间距的吻合度
	idealSpacing := 360.0 / float64(n)
	totalDeviation := 0.0

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			diff := math.Abs(hues[i] - hues[j])
			if diff > 180 {
				diff = 360 - diff
			}
			// 计算与最近理想间距倍数的偏差
			nearest := math.Round(diff/idealSpacing) * idealSpacing
			if nearest == 0 {
				nearest = idealSpacing
			}
			deviation := math.Abs(diff - nearest)
			totalDeviation += deviation
		}
	}

	pairs := float64(n * (n - 1) / 2)
	avgDeviation := totalDeviation / pairs
	// 将偏差映射到 [0,1] 分数
	score := 1.0 - math.Min(avgDeviation/90.0, 1.0)
	return score
}

// ComplementaryHarmony scores how close two colors are to being
// complementary (180° apart). Returns 1 for perfect complements, 0 for
// identical hues.
func ComplementaryHarmony(c1, c2 colorparse.Color) float64 {
	h1, _, _ := colorconv.RGBToHSL(c1.R, c1.G, c1.B)
	h2, _, _ := colorconv.RGBToHSL(c2.R, c2.G, c2.B)
	diff := math.Abs(h1 - h2)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff / 180.0
}

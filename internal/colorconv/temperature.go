package colorconv

import (
	"math"

	"colorsci/internal/colorparse"
)

// Temperature returns the perceived color temperature of a color as a
// value in [-1, 1]. Negative values indicate cool colors (blue-green),
// positive values indicate warm colors (red-orange-yellow), and values
// near zero indicate neutral colors.
func Temperature(c colorparse.Color) float64 {
	h, s, _ := RGBToHSL(c.R, c.G, c.B)
	if s < 0.05 {
		// 低饱和度视为中性
		return 0
	}
	// 暖色: 0-90° (红-橙-黄), 冷色: 150-270° (青-蓝-紫)
	// 中间过渡: 90-150° 和 270-360°
	var temp float64
	switch {
	case h <= 90:
		// 暖色区，0° 最暖（红），90° 过渡到中性
		temp = 1.0 - h/90.0
	case h <= 150:
		// 从暖过渡到冷
		temp = -(h - 90) / 60.0
	case h <= 270:
		// 冷色区
		temp = -1.0
	default:
		// 从冷过渡回暖 (270-360°)
		temp = (h - 270) / 90.0
	}
	// 饱和度越高温度感越强
	return temp * math.Min(s*1.5, 1.0)
}

// IsWarm reports whether the color is perceived as warm (temperature > 0.2).
func IsWarm(c colorparse.Color) bool {
	return Temperature(c) > 0.2
}

// IsCool reports whether the color is perceived as cool (temperature < -0.2).
func IsCool(c colorparse.Color) bool {
	return Temperature(c) < -0.2
}

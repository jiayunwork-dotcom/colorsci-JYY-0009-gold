// Example: parse colors, convert to Lab and compute a DeltaE2000 score.
package main

import (
	"fmt"

	"colorsci/internal/colorconv"
	"colorsci/internal/colormetric"
	"colorsci/internal/colorparse"
)

func main() {
	colors := []string{
		"#ff0000",
		"hsl(120 100% 50%)",
		"rgb(0 0 255 / 0.5)",
		"rebeccapurple",
		"#f80",
	}
	for _, s := range colors {
		c, err := colorparse.Parse(s)
		if err != nil {
			fmt.Printf("%-18s parse error: %v\n", s, err)
			continue
		}
		l, a, b := colorconv.ColorToLab(c)
		fmt.Printf("%-18s -> %-9s Lab(%.2f, %.2f, %.2f) alpha=%.2f\n",
			s, colorparse.FormatHex(c), l, a, b, c.A)
	}

	red, _ := colorparse.Parse("red")
	green, _ := colorparse.Parse("lime")
	white, _ := colorparse.Parse("white")
	d := colormetric.DeltaE2000(colormetric.LabOf(red), colormetric.LabOf(green))
	ratio := colormetric.ContrastRatio(white, red)
	fmt.Printf("DeltaE2000(red, green) = %.4f\n", d)
	fmt.Printf("contrast(white, red)   = %.3f:1\n", ratio)
}

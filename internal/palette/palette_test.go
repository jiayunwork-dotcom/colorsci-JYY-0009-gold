package palette

import (
	"colorsci/internal/colorparse"
	"testing"
)

func TestAnalogous(t *testing.T) {
	base := colorparse.Color{R: 255, G: 0, B: 0}
	colors := Analogous(base, 5, 60)
	if len(colors) != 5 {
		t.Fatalf("expected 5, got %d", len(colors))
	}
}

func TestComplementary(t *testing.T) {
	red := colorparse.Color{R: 255, G: 0, B: 0}
	comp := Complementary(red)
	// Complement of red should be cyan-ish.
	if comp.R > comp.G {
		t.Fatalf("complement of red should have more green/blue: %+v", comp)
	}
}

func TestTriadic(t *testing.T) {
	base := colorparse.Color{R: 255, G: 0, B: 0}
	tri := Triadic(base)
	// Three distinct colors.
	if tri[0] == tri[1] || tri[1] == tri[2] {
		t.Fatal("triadic should produce distinct colors")
	}
}

func TestMonochromatic(t *testing.T) {
	base := colorparse.Color{R: 100, G: 150, B: 200}
	mono := Monochromatic(base, 5)
	if len(mono) != 5 {
		t.Fatalf("expected 5, got %d", len(mono))
	}
}

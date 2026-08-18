package colormetric

import (
	"math"
	"testing"

	"colorsci/internal/colorparse"
)

func mustParse(t *testing.T, s string) colorparse.Color {
	t.Helper()
	c, err := colorparse.Parse(s)
	if err != nil {
		t.Fatalf("Parse(%q): %v", s, err)
	}
	return c
}

func TestRelativeLuminanceBlackWhite(t *testing.T) {
	black := mustParse(t, "#000000")
	white := mustParse(t, "#ffffff")
	if got := RelativeLuminance(black); got != 0 {
		t.Errorf("black luminance = %v, want 0", got)
	}
	if got := RelativeLuminance(white); math.Abs(got-1) > 1e-9 {
		t.Errorf("white luminance = %v, want 1", got)
	}
}

func TestRelativeLuminanceKnown(t *testing.T) {
	gray := mustParse(t, "#808080")
	// 128/255 linearized: ~0.21586
	if got := RelativeLuminance(gray); math.Abs(got-0.21586) > 1e-4 {
		t.Errorf("gray luminance = %v, want ~0.21586", got)
	}
}

func TestContrastRatioBlackWhite(t *testing.T) {
	black := mustParse(t, "#000000")
	white := mustParse(t, "#ffffff")
	if got := ContrastRatio(black, white); math.Abs(got-21) > 1e-6 {
		t.Errorf("black/white contrast = %v, want 21", got)
	}
}

func TestContrastRatioSymmetric(t *testing.T) {
	a := mustParse(t, "#ff0000")
	b := mustParse(t, "#00ff00")
	if d1, d2 := ContrastRatio(a, b), ContrastRatio(b, a); math.Abs(d1-d2) > 1e-12 {
		t.Errorf("contrast not symmetric: %v vs %v", d1, d2)
	}
}

func TestContrastRatioAAThreshold(t *testing.T) {
	// WCAG example: #767676 on white is just above 4.5:1.
	fg := mustParse(t, "#767676")
	bg := mustParse(t, "#ffffff")
	r := ContrastRatio(fg, bg)
	if r < 4.5 {
		t.Errorf("#767676 on white = %.3f, want >= 4.5", r)
	}
}

func TestWCAGLevels(t *testing.T) {
	lv := WCAGLevels(4.5)
	if !lv.AANormal || !lv.AALarge || !lv.AAALarge || lv.AAANormal {
		t.Errorf("levels for 4.5 = %+v", lv)
	}
	lv2 := WCAGLevels(7.0)
	if !lv2.AAANormal {
		t.Errorf("7.0 should satisfy AAA normal: %+v", lv2)
	}
	lv3 := WCAGLevels(2.9)
	if lv3.AALarge || lv3.AANormal {
		t.Errorf("2.9 should satisfy nothing: %+v", lv3)
	}
}

func TestLevelsPasses(t *testing.T) {
	lv := WCAGLevels(4.5)
	if !lv.Passes(false, false) {
		t.Errorf("4.5 should pass AA normal")
	}
	if lv.Passes(false, true) {
		t.Errorf("4.5 should not pass AAA normal")
	}
	lvAaa := WCAGLevels(7.1)
	if !lvAaa.Passes(false, true) {
		t.Errorf("7.1 should pass AAA normal")
	}
}

func TestBestPair(t *testing.T) {
	fg := []colorparse.Color{mustParse(t, "#333333"), mustParse(t, "#ff0000")}
	bg := []colorparse.Color{mustParse(t, "#ffffff")}
	bf, bb, r := BestPair(fg, bg)
	if r < 4.5 {
		t.Errorf("best pair ratio = %.2f, want >= 4.5", r)
	}
	if bb != mustParse(t, "#ffffff") {
		t.Errorf("best background = %v, want white", bb)
	}
	if bf == mustParse(t, "#ff0000") {
		t.Errorf("red on white should not be the best pair")
	}
}

func TestBestPairEmpty(t *testing.T) {
	if _, _, r := BestPair(nil, nil); !math.IsInf(r, -1) {
		t.Errorf("empty BestPair ratio = %v, want -Inf", r)
	}
}

func TestLabOf(t *testing.T) {
	c := mustParse(t, "#ff0000")
	l := LabOf(c)
	if math.Abs(l.L-53.233) > 0.02 {
		t.Errorf("LabOf(red).L = %v, want ~53.233", l.L)
	}
}

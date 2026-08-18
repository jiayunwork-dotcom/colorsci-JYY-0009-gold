package colorparse

import (
	"errors"
	"testing"
)

func TestNamedCommon(t *testing.T) {
	cases := map[string]Color{
		"red":           {R: 0xff, G: 0x00, B: 0x00, A: 1},
		"white":         {R: 0xff, G: 0xff, B: 0xff, A: 1},
		"black":         {R: 0x00, G: 0x00, B: 0x00, A: 1},
		"navy":          {R: 0x00, G: 0x00, B: 0x80, A: 1},
		"aqua":          {R: 0x00, G: 0xff, B: 0xff, A: 1},
		"fuchsia":       {R: 0xff, G: 0x00, B: 0xff, A: 1},
		"rebeccapurple": {R: 0x66, G: 0x33, B: 0x99, A: 1},
	}
	for name, want := range cases {
		c, err := ParseNamed(name)
		if err != nil {
			t.Errorf("ParseNamed(%q) error: %v", name, err)
			continue
		}
		if c != want {
			t.Errorf("ParseNamed(%q) = %v, want %v", name, c, want)
		}
	}
}

func TestNamedCaseInsensitive(t *testing.T) {
	c, err := ParseNamed("CornFlowerBlue")
	if err != nil {
		t.Fatalf("ParseNamed(CornFlowerBlue) error: %v", err)
	}
	if c.R != 0x64 || c.G != 0x95 || c.B != 0xed {
		t.Errorf("cornflowerblue = %v", c)
	}
}

func TestNamedUnknown(t *testing.T) {
	if _, err := ParseNamed("notacolor"); !errors.Is(err, ErrUnknown) {
		t.Errorf("err = %v, want ErrUnknown", err)
	}
}

func TestNamedThroughParse(t *testing.T) {
	c, err := Parse("tomato")
	if err != nil {
		t.Fatalf("Parse(tomato) error: %v", err)
	}
	if c.R != 0xff || c.G != 0x63 || c.B != 0x47 {
		t.Errorf("tomato = %v", c)
	}
}

func TestNamedCount(t *testing.T) {
	// The W3C extended keyword set has 148 entries including grey/gray
	// aliases; require a full table.
	if got := len(named); got < 140 {
		t.Errorf("named table has %d entries, want >= 140", got)
	}
}

func TestNamesSortedUnique(t *testing.T) {
	names := Names()
	if len(names) != len(named) {
		t.Errorf("Names() returned %d, table has %d", len(names), len(named))
	}
	for i := 1; i < len(names); i++ {
		if names[i-1] >= names[i] {
			t.Errorf("Names() not sorted at %d: %q >= %q", i, names[i-1], names[i])
		}
	}
}

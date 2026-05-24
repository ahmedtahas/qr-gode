package colors

import "testing"

func TestNewSolid(t *testing.T) {
	s := NewSolid("#ff0000")
	if s.Hex != "#ff0000" {
		t.Errorf("expected #ff0000, got %s", s.Hex)
	}
}

func TestSolid_ColorAt(t *testing.T) {
	s := NewSolid("#abcdef")
	// Solid color ignores coordinates.
	for _, xy := range [][2]float64{{0, 0}, {0.5, 0.5}, {1, 1}, {-1, 2}} {
		if got := s.ColorAt(xy[0], xy[1]); got != "#abcdef" {
			t.Errorf("ColorAt(%v,%v) = %s, want #abcdef", xy[0], xy[1], got)
		}
	}
}

func TestSolid_Type(t *testing.T) {
	if got := NewSolid("#000").Type(); got != "solid" {
		t.Errorf("Type = %s, want solid", got)
	}
}

func TestSolid_SVGDefs(t *testing.T) {
	if got := NewSolid("#000").SVGDefs("foo"); got != "" {
		t.Errorf("SVGDefs = %q, want empty", got)
	}
}

func TestSolid_SVGFill(t *testing.T) {
	if got := NewSolid("#112233").SVGFill("ignored"); got != "#112233" {
		t.Errorf("SVGFill = %s, want #112233", got)
	}
}

func TestParseHex(t *testing.T) {
	cases := []struct {
		in      string
		wantErr bool
	}{
		{"#abc", false},
		{"#aabbcc", false},
		{"#aabbccdd", false},
		{"", true},
		{"abc", true},      // missing #
		{"#ab", true},      // wrong length
		{"#abcde", true},   // wrong length
		{"#abcdefg", true}, // wrong length
	}
	for _, tc := range cases {
		got, err := ParseHex(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseHex(%q) expected error, got %q", tc.in, got)
			}
		} else {
			if err != nil {
				t.Errorf("ParseHex(%q) unexpected error: %v", tc.in, err)
			}
			if got != tc.in {
				t.Errorf("ParseHex(%q) = %q, want %q", tc.in, got, tc.in)
			}
		}
	}
}

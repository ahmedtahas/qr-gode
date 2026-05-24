package colors

import (
	"strings"
	"testing"
)

func TestNewLinearGradient(t *testing.T) {
	g := NewLinearGradient(45, []string{"#000", "#fff"})
	if g.Angle != 45 {
		t.Errorf("Angle = %v, want 45", g.Angle)
	}
	if len(g.Stops) != 2 {
		t.Errorf("Stops len = %d, want 2", len(g.Stops))
	}
}

func TestLinearGradient_ColorAt_EmptyStops(t *testing.T) {
	g := NewLinearGradient(0, nil)
	if got := g.ColorAt(0.5, 0.5); got != "#000000" {
		t.Errorf("empty stops fallback = %s, want #000000", got)
	}
}

func TestLinearGradient_ColorAt_SingleStop(t *testing.T) {
	g := NewLinearGradient(0, []string{"#abcdef"})
	if got := g.ColorAt(0.3, 0.7); got != "#abcdef" {
		t.Errorf("single stop = %s, want #abcdef", got)
	}
}

func TestLinearGradient_ColorAt_Interpolation(t *testing.T) {
	g := NewLinearGradient(0, []string{"#aaa", "#bbb", "#ccc"})
	// At angle 0, pos = x. Index 0 at x=0, last index at x=1.
	if got := g.ColorAt(0, 0); got != "#aaa" {
		t.Errorf("ColorAt(0,0) = %s, want #aaa", got)
	}
	if got := g.ColorAt(1, 0); got != "#ccc" {
		t.Errorf("ColorAt(1,0) = %s, want #ccc", got)
	}
}

func TestLinearGradient_ColorAt_Clamps(t *testing.T) {
	g := NewLinearGradient(0, []string{"#000", "#fff"})
	// Out-of-range coordinates should clamp to a valid stop, not panic.
	_ = g.ColorAt(-5, -5)
	_ = g.ColorAt(5, 5)
}

func TestLinearGradient_Type(t *testing.T) {
	if got := NewLinearGradient(0, nil).Type(); got != "linear-gradient" {
		t.Errorf("Type = %s", got)
	}
}

func TestLinearGradient_SVGDefs(t *testing.T) {
	g := NewLinearGradient(45, []string{"#ff0000", "#00ff00", "#0000ff"})
	out := g.SVGDefs("grad-1")
	if !strings.HasPrefix(out, `<linearGradient id="grad-1"`) {
		t.Errorf("missing opening tag: %s", out)
	}
	if !strings.Contains(out, "</linearGradient>") {
		t.Errorf("missing closing tag: %s", out)
	}
	for _, hex := range []string{"#ff0000", "#00ff00", "#0000ff"} {
		if !strings.Contains(out, hex) {
			t.Errorf("SVGDefs missing stop color %s: %s", hex, out)
		}
	}
}

func TestLinearGradient_SVGFill(t *testing.T) {
	if got := NewLinearGradient(0, nil).SVGFill("g1"); got != "url(#g1)" {
		t.Errorf("SVGFill = %s", got)
	}
}

func TestNewRadialGradient(t *testing.T) {
	g := NewRadialGradient(0.4, 0.6, []string{"#a", "#b"})
	if g.CenterX != 0.4 || g.CenterY != 0.6 {
		t.Errorf("center = (%v,%v), want (0.4,0.6)", g.CenterX, g.CenterY)
	}
	if len(g.Stops) != 2 {
		t.Errorf("Stops len = %d", len(g.Stops))
	}
}

func TestRadialGradient_ColorAt_EmptyStops(t *testing.T) {
	g := NewRadialGradient(0.5, 0.5, nil)
	if got := g.ColorAt(0.5, 0.5); got != "#000000" {
		t.Errorf("fallback = %s", got)
	}
}

func TestRadialGradient_ColorAt_CenterAndEdge(t *testing.T) {
	g := NewRadialGradient(0.5, 0.5, []string{"#center", "#mid", "#edge"})
	// At the center, distance = 0 → first stop.
	if got := g.ColorAt(0.5, 0.5); got != "#center" {
		t.Errorf("center = %s, want #center", got)
	}
	// Far corner → clamps to last stop.
	if got := g.ColorAt(1, 1); got != "#edge" {
		t.Errorf("corner = %s, want #edge", got)
	}
}

func TestRadialGradient_Type(t *testing.T) {
	if got := NewRadialGradient(0, 0, nil).Type(); got != "radial-gradient" {
		t.Errorf("Type = %s", got)
	}
}

func TestRadialGradient_SVGDefs(t *testing.T) {
	g := NewRadialGradient(0.25, 0.75, []string{"#aaa", "#bbb"})
	out := g.SVGDefs("rg")
	if !strings.HasPrefix(out, `<radialGradient id="rg"`) {
		t.Errorf("missing opening tag: %s", out)
	}
	if !strings.Contains(out, "</radialGradient>") {
		t.Errorf("missing closing tag: %s", out)
	}
	if !strings.Contains(out, `cx="25%"`) || !strings.Contains(out, `cy="75%"`) {
		t.Errorf("center coords missing: %s", out)
	}
}

func TestRadialGradient_SVGFill(t *testing.T) {
	if got := NewRadialGradient(0, 0, nil).SVGFill("rg"); got != "url(#rg)" {
		t.Errorf("SVGFill = %s", got)
	}
}

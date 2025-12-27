package colors

import (
	"fmt"
	"math"
	"strings"
)

// LinearGradient represents a linear color gradient.
type LinearGradient struct {
	Angle float64  // Angle in degrees (0 = right, 90 = down, 180 = left, 270 = up)
	Stops []string // Color stops (hex values)
}

// NewLinearGradient creates a linear gradient.
func NewLinearGradient(angle float64, stops []string) *LinearGradient {
	return &LinearGradient{
		Angle: angle,
		Stops: stops,
	}
}

func (g *LinearGradient) ColorAt(x, y float64) string {
	if len(g.Stops) == 0 {
		return "#000000"
	}
	if len(g.Stops) == 1 {
		return g.Stops[0]
	}
	// Simple linear interpolation along gradient direction
	rad := g.Angle * math.Pi / 180
	pos := x*math.Cos(rad) + y*math.Sin(rad)
	pos = math.Max(0, math.Min(1, pos))
	idx := int(pos * float64(len(g.Stops)-1))
	return g.Stops[idx]
}

func (g *LinearGradient) Type() string {
	return "linear-gradient"
}

func (g *LinearGradient) SVGDefs(id string) string {
	// Convert angle to x1,y1,x2,y2 coordinates
	rad := g.Angle * math.Pi / 180
	x1 := 50 - 50*math.Cos(rad)
	y1 := 50 - 50*math.Sin(rad)
	x2 := 50 + 50*math.Cos(rad)
	y2 := 50 + 50*math.Sin(rad)

	var sb strings.Builder
	fmt.Fprintf(&sb, `<linearGradient id="%s" x1="%.0f%%" y1="%.0f%%" x2="%.0f%%" y2="%.0f%%">`,
		id, x1, y1, x2, y2)

	for i, stop := range g.Stops {
		offset := float64(i) / float64(len(g.Stops)-1) * 100
		fmt.Fprintf(&sb, `<stop offset="%.0f%%" stop-color="%s"/>`, offset, stop)
	}

	sb.WriteString("</linearGradient>")
	return sb.String()
}

func (g *LinearGradient) SVGFill(id string) string {
	return fmt.Sprintf("url(#%s)", id)
}

// RadialGradient represents a radial color gradient.
type RadialGradient struct {
	CenterX float64  // Center X (0.0-1.0)
	CenterY float64  // Center Y (0.0-1.0)
	Stops   []string // Color stops (hex values)
}

// NewRadialGradient creates a radial gradient from center point.
func NewRadialGradient(cx, cy float64, stops []string) *RadialGradient {
	return &RadialGradient{
		CenterX: cx,
		CenterY: cy,
		Stops:   stops,
	}
}

func (g *RadialGradient) ColorAt(x, y float64) string {
	if len(g.Stops) == 0 {
		return "#000000"
	}
	dx := x - g.CenterX
	dy := y - g.CenterY
	dist := math.Sqrt(dx*dx+dy*dy) / 0.707 // Normalize to roughly 0-1
	dist = math.Max(0, math.Min(1, dist))
	idx := int(dist * float64(len(g.Stops)-1))
	return g.Stops[idx]
}

func (g *RadialGradient) Type() string {
	return "radial-gradient"
}

func (g *RadialGradient) SVGDefs(id string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, `<radialGradient id="%s" cx="%.0f%%" cy="%.0f%%" r="70%%">`,
		id, g.CenterX*100, g.CenterY*100)

	for i, stop := range g.Stops {
		offset := float64(i) / float64(len(g.Stops)-1) * 100
		fmt.Fprintf(&sb, `<stop offset="%.0f%%" stop-color="%s"/>`, offset, stop)
	}

	sb.WriteString("</radialGradient>")
	return sb.String()
}

func (g *RadialGradient) SVGFill(id string) string {
	return fmt.Sprintf("url(#%s)", id)
}

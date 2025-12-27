package elements

import "github.com/ahmedtahas/qr-gode/internal/colors"

// FinderRenderer renders finder patterns.
type FinderRenderer struct {
	// Simple mode (all layers same style)
	Shape string
	Color colors.Color

	// Detailed mode (per-layer styling)
	OuterShape        string
	OuterColor        colors.Color
	OuterCornerRadius float64

	MiddleShape        string
	MiddleColor        colors.Color
	MiddleCornerRadius float64

	CenterShape string
	CenterColor colors.Color

	CellSize float64
	Detailed bool // true = use per-layer styling
}

// Render renders a finder pattern at the given position.
// position: "top-left", "top-right", "bottom-left"
func (r *FinderRenderer) Render(position string) string {
	// TODO:
	// Finder pattern is 7x7 modules:
	// - Outer ring: 7x7 dark border
	// - Middle ring: 5x5 light area (inside outer)
	// - Center: 3x3 dark square
	//
	// For simple mode: render all with same shape/color
	// For detailed mode: render each layer separately
	//
	// Position determines top-left corner:
	// - top-left: (0, 0)
	// - top-right: (size-7, 0)
	// - bottom-left: (0, size-7)
	return ""
}

// RenderAll renders all three finder patterns.
func (r *FinderRenderer) RenderAll(matrixSize int) string {
	// TODO:
	// 1. Render top-left finder
	// 2. Render top-right finder
	// 3. Render bottom-left finder
	// 4. Return <g> containing all
	return ""
}

// renderOuterRing renders the 7x7 outer border.
func (r *FinderRenderer) renderOuterRing(x, y float64) string {
	// TODO
	return ""
}

// renderMiddleRing renders the 5x5 light middle area.
func (r *FinderRenderer) renderMiddleRing(x, y float64) string {
	// TODO
	return ""
}

// renderCenter renders the 3x3 dark center.
func (r *FinderRenderer) renderCenter(x, y float64) string {
	// TODO
	return ""
}

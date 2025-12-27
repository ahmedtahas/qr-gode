package elements

import "github.com/ahmedtahas/qr-gode/internal/colors"

// AlignmentRenderer renders alignment patterns.
type AlignmentRenderer struct {
	// Simple mode
	Shape string
	Color colors.Color

	// Detailed mode
	OuterShape  string
	OuterColor  colors.Color
	CenterShape string
	CenterColor colors.Color

	CellSize float64
	Detailed bool
}

// Render renders an alignment pattern centered at (cx, cy) in module coordinates.
func (r *AlignmentRenderer) Render(cx, cy int) string {
	// TODO:
	// Alignment pattern is 5x5 modules:
	// - Outer ring: 5x5 dark border
	// - Middle: 3x3 light area
	// - Center: 1x1 dark module
	return ""
}

// RenderAll renders all alignment patterns for the given version.
func (r *AlignmentRenderer) RenderAll(version int, matrixSize int) string {
	// TODO:
	// 1. Get alignment pattern positions for version
	// 2. Skip positions that overlap finder patterns
	// 3. Render each alignment pattern
	// 4. Return <g> containing all
	return ""
}

// getAlignmentPositions returns center positions for alignment patterns.
func getAlignmentPositions(version int) []int {
	// TODO: Return positions from lookup table
	// Version 1: no alignment patterns
	// Version 2: [6, 18]
	// etc.
	return nil
}

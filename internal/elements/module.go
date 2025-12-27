package elements

import "github.com/ahmedtahas/qr-gode/internal/colors"

// ModuleRenderer renders data modules.
type ModuleRenderer struct {
	Shape    string       // Shape name or SVG path
	Color    colors.Color // Color source
	Size     float64      // Size as fraction of cell
	CellSize float64      // Actual cell size in SVG units
}

// NewModuleRenderer creates a module renderer.
func NewModuleRenderer(shape string, color colors.Color, size, cellSize float64) *ModuleRenderer {
	return &ModuleRenderer{
		Shape:    shape,
		Color:    color,
		Size:     size,
		CellSize: cellSize,
	}
}

// Render renders a single module at grid position (col, row).
// Returns SVG element string.
func (r *ModuleRenderer) Render(col, row int) string {
	// TODO:
	// 1. Calculate center position: (col + 0.5) * CellSize, (row + 0.5) * CellSize
	// 2. Get shape SVG path
	// 3. Get color at normalized position
	// 4. Generate SVG element with transform (translate, scale)
	return ""
}

// RenderAll renders all data modules from the matrix.
// Returns combined SVG group.
func (r *ModuleRenderer) RenderAll(modules [][]bool, reserved [][]bool) string {
	// TODO:
	// 1. Iterate through all modules
	// 2. Skip reserved (function pattern) modules
	// 3. Render only dark modules
	// 4. Return <g> containing all module elements
	return ""
}

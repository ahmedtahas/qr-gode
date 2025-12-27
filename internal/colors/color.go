package colors

// Color represents a color source for rendering modules.
// It can be solid, gradient, or image-sampled.
type Color interface {
	// ColorAt returns the hex color for a module at position (x, y).
	// Coordinates are normalized 0.0-1.0 across the QR code.
	ColorAt(x, y float64) string

	// Type returns the color type identifier.
	Type() string

	// SVGDefs returns any SVG <defs> needed (e.g., gradient definitions).
	// Returns empty string if no defs needed.
	SVGDefs(id string) string

	// SVGFill returns the fill attribute value.
	// Could be a color like "#ff0000" or a reference like "url(#gradient-1)".
	SVGFill(id string) string
}

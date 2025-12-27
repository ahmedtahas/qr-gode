package shapes

import "fmt"

func init() {
	// Register all built-in shapes
	Register(&square{})
	Register(&circle{})
	Register(&roundedSquare{radius: 0.3})
	Register(&diamond{})
	Register(&dot{scale: 0.7})
	Register(&star{})
	Register(&heart{})
}

// Square - standard QR module
type square struct{}

func (s *square) Name() string { return "square" }
func (s *square) SVGPath() string {
	return "M0 0h1v1h-1z"
}

// Circle - circular module
type circle struct{}

func (s *circle) Name() string { return "circle" }
func (s *circle) SVGPath() string {
	// Circle inscribed in unit square, centered at (0.5, 0.5) with radius 0.5
	// Using arc commands: M start A rx ry rotation large-arc sweep end
	return "M0.5 0A0.5 0.5 0 0 1 0.5 1A0.5 0.5 0 0 1 0.5 0z"
}

// RoundedSquare - square with rounded corners
type roundedSquare struct {
	radius float64 // corner radius as fraction of size (0-0.5)
}

func (s *roundedSquare) Name() string { return "rounded" }
func (s *roundedSquare) SVGPath() string {
	r := s.radius
	// Rounded rectangle path
	return fmt.Sprintf("M%.2f 0h%.2fq%.2f 0 %.2f %.2fv%.2fq0 %.2f -%.2f %.2fh-%.2fq-%.2f 0 -%.2f -%.2fv-%.2fq0 -%.2f %.2f -%.2fz",
		r, 1-2*r, r, r, r, 1-2*r, r, r, r, 1-2*r, r, r, r, 1-2*r, r, r, r)
}

// Diamond - 45-degree rotated square
type diamond struct{}

func (s *diamond) Name() string { return "diamond" }
func (s *diamond) SVGPath() string {
	// Diamond centered in unit square
	return "M0.5 0L1 0.5L0.5 1L0 0.5z"
}

// Dot - small centered circle
type dot struct {
	scale float64 // size as fraction of cell (0-1)
}

func (s *dot) Name() string { return "dot" }
func (s *dot) SVGPath() string {
	r := s.scale / 2
	offset := (1 - s.scale) / 2
	// Circle centered at (0.5, 0.5) with scaled radius
	return fmt.Sprintf("M%.2f %.2fA%.2f %.2f 0 0 1 %.2f %.2fA%.2f %.2f 0 0 1 %.2f %.2fz",
		0.5, offset, r, r, 0.5, offset+s.scale, r, r, 0.5, offset)
}

// Star - 4-pointed star
type star struct{}

func (s *star) Name() string { return "star" }
func (s *star) SVGPath() string {
	// 4-pointed star
	return "M0.5 0L0.6 0.4L1 0.5L0.6 0.6L0.5 1L0.4 0.6L0 0.5L0.4 0.4z"
}

// Heart shape
type heart struct{}

func (s *heart) Name() string { return "heart" }
func (s *heart) SVGPath() string {
	// Heart shape scaled to unit square
	return "M0.5 0.2C0.5 0.1 0.4 0 0.25 0C0.1 0 0 0.15 0 0.3C0 0.55 0.5 1 0.5 1C0.5 1 1 0.55 1 0.3C1 0.15 0.9 0 0.75 0C0.6 0 0.5 0.1 0.5 0.2z"
}

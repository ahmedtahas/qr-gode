package colors

import "fmt"

// Solid represents a single hex color.
type Solid struct {
	Hex string
}

// NewSolid creates a solid color from hex string.
func NewSolid(hex string) *Solid {
	// TODO: Validate hex format
	return &Solid{Hex: hex}
}

func (s *Solid) ColorAt(x, y float64) string {
	return s.Hex
}

func (s *Solid) Type() string {
	return "solid"
}

func (s *Solid) SVGDefs(id string) string {
	return ""
}

func (s *Solid) SVGFill(id string) string {
	return s.Hex
}

// ParseHex validates and normalizes a hex color string.
func ParseHex(hex string) (string, error) {
	// TODO:
	// Accept: #RGB, #RRGGBB, #RRGGBBAA
	// Normalize to #RRGGBB or #RRGGBBAA
	if len(hex) == 0 || hex[0] != '#' {
		return "", fmt.Errorf("invalid hex color: %s", hex)
	}
	return hex, nil
}

package colors

import "fmt"

// Solid represents a single hex color.
type Solid struct {
	Hex string
}

// NewSolid creates a solid color from hex string.
func NewSolid(hex string) *Solid {
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
// Accepts: #RGB, #RRGGBB, #RRGGBBAA
func ParseHex(hex string) (string, error) {
	if len(hex) == 0 || hex[0] != '#' {
		return "", fmt.Errorf("invalid hex color: %s", hex)
	}
	switch len(hex) {
	case 4, 7, 9: // #RGB, #RRGGBB, #RRGGBBAA
		return hex, nil
	default:
		return "", fmt.Errorf("invalid hex color length: %s", hex)
	}
}

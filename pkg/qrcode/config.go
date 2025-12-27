package qrcode

import "github.com/ahmedtahas/qr-gode/internal/colors"

// ErrorCorrectionLevel defines the amount of redundancy in the QR code.
// Higher levels provide more error recovery but result in denser codes.
type ErrorCorrectionLevel int

const (
	LevelL ErrorCorrectionLevel = iota // ~7% recovery capacity
	LevelM                             // ~15% recovery capacity (default)
	LevelQ                             // ~25% recovery capacity
	LevelH                             // ~30% recovery capacity
)

// Color is an alias to the internal color interface for advanced usage.
// Most users should use the builder methods like Foreground(), LinearGradient(), etc.
type Color = colors.Color

// NewSolidColor creates a solid color from a hex string.
// This is exported for advanced configuration via Config struct.
func NewSolidColor(hex string) Color {
	return colors.NewSolid(hex)
}

// NewLinearGradientColor creates a linear gradient color.
// Angle is in degrees (0=right, 90=down, 180=left, 270=up).
func NewLinearGradientColor(angle float64, stops []string) Color {
	return colors.NewLinearGradient(angle, stops)
}

// NewRadialGradientColor creates a radial gradient color.
// Center coordinates are fractions from 0.0 to 1.0.
func NewRadialGradientColor(centerX, centerY float64, stops []string) Color {
	return colors.NewRadialGradient(centerX, centerY, stops)
}

// Config holds all configuration for QR code generation.
type Config struct {
	// QR data settings
	ErrorCorrection ErrorCorrectionLevel

	// Overall dimensions
	Size      int  // Output size in pixels
	QuietZone int  // Margin around QR (in modules)

	// Styling
	Background colors.Color
	Modules    ModuleStyle
	Finders    FinderStyle
	Alignment  AlignmentStyle
	Timing     TimingStyle
	Logo       *LogoConfig

	// Custom images for elements
	Images *CustomImages
}

// ModuleStyle defines how data modules are rendered.
type ModuleStyle struct {
	Shape string       // Shape name or SVG path
	Color colors.Color // Solid, gradient, or image-sampled
	Size  float64      // Size as fraction of cell (0.0-1.0)
}

// FinderStyle defines how finder patterns are rendered.
type FinderStyle struct {
	// Simple mode: style all three layers uniformly
	Shape string
	Color colors.Color

	// Detailed mode: style each layer separately
	Outer  *FinderLayerStyle
	Middle *FinderLayerStyle
	Center *FinderLayerStyle
}

// FinderLayerStyle defines one layer of a finder pattern.
type FinderLayerStyle struct {
	Shape        string
	Color        colors.Color
	CornerRadius float64 // For rounded shapes
}

// AlignmentStyle defines how alignment patterns are rendered.
type AlignmentStyle struct {
	// Simple mode
	Shape string
	Color colors.Color

	// Detailed mode
	Outer  *AlignmentLayerStyle
	Center *AlignmentLayerStyle
}

// AlignmentLayerStyle defines one layer of an alignment pattern.
type AlignmentLayerStyle struct {
	Shape string
	Color colors.Color
}

// TimingStyle defines how timing patterns are rendered.
type TimingStyle struct {
	Shape string
	Color colors.Color
}

// LogoConfig defines the logo overlay settings.
type LogoConfig struct {
	Path       string // Path to logo image (required)
	Width      int    // Optional: logo width in pixels (0 = auto-calculate)
	Height     int    // Optional: logo height in pixels (0 = auto-calculate)
	Background string // Background color behind logo (hex or "transparent", default white)
}

// CustomImages defines custom PNG images for different QR elements.
type CustomImages struct {
	Finder    string // Path to PNG for finder pattern modules (7x7 outer squares)
	Module    string // Path to PNG for regular data modules
	Alignment string // Path to PNG for alignment pattern modules (5x5 squares)
}

// LoadConfig loads a Config from a TOML file.
func LoadConfig(path string) (*Config, error) {
	// TODO: Parse TOML file
	return nil, nil
}

// DefaultConfig returns a config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		ErrorCorrection: LevelM,
		Size:            256,
		QuietZone:       4,
		Background:      NewSolidColor("#FFFFFF"),
		Modules: ModuleStyle{
			Shape: "square",
			Color: NewSolidColor("#000000"),
			Size:  1.0,
		},
		Finders: FinderStyle{
			Shape: "square",
			Color: NewSolidColor("#000000"),
		},
		Alignment: AlignmentStyle{
			Shape: "square",
			Color: NewSolidColor("#000000"),
		},
		Timing: TimingStyle{
			Shape: "square",
			Color: NewSolidColor("#000000"),
		},
	}
}

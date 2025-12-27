package qrgode

import (
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/ahmedtahas/qr-gode/internal/colors"
	"github.com/ahmedtahas/qr-gode/internal/encoder"
)

// QRCode represents a QR code generator with fluent configuration.
type QRCode struct {
	data   string
	config *Config
	errs   []error // Accumulated validation errors
}

// New creates a new QR code generator for the given data.
//
// Example:
//
//	qr := qrcode.New("https://example.com")
//	svg, err := qr.SVG()
func New(data string) *QRCode {
	return &QRCode{
		data:   data,
		config: DefaultConfig(),
	}
}

// Size sets the output size in pixels. Default is 256.
func (q *QRCode) Size(pixels int) *QRCode {
	q.config.Size = pixels
	return q
}

// QuietZone sets the margin around the QR code in modules. Default is 4.
func (q *QRCode) QuietZone(modules int) *QRCode {
	q.config.QuietZone = modules
	return q
}

// ErrorCorrection sets the error correction level. Default is LevelM.
//
// Available levels:
//   - LevelL: ~7% recovery
//   - LevelM: ~15% recovery (default)
//   - LevelQ: ~25% recovery
//   - LevelH: ~30% recovery
func (q *QRCode) ErrorCorrection(level ErrorCorrectionLevel) *QRCode {
	q.config.ErrorCorrection = level
	return q
}

// Shape sets the module shape.
//
// Available shapes: ShapeSquare, ShapeCircle, ShapeRounded, ShapeDiamond, ShapeDot, ShapeStar, ShapeHeart
func (q *QRCode) Shape(shape Shape) *QRCode {
	q.config.Modules.Shape = string(shape)
	return q
}

// Foreground sets the foreground color (modules) as a hex string.
//
// Example: "#000000", "#3498db", "rgb(52, 152, 219)"
func (q *QRCode) Foreground(hex string) *QRCode {
	q.config.Modules.Color = colors.NewSolid(hex)
	return q
}

// Background sets the background color as a hex string.
//
// Example: "#ffffff", "#f0f0f0", "transparent"
func (q *QRCode) Background(hex string) *QRCode {
	q.config.Background = colors.NewSolid(hex)
	return q
}

// LinearGradient sets a linear gradient for the modules.
//
// Parameters:
//   - angle: Direction in degrees (0=right, 90=down, 180=left, 270=up)
//   - colorStops: At least 2 hex color strings
//
// Example:
//
//	qr.LinearGradient(45, "#ff0000", "#0000ff")
func (q *QRCode) LinearGradient(angle float64, colorStops ...string) *QRCode {
	q.config.Modules.Color = colors.NewLinearGradient(angle, colorStops)
	return q
}

// RadialGradient sets a radial gradient for the modules.
//
// Parameters:
//   - centerX, centerY: Center position as fractions (0.0-1.0)
//   - colorStops: At least 2 hex color strings
//
// Example:
//
//	qr.RadialGradient(0.5, 0.5, "#ff0000", "#0000ff")
func (q *QRCode) RadialGradient(centerX, centerY float64, colorStops ...string) *QRCode {
	q.config.Modules.Color = colors.NewRadialGradient(centerX, centerY, colorStops)
	return q
}

// ModuleImage sets a custom PNG/JPG/SVG image for data modules.
// Each module will be rendered using this image.
// The image is validated immediately; errors are collected and returned by SVG()/SaveAs().
func (q *QRCode) ModuleImage(path string) *QRCode {
	q.ensureImages()
	if err := ValidateModuleImage(path); err != nil {
		q.errs = append(q.errs, err)
	} else {
		q.config.Images.Module = path
	}
	return q
}

// FinderImage sets a custom image for finder patterns.
// The image represents the full 7x7 finder pattern (the big squares in corners).
// It will be mirrored appropriately for each corner.
// The image is validated immediately; errors are collected and returned by SVG()/SaveAs().
func (q *QRCode) FinderImage(path string) *QRCode {
	q.ensureImages()
	if err := ValidateFinderImage(path); err != nil {
		q.errs = append(q.errs, err)
	} else {
		q.config.Images.Finder = path
	}
	return q
}

// AlignmentImage sets a custom image for alignment patterns.
// The image represents the full 5x5 alignment pattern (smaller squares in larger QR codes).
// The image is validated immediately; errors are collected and returned by SVG()/SaveAs().
func (q *QRCode) AlignmentImage(path string) *QRCode {
	q.ensureImages()
	if err := ValidateAlignmentImage(path); err != nil {
		q.errs = append(q.errs, err)
	} else {
		q.config.Images.Alignment = path
	}
	return q
}

// Logo sets a logo image from a file path to display in the center of the QR code.
// The logo will have a white background padding by default.
func (q *QRCode) Logo(path string) *QRCode {
	q.ensureLogo()
	if err := ValidateLogoImage(path); err != nil {
		q.errs = append(q.errs, err)
	} else {
		q.config.Logo.Path = path
	}
	return q
}

// LogoImage sets an in-memory image as the logo to display in the center of the QR code.
// This takes precedence over Logo() if both are set.
// The logo will have a white background padding by default.
func (q *QRCode) LogoImage(img image.Image) *QRCode {
	q.ensureLogo()
	if img == nil {
		return q
	}
	q.config.Logo.Image = img
	return q
}

// LogoWidth sets the logo width in pixels. Height will be calculated to preserve aspect ratio.
// If not set, logo size is auto-calculated to fit within 15-30% of QR size.
func (q *QRCode) LogoWidth(pixels int) *QRCode {
	q.ensureLogo()
	q.config.Logo.Width = pixels
	return q
}

// LogoHeight sets the logo height in pixels. Width will be calculated to preserve aspect ratio.
// If not set, logo size is auto-calculated to fit within 15-30% of QR size.
func (q *QRCode) LogoHeight(pixels int) *QRCode {
	q.ensureLogo()
	q.config.Logo.Height = pixels
	return q
}

// LogoDimensions sets both logo width and height in pixels.
func (q *QRCode) LogoDimensions(width, height int) *QRCode {
	q.ensureLogo()
	q.config.Logo.Width = width
	q.config.Logo.Height = height
	return q
}

// LogoBackground sets the background color behind the logo.
// Use hex color (e.g., "#ffffff") or "transparent" for no background.
func (q *QRCode) LogoBackground(color string) *QRCode {
	q.ensureLogo()
	q.config.Logo.Background = color
	return q
}

func (q *QRCode) ensureImages() {
	if q.config.Images == nil {
		q.config.Images = &CustomImages{}
	}
}

func (q *QRCode) ensureLogo() {
	if q.config.Logo == nil {
		q.config.Logo = &LogoConfig{}
	}
}

// Validate checks the configuration and returns any validation errors.
// This is called automatically by SVG() and SaveAs(), but can be called
// explicitly to check configuration before generation.
func (q *QRCode) Validate() []error {
	// Return accumulated errors from builder methods
	if len(q.errs) > 0 {
		return q.errs
	}
	// Run full config validation
	return ValidateConfig(q.config)
}

// SVG generates and returns the QR code as SVG bytes.
// Returns an error if validation fails or encoding fails.
func (q *QRCode) SVG() ([]byte, error) {
	// Check for validation errors
	if errs := q.Validate(); len(errs) > 0 {
		return nil, errs[0] // Return first error
	}

	ecl := encoder.ErrorCorrectionLevel(q.config.ErrorCorrection)
	enc := encoder.New(q.data, ecl)
	matrix, err := enc.Encode()
	if err != nil {
		return nil, err
	}

	renderer := NewRenderer(matrix, q.config)
	return renderer.RenderSVG()
}

// SVGString generates and returns the QR code as an SVG string.
func (q *QRCode) SVGString() (string, error) {
	svg, err := q.SVG()
	if err != nil {
		return "", err
	}
	return string(svg), nil
}

// SaveAs generates the QR code and saves it to the specified file.
// Currently supports .svg files. PNG support is planned.
func (q *QRCode) SaveAs(path string) error {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".png" {
		return &UnsupportedFormatError{Format: "png"}
	}

	svg, err := q.SVG()
	if err != nil {
		return err
	}
	return os.WriteFile(path, svg, 0644)
}

// GetConfig returns the underlying configuration for advanced customization.
func (q *QRCode) GetConfig() *Config {
	return q.config
}

// UnsupportedFormatError is returned when an unsupported output format is requested.
type UnsupportedFormatError struct {
	Format string
}

func (e *UnsupportedFormatError) Error() string {
	return "unsupported format: " + e.Format + " (only SVG is currently supported)"
}

package qrgode

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ahmedtahas/qr-gode/internal/encoder"
)

// Generate creates a QR code from the given data and config.
// Returns SVG as a byte slice.
// If cfg is nil, DefaultConfig() is used.
func Generate(data string, cfg *Config) ([]byte, error) {
	if data == "" {
		return nil, &ValidationError{Field: "Data", Message: "cannot be empty"}
	}

	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Validate configuration
	if errs := ValidateConfig(cfg); len(errs) > 0 {
		return nil, errs[0]
	}

	// Convert public ECL to internal ECL
	ecl := encoder.ErrorCorrectionLevel(cfg.ErrorCorrection)

	// Encode data using internal/encoder
	enc := encoder.New(data, ecl)
	matrix, err := enc.Encode()
	if err != nil {
		return nil, err
	}

	// Render to SVG
	renderer := newRenderer(matrix, cfg)
	return renderer.renderSVG()
}

// GenerateToFile creates a QR code and writes it to the specified path.
// The format is chosen from the file extension: .svg or .png.
func GenerateToFile(data string, cfg *Config, path string) error {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".svg", "":
		svg, err := Generate(data, cfg)
		if err != nil {
			return err
		}
		return os.WriteFile(path, svg, 0644)
	case ".png":
		png, err := GeneratePNG(data, cfg)
		if err != nil {
			return err
		}
		return os.WriteFile(path, png, 0644)
	default:
		return &UnsupportedFormatError{Format: strings.TrimPrefix(ext, ".")}
	}
}

// GeneratePNG creates a QR code and returns it as PNG bytes.
// If cfg is nil, DefaultConfig() is used.
func GeneratePNG(data string, cfg *Config) ([]byte, error) {
	if data == "" {
		return nil, &ValidationError{Field: "Data", Message: "cannot be empty"}
	}
	if cfg == nil {
		cfg = DefaultConfig()
	}
	if errs := ValidateConfig(cfg); len(errs) > 0 {
		return nil, errs[0]
	}
	ecl := encoder.ErrorCorrectionLevel(cfg.ErrorCorrection)
	enc := encoder.New(data, ecl)
	matrix, err := enc.Encode()
	if err != nil {
		return nil, err
	}
	r := newRenderer(matrix, cfg)
	return r.renderPNG()
}

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
	renderer := NewRenderer(matrix, cfg)
	return renderer.RenderSVG()
}

// GenerateToFile creates a QR code and writes it to the specified path.
// Supports .svg and .png extensions.
func GenerateToFile(data string, cfg *Config, path string) error {
	svg, err := Generate(data, cfg)
	if err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".png" {
		// PNG not yet implemented
		renderer := NewRenderer(nil, cfg)
		_, err := renderer.RenderPNG()
		return err
	}

	return os.WriteFile(path, svg, 0644)
}

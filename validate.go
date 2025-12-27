package qrgode

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"os"
	"path/filepath"
	"strings"
)

// ValidationError represents an error during configuration validation.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateImage checks if the given path points to a valid image file.
// Returns the image dimensions if valid.
func ValidateImage(path string) (width, height int, err error) {
	if path == "" {
		return 0, 0, errors.New("empty path")
	}

	// Check file exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, 0, fmt.Errorf("file not found: %s", path)
		}
		return 0, 0, fmt.Errorf("cannot access file: %w", err)
	}

	if info.IsDir() {
		return 0, 0, fmt.Errorf("path is a directory, not a file: %s", path)
	}

	// Check extension
	ext := strings.ToLower(filepath.Ext(path))
	validExts := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".svg":  true,
	}
	if !validExts[ext] {
		return 0, 0, fmt.Errorf("unsupported image format: %s (supported: png, jpg, jpeg, svg)", ext)
	}

	// For SVG, we can't easily get dimensions, just check it exists
	if ext == ".svg" {
		return 0, 0, nil
	}

	// Open and decode image to verify it's valid
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	// Decode image config (doesn't load full image, just header)
	cfg, format, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid image file: %w", err)
	}

	_ = format // We already know format from extension
	return cfg.Width, cfg.Height, nil
}

// ValidateFinderImage validates an image intended for finder patterns.
// Finder patterns are 7x7 modules, so ideally the image should be square.
func ValidateFinderImage(path string) error {
	w, h, err := ValidateImage(path)
	if err != nil {
		return &ValidationError{
			Field:   "FinderImage",
			Message: err.Error(),
		}
	}

	// SVG doesn't return dimensions
	if w == 0 && h == 0 {
		return nil
	}

	// Warn if not square (but don't fail)
	if w != h {
		// We could log a warning here, but for now just accept it
		// The image will be stretched to fit
	}

	return nil
}

// ValidateAlignmentImage validates an image intended for alignment patterns.
// Alignment patterns are 5x5 modules, so ideally the image should be square.
func ValidateAlignmentImage(path string) error {
	w, h, err := ValidateImage(path)
	if err != nil {
		return &ValidationError{
			Field:   "AlignmentImage",
			Message: err.Error(),
		}
	}

	// SVG doesn't return dimensions
	if w == 0 && h == 0 {
		return nil
	}

	// Warn if not square (but don't fail)
	if w != h {
		// We could log a warning here, but for now just accept it
	}

	return nil
}

// ValidateModuleImage validates an image intended for data modules.
func ValidateModuleImage(path string) error {
	_, _, err := ValidateImage(path)
	if err != nil {
		return &ValidationError{
			Field:   "ModuleImage",
			Message: err.Error(),
		}
	}
	return nil
}

// ValidateLogoImage validates an image intended for the center logo.
func ValidateLogoImage(path string) error {
	_, _, err := ValidateImage(path)
	if err != nil {
		return &ValidationError{
			Field:   "Logo",
			Message: err.Error(),
		}
	}
	return nil
}

// ValidateConfig validates the entire configuration.
// Returns a list of all validation errors found.
func ValidateConfig(cfg *Config) []error {
	var errs []error

	// Validate size
	if cfg.Size <= 0 {
		errs = append(errs, &ValidationError{
			Field:   "Size",
			Message: "must be positive",
		})
	}

	// Validate quiet zone
	if cfg.QuietZone < 0 {
		errs = append(errs, &ValidationError{
			Field:   "QuietZone",
			Message: "cannot be negative",
		})
	}

	// Validate custom images if provided
	if cfg.Images != nil {
		if cfg.Images.Module != "" {
			if err := ValidateModuleImage(cfg.Images.Module); err != nil {
				errs = append(errs, err)
			}
		}
		if cfg.Images.Finder != "" {
			if err := ValidateFinderImage(cfg.Images.Finder); err != nil {
				errs = append(errs, err)
			}
		}
		if cfg.Images.Alignment != "" {
			if err := ValidateAlignmentImage(cfg.Images.Alignment); err != nil {
				errs = append(errs, err)
			}
		}
	}

	// Validate logo if provided
	if cfg.Logo != nil && cfg.Logo.Path != "" {
		if err := ValidateLogoImage(cfg.Logo.Path); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

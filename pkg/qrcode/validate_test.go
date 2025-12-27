package qrcode

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateImage_NotFound(t *testing.T) {
	_, _, err := ValidateImage("/nonexistent/path/image.png")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got: %v", err)
	}
}

func TestValidateImage_EmptyPath(t *testing.T) {
	_, _, err := ValidateImage("")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestValidateImage_Directory(t *testing.T) {
	_, _, err := ValidateImage(os.TempDir())
	if err == nil {
		t.Error("expected error for directory")
	}
	if !strings.Contains(err.Error(), "directory") {
		t.Errorf("expected 'directory' error, got: %v", err)
	}
}

func TestValidateImage_UnsupportedFormat(t *testing.T) {
	// Create a temp file with wrong extension
	tmpFile := filepath.Join(os.TempDir(), "test.txt")
	os.WriteFile(tmpFile, []byte("test"), 0644)
	defer os.Remove(tmpFile)

	_, _, err := ValidateImage(tmpFile)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("expected 'unsupported' error, got: %v", err)
	}
}

func TestValidateImage_InvalidPNG(t *testing.T) {
	// Create a temp file with PNG extension but invalid content
	tmpFile := filepath.Join(os.TempDir(), "invalid.png")
	os.WriteFile(tmpFile, []byte("not a png"), 0644)
	defer os.Remove(tmpFile)

	_, _, err := ValidateImage(tmpFile)
	if err == nil {
		t.Error("expected error for invalid PNG")
	}
	if !strings.Contains(err.Error(), "invalid image") {
		t.Errorf("expected 'invalid image' error, got: %v", err)
	}
}

func TestValidateImage_ValidPNG(t *testing.T) {
	// Use the test images in the project root if they exist
	testImages := []string{
		"../../test_module.png",
		"../../test_finder.png",
		"../../test_align.png",
	}

	for _, img := range testImages {
		// Get absolute path
		absPath, err := filepath.Abs(img)
		if err != nil {
			continue
		}

		// Skip if file doesn't exist
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			t.Logf("skipping %s (not found)", img)
			continue
		}

		w, h, err := ValidateImage(absPath)
		if err != nil {
			t.Errorf("unexpected error for %s: %v", img, err)
			continue
		}
		if w <= 0 || h <= 0 {
			t.Errorf("expected positive dimensions for %s, got %dx%d", img, w, h)
		}
	}
}

func TestValidateFinderImage_NotFound(t *testing.T) {
	err := ValidateFinderImage("/nonexistent/finder.png")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
	if valErr.Field != "FinderImage" {
		t.Errorf("expected field 'FinderImage', got '%s'", valErr.Field)
	}
}

func TestValidateAlignmentImage_NotFound(t *testing.T) {
	err := ValidateAlignmentImage("/nonexistent/align.png")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
	if valErr.Field != "AlignmentImage" {
		t.Errorf("expected field 'AlignmentImage', got '%s'", valErr.Field)
	}
}

func TestValidateModuleImage_NotFound(t *testing.T) {
	err := ValidateModuleImage("/nonexistent/module.png")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
	if valErr.Field != "ModuleImage" {
		t.Errorf("expected field 'ModuleImage', got '%s'", valErr.Field)
	}
}

func TestValidateConfig_NegativeSize(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Size = -1

	errs := ValidateConfig(cfg)
	if len(errs) == 0 {
		t.Error("expected error for negative size")
	}
}

func TestValidateConfig_NegativeQuietZone(t *testing.T) {
	cfg := DefaultConfig()
	cfg.QuietZone = -1

	errs := ValidateConfig(cfg)
	if len(errs) == 0 {
		t.Error("expected error for negative quiet zone")
	}
}

func TestValidateConfig_InvalidImages(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Images = &CustomImages{
		Module:    "/nonexistent/module.png",
		Finder:    "/nonexistent/finder.png",
		Alignment: "/nonexistent/align.png",
	}

	errs := ValidateConfig(cfg)
	if len(errs) != 3 {
		t.Errorf("expected 3 errors, got %d", len(errs))
	}
}

func TestBuilder_InvalidModuleImage(t *testing.T) {
	qr := New("test").ModuleImage("/nonexistent/module.png")

	_, err := qr.SVG()
	if err == nil {
		t.Error("expected error for invalid module image")
	}
}

func TestBuilder_InvalidFinderImage(t *testing.T) {
	qr := New("test").FinderImage("/nonexistent/finder.png")

	_, err := qr.SVG()
	if err == nil {
		t.Error("expected error for invalid finder image")
	}
}

func TestBuilder_InvalidAlignmentImage(t *testing.T) {
	qr := New("test").AlignmentImage("/nonexistent/align.png")

	_, err := qr.SVG()
	if err == nil {
		t.Error("expected error for invalid alignment image")
	}
}

func TestBuilder_Validate(t *testing.T) {
	qr := New("test").
		ModuleImage("/nonexistent/module.png").
		FinderImage("/nonexistent/finder.png")

	errs := qr.Validate()
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Field:   "TestField",
		Message: "test message",
	}

	expected := "TestField: test message"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

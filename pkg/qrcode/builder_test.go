package qrcode

import (
	"os"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	qr := New("test")
	if qr.data != "test" {
		t.Errorf("expected data 'test', got '%s'", qr.data)
	}
	if qr.config == nil {
		t.Error("expected config to be initialized")
	}
}

func TestBuilderChaining(t *testing.T) {
	qr := New("https://example.com").
		Size(512).
		QuietZone(2).
		ErrorCorrection(LevelH).
		Shape("circle").
		Foreground("#ff0000").
		Background("#ffffff")

	if qr.config.Size != 512 {
		t.Errorf("expected size 512, got %d", qr.config.Size)
	}
	if qr.config.QuietZone != 2 {
		t.Errorf("expected quiet zone 2, got %d", qr.config.QuietZone)
	}
	if qr.config.ErrorCorrection != LevelH {
		t.Errorf("expected LevelH, got %d", qr.config.ErrorCorrection)
	}
	if qr.config.Modules.Shape != "circle" {
		t.Errorf("expected shape 'circle', got '%s'", qr.config.Modules.Shape)
	}
}

func TestLinearGradient(t *testing.T) {
	qr := New("test").LinearGradient(45, "#ff0000", "#0000ff")

	if qr.config.Modules.Color == nil {
		t.Error("expected color to be set")
	}
	if qr.config.Modules.Color.Type() != "linear-gradient" {
		t.Errorf("expected linear-gradient, got %s", qr.config.Modules.Color.Type())
	}
}

func TestRadialGradient(t *testing.T) {
	qr := New("test").RadialGradient(0.5, 0.5, "#ff0000", "#0000ff")

	if qr.config.Modules.Color == nil {
		t.Error("expected color to be set")
	}
	if qr.config.Modules.Color.Type() != "radial-gradient" {
		t.Errorf("expected radial-gradient, got %s", qr.config.Modules.Color.Type())
	}
}

func TestCustomImages(t *testing.T) {
	// Test with invalid paths - should accumulate errors
	qr := New("test").
		ModuleImage("nonexistent_module.png").
		FinderImage("nonexistent_finder.png").
		AlignmentImage("nonexistent_align.png")

	// Images struct should be initialized
	if qr.config.Images == nil {
		t.Error("expected images struct to be initialized")
	}

	// Paths should NOT be set because validation failed
	if qr.config.Images.Module != "" {
		t.Errorf("expected empty module path due to validation, got %s", qr.config.Images.Module)
	}
	if qr.config.Images.Finder != "" {
		t.Errorf("expected empty finder path due to validation, got %s", qr.config.Images.Finder)
	}
	if qr.config.Images.Alignment != "" {
		t.Errorf("expected empty alignment path due to validation, got %s", qr.config.Images.Alignment)
	}

	// Should have 3 validation errors
	errs := qr.Validate()
	if len(errs) != 3 {
		t.Errorf("expected 3 validation errors, got %d", len(errs))
	}
}

func TestSVGGeneration(t *testing.T) {
	svg, err := New("https://example.com").
		Size(256).
		Shape("square").
		SVG()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(svg) == 0 {
		t.Error("expected non-empty SVG")
	}
	if !strings.Contains(string(svg), "<svg") {
		t.Error("expected SVG content")
	}
}

func TestSVGString(t *testing.T) {
	svgStr, err := New("test").SVGString()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(svgStr, "<svg") {
		t.Error("expected SVG string to start with <svg")
	}
}

func TestSaveAs(t *testing.T) {
	tmpFile := os.TempDir() + "/test_qr.svg"
	defer os.Remove(tmpFile)

	err := New("test").SaveAs(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("expected file to be created")
	}
}

func TestSaveAsPNG(t *testing.T) {
	err := New("test").SaveAs("test.png")

	if err == nil {
		t.Error("expected error for PNG format")
	}
	if _, ok := err.(*UnsupportedFormatError); !ok {
		t.Errorf("expected UnsupportedFormatError, got %T", err)
	}
}

func TestGetConfig(t *testing.T) {
	qr := New("test")
	cfg := qr.GetConfig()

	if cfg != qr.config {
		t.Error("expected GetConfig to return same config instance")
	}
}

func TestAllShapes(t *testing.T) {
	shapes := []string{"square", "circle", "rounded", "diamond", "dot", "star", "heart"}

	for _, shape := range shapes {
		t.Run(shape, func(t *testing.T) {
			svg, err := New("test").Shape(shape).SVG()
			if err != nil {
				t.Fatalf("unexpected error for shape %s: %v", shape, err)
			}
			if len(svg) == 0 {
				t.Errorf("expected non-empty SVG for shape %s", shape)
			}
		})
	}
}

func TestAllErrorCorrectionLevels(t *testing.T) {
	levels := []ErrorCorrectionLevel{LevelL, LevelM, LevelQ, LevelH}

	for _, level := range levels {
		t.Run(string('L'+rune(level)), func(t *testing.T) {
			svg, err := New("test").ErrorCorrection(level).SVG()
			if err != nil {
				t.Fatalf("unexpected error for level %d: %v", level, err)
			}
			if len(svg) == 0 {
				t.Errorf("expected non-empty SVG for level %d", level)
			}
		})
	}
}

package qrgode

import (
	"bytes"
	"image"
	"io"
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

func TestPNG(t *testing.T) {
	data, err := New("https://example.com").Size(256).PNG()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) < 8 {
		t.Fatalf("PNG output too small: %d bytes", len(data))
	}
	// PNG magic number: 89 50 4E 47 0D 0A 1A 0A
	magic := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}
	if !bytes.Equal(data[:8], magic) {
		t.Errorf("output does not have PNG magic bytes: %x", data[:8])
	}
}

func TestSaveAsPNG(t *testing.T) {
	tmpFile := os.TempDir() + "/test_qr.png"
	defer os.Remove(tmpFile)

	if err := New("test").SaveAs(tmpFile); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("expected non-empty PNG file")
	}
}

func TestSaveAsUnsupported(t *testing.T) {
	err := New("test").SaveAs("test.jpg")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
	if _, ok := err.(*UnsupportedFormatError); !ok {
		t.Errorf("expected UnsupportedFormatError, got %T", err)
	}
}

func TestWriteTo(t *testing.T) {
	var buf bytes.Buffer
	n, err := New("https://example.com").Size(256).WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == 0 {
		t.Error("expected non-zero bytes written")
	}
	if int64(buf.Len()) != n {
		t.Errorf("WriteTo reported %d bytes but buffer has %d", n, buf.Len())
	}
	if !bytes.HasPrefix(buf.Bytes(), []byte("<svg")) {
		t.Error("expected SVG output to begin with <svg")
	}

	// Confirm it implements io.WriterTo.
	var _ io.WriterTo = New("test")
}

func TestWriteToError(t *testing.T) {
	// Empty data should fail before writing.
	var buf bytes.Buffer
	n, err := New("").WriteTo(&buf)
	if err == nil {
		t.Error("expected error for empty data")
	}
	if n != 0 {
		t.Errorf("expected 0 bytes written on error, got %d", n)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty buffer on error, got %d bytes", buf.Len())
	}
}

func TestGetConfig(t *testing.T) {
	qr := New("test")
	cfg := qr.GetConfig()

	if cfg != qr.config {
		t.Error("expected GetConfig to return same config instance")
	}
}

func TestLogoWidth(t *testing.T) {
	qr := New("test").LogoWidth(120)
	if qr.config.Logo == nil {
		t.Fatal("Logo config not initialized")
	}
	if qr.config.Logo.Width != 120 {
		t.Errorf("Width = %d, want 120", qr.config.Logo.Width)
	}
}

func TestLogoHeight(t *testing.T) {
	qr := New("test").LogoHeight(80)
	if qr.config.Logo == nil {
		t.Fatal("Logo config not initialized")
	}
	if qr.config.Logo.Height != 80 {
		t.Errorf("Height = %d, want 80", qr.config.Logo.Height)
	}
}

func TestLogoBackground(t *testing.T) {
	qr := New("test").LogoBackground("transparent")
	if qr.config.Logo == nil {
		t.Fatal("Logo config not initialized")
	}
	if qr.config.Logo.Background != "transparent" {
		t.Errorf("Background = %q, want transparent", qr.config.Logo.Background)
	}
}

func TestUnsupportedFormatError(t *testing.T) {
	err := &UnsupportedFormatError{Format: "jpg"}
	got := err.Error()
	if !strings.Contains(got, "jpg") {
		t.Errorf("error message %q should contain format name", got)
	}
	if !strings.Contains(got, "svg") || !strings.Contains(got, "png") {
		t.Errorf("error message %q should list supported formats", got)
	}
}

func TestLogoMode(t *testing.T) {
	qr := New("test")
	qr.LogoMode(LogoOverlay)
	if qr.config.Logo == nil {
		t.Fatal("expected Logo config to be initialized")
	}
	if qr.config.Logo.Mode != LogoOverlay {
		t.Errorf("expected LogoOverlay, got %d", qr.config.Logo.Mode)
	}
}

func TestLogoPadding(t *testing.T) {
	qr := New("test").LogoPadding(0.05)
	if qr.config.Logo.Padding != 0.05 {
		t.Errorf("expected padding 0.05, got %f", qr.config.Logo.Padding)
	}
}

func TestLogoModeAffectsRendering(t *testing.T) {
	// Same logo, different modes. Overlay renders modules in the center too,
	// so the SVG path string is measurably longer.
	logo := image.NewRGBA(image.Rect(0, 0, 64, 64))
	build := func(mode LogoMode) []byte {
		qr := New("https://example.com").Size(320).ErrorCorrection(LevelH).
			LogoImage(logo).LogoDimensions(100, 100).LogoMode(mode)
		svg, err := qr.SVG()
		if err != nil {
			t.Fatalf("svg: %v", err)
		}
		return svg
	}
	excludeLen := len(build(LogoExclude))
	overlayLen := len(build(LogoOverlay))
	if overlayLen <= excludeLen {
		t.Errorf("expected overlay SVG (%d) to be longer than exclude SVG (%d)", overlayLen, excludeLen)
	}
}

func TestAllShapes(t *testing.T) {
	shapes := []Shape{ShapeSquare, ShapeCircle, ShapeRounded, ShapeDiamond, ShapeDot, ShapeStar, ShapeHeart}

	for _, shape := range shapes {
		t.Run(string(shape), func(t *testing.T) {
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

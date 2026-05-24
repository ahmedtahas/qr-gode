package qrgode

import (
	"bytes"
	"testing"
)

func TestGenerateWithOptions(t *testing.T) {
	svg, err := GenerateWithOptions("https://example.com",
		WithSize(256),
		WithQuietZone(2),
		WithErrorCorrection(LevelH),
		WithModuleShape("circle"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.HasPrefix(svg, []byte("<svg")) {
		t.Error("expected SVG output")
	}
}

func TestWithLogo_OptionSetsLogoConfig(t *testing.T) {
	logoPath := writeTempPNG(t, "logo.png")
	cfg := DefaultConfig()
	WithLogo(logoPath)(cfg)
	if cfg.Logo == nil || cfg.Logo.Path != logoPath {
		t.Errorf("WithLogo did not set Logo.Path: %+v", cfg.Logo)
	}
}

func TestWithLogoSize_OptionSetsDimensions(t *testing.T) {
	logoPath := writeTempPNG(t, "logo.png")
	cfg := DefaultConfig()
	WithLogoSize(logoPath, 120, 80)(cfg)
	if cfg.Logo == nil {
		t.Fatal("Logo nil")
	}
	if cfg.Logo.Width != 120 || cfg.Logo.Height != 80 {
		t.Errorf("dimensions = %dx%d, want 120x80", cfg.Logo.Width, cfg.Logo.Height)
	}
}

func TestGenerateWithOptions_WithLogo(t *testing.T) {
	logoPath := writeTempPNG(t, "logo.png")
	svg, err := GenerateWithOptions("https://example.com",
		WithSize(256),
		WithErrorCorrection(LevelH),
		WithLogo(logoPath),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(svg) == 0 {
		t.Error("expected non-empty SVG")
	}
}

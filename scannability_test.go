package qrgode

import (
	"strings"
	"testing"
)

func TestScannabilityWarnings_NoLogo(t *testing.T) {
	qr := New("https://example.com")
	if w := qr.ScannabilityWarnings(); w != nil {
		t.Errorf("expected no warnings without a logo, got %v", w)
	}
}

func TestScannabilityWarnings_AutoSizedLogoLowECL(t *testing.T) {
	// Auto-sized logo ≈ 22.5% × 22.5% ≈ 5% of area — above L's 5% limit,
	// equal-ish at M, below Q/H.
	qr := New("https://example.com").Logo("nonexistent.png") // path validation fails but Logo struct exists
	// Force path on directly to bypass the validator failing the path setter.
	qr.config.Logo = &LogoConfig{Path: "auto.png"}

	qr.ErrorCorrection(LevelL)
	if w := qr.ScannabilityWarnings(); len(w) == 0 {
		t.Error("expected warning at LevelL with auto-sized logo")
	}

	qr.ErrorCorrection(LevelH)
	if w := qr.ScannabilityWarnings(); len(w) != 0 {
		t.Errorf("expected no warning at LevelH with auto-sized logo, got %v", w)
	}
}

func TestScannabilityWarnings_OversizedLogo(t *testing.T) {
	qr := New("https://example.com").Size(512).ErrorCorrection(LevelH)
	// Manually set a logo path + huge dimensions (skip Logo() validation).
	qr.config.Logo = &LogoConfig{
		Path:   "anything.png",
		Width:  300,
		Height: 300,
	}
	// 300² / 512² ≈ 34% — exceeds H's 23% limit.
	warnings := qr.ScannabilityWarnings()
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d: %v", len(warnings), warnings)
	}
	if !strings.Contains(warnings[0], "ECL H") {
		t.Errorf("expected warning to mention ECL H, got %q", warnings[0])
	}
	if !strings.Contains(warnings[0], "%") {
		t.Errorf("expected warning to include a percentage, got %q", warnings[0])
	}
}

func TestScannabilityWarnings_SafeExplicitLogo(t *testing.T) {
	qr := New("test").Size(512).ErrorCorrection(LevelQ)
	qr.config.Logo = &LogoConfig{
		Path:   "anything.png",
		Width:  150,
		Height: 150,
	}
	// 150² / 512² ≈ 8.6% — well under Q's 18% limit.
	if w := qr.ScannabilityWarnings(); w != nil {
		t.Errorf("expected no warning for safe logo at LevelQ, got %v", w)
	}
}

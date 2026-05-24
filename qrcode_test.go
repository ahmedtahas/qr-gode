package qrgode

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerate_NilConfig(t *testing.T) {
	svg, err := Generate("hello", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.HasPrefix(svg, []byte("<svg")) {
		t.Error("expected SVG output")
	}
}

func TestGenerate_EmptyData(t *testing.T) {
	_, err := Generate("", nil)
	if err == nil {
		t.Error("expected error for empty data")
	}
	if _, ok := err.(*ValidationError); !ok {
		t.Errorf("expected *ValidationError, got %T", err)
	}
}

func TestGenerate_InvalidConfig(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Size = -1 // forces ValidateConfig to fail
	if _, err := Generate("hi", cfg); err == nil {
		t.Error("expected error for invalid config")
	}
}

func TestGeneratePNG_Basic(t *testing.T) {
	data, err := GeneratePNG("https://example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) < 8 || !bytes.Equal(data[:8], []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}) {
		t.Errorf("output missing PNG magic: %x", data[:min(8, len(data))])
	}
}

func TestGeneratePNG_EmptyData(t *testing.T) {
	if _, err := GeneratePNG("", nil); err == nil {
		t.Error("expected error for empty data")
	}
}

func TestGeneratePNG_InvalidConfig(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Size = -1
	if _, err := GeneratePNG("hi", cfg); err == nil {
		t.Error("expected error for invalid config")
	}
}

func TestGenerateToFile_SVG(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.svg")
	if err := GenerateToFile("test", nil, path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Size() == 0 {
		t.Errorf("file missing or empty: %v", err)
	}
}

func TestGenerateToFile_PNG(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.png")
	if err := GenerateToFile("test", nil, path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Size() == 0 {
		t.Errorf("file missing or empty: %v", err)
	}
}

func TestGenerateToFile_NoExtensionDefaultsToSVG(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out")
	if err := GenerateToFile("test", nil, path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGenerateToFile_UnsupportedExtension(t *testing.T) {
	err := GenerateToFile("test", nil, filepath.Join(t.TempDir(), "out.bmp"))
	if err == nil {
		t.Error("expected error for unsupported extension")
	}
	if _, ok := err.(*UnsupportedFormatError); !ok {
		t.Errorf("expected *UnsupportedFormatError, got %T", err)
	}
}

func TestGenerateToFile_PropagatesGenerateError(t *testing.T) {
	// Empty data fails inside Generate — wrapper should propagate.
	if err := GenerateToFile("", nil, filepath.Join(t.TempDir(), "x.svg")); err == nil {
		t.Error("expected error from underlying Generate")
	}
	if err := GenerateToFile("", nil, filepath.Join(t.TempDir(), "x.png")); err == nil {
		t.Error("expected error from underlying GeneratePNG")
	}
}

// Encoding long data forces version 7+, which exercises the version-info
// reservation path in the encoder.
func TestGenerate_Version7Plus(t *testing.T) {
	longData := strings.Repeat("A", 200)
	svg, err := Generate(longData, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(svg) == 0 {
		t.Error("expected non-empty SVG for v7+ QR")
	}
}

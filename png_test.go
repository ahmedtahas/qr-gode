package qrgode

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

func TestPNG_WithLogoImage(t *testing.T) {
	logo := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := range 64 {
		for x := range 64 {
			logo.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	data, err := New("https://example.com").
		Size(320).
		ErrorCorrection(LevelH).
		LogoImage(logo).
		PNG()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) < 8 || !bytes.Equal(data[:8], []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}) {
		t.Error("expected PNG magic bytes")
	}
}

func TestPNG_WithLogoFromFile(t *testing.T) {
	logoPath := writeTempPNG(t, "logo.png")
	data, err := New("https://example.com").
		Size(320).
		ErrorCorrection(LevelH).
		Logo(logoPath).
		PNG()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("empty PNG")
	}
}

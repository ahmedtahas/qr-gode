package qrgode

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// writeTempPNG creates a small PNG fixture and returns its path. Used by tests
// that need real bytes for file-stat validation.
func writeTempPNG(t *testing.T, name string) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := range 32 {
		for x := range 32 {
			img.Set(x, y, color.Black)
		}
	}
	path := filepath.Join(t.TempDir(), name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode: %v", err)
	}
	return path
}

package qrgode

import (
	"strings"
	"testing"
)

func TestSVG_WithCustomImages(t *testing.T) {
	finder := writeTempPNG(t, "finder.png")
	align := writeTempPNG(t, "align.png")
	mod := writeTempPNG(t, "module.png")

	svg, err := New("https://example.com").
		Size(320).
		FinderImage(finder).
		AlignmentImage(align).
		ModuleImage(mod).
		SVG()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Custom images render as <image href="..."> tags in the SVG. The exact
	// path will be base64-encoded, but the marker should be present.
	if !strings.Contains(string(svg), "<image") {
		t.Errorf("expected <image> tags in SVG with custom images, got: %s", string(svg)[:min(200, len(svg))])
	}
}

func TestSVG_WithFinderImageOnly(t *testing.T) {
	finder := writeTempPNG(t, "finder.png")
	svg, err := New("test").FinderImage(finder).SVG()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(svg), "<image") {
		t.Error("expected <image> tag for custom finder")
	}
}

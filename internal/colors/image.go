package colors

import (
	"fmt"
	"image"
)

// ImageSampler samples colors from an image file.
type ImageSampler struct {
	Path  string
	img   image.Image
	ready bool
}

// NewImageSampler creates a color sampler from an image path.
func NewImageSampler(path string) *ImageSampler {
	return &ImageSampler{
		Path: path,
	}
}

// Load loads the image from disk.
func (s *ImageSampler) Load() error {
	// TODO:
	// 1. Open file
	// 2. Decode image (PNG, JPEG, etc.)
	// 3. Store in s.img
	// 4. Set s.ready = true
	return nil
}

func (s *ImageSampler) ColorAt(x, y float64) string {
	// TODO:
	// 1. Map normalized coordinates to image pixels
	// 2. Sample pixel color
	// 3. Return as hex string
	if !s.ready {
		return "#000000"
	}
	return "#000000"
}

func (s *ImageSampler) Type() string {
	return "image"
}

func (s *ImageSampler) SVGDefs(id string) string {
	// Note: For image sampling, we don't use SVG patterns.
	// Instead, each module gets its own fill color.
	// This method could define a pattern if we wanted tiled images.
	return ""
}

func (s *ImageSampler) SVGFill(id string) string {
	// This won't be used directly since each module
	// gets individual ColorAt() calls
	return fmt.Sprintf("url(#%s)", id)
}

// rgbToHex converts RGB values to hex string.
func rgbToHex(r, g, b uint8) string {
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

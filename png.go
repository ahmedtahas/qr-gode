package qrgode

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// rasterizeSVGToPNG converts SVG bytes to a PNG of the given pixel size.
// It uses oksvg + rasterx (pure Go) so PNG output inherits every styling
// feature from the SVG pipeline — gradients, custom shapes, logos, etc.
func rasterizeSVGToPNG(svg []byte, width, height int) ([]byte, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("png: invalid dimensions %dx%d", width, height)
	}

	icon, err := oksvg.ReadIconStream(bytes.NewReader(svg))
	if err != nil {
		return nil, fmt.Errorf("png: parse svg: %w", err)
	}
	icon.SetTarget(0, 0, float64(width), float64(height))

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(width, height, img, img.Bounds())
	raster := rasterx.NewDasher(width, height, scanner)
	icon.Draw(raster, 1.0)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("png: encode: %w", err)
	}
	return buf.Bytes(), nil
}

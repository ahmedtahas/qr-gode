package qrgode

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	xdraw "golang.org/x/image/draw"
)

// renderPNG generates the QR code as PNG bytes. It rasterizes the SVG
// pipeline via oksvg+rasterx, then composites the logo image onto the
// result — oksvg doesn't rasterize embedded <image href="data:..."> data
// URIs, so the logo has to be drawn manually for it to appear.
func (r *renderer) renderPNG() ([]byte, error) {
	svg, err := r.renderSVG()
	if err != nil {
		return nil, err
	}
	width, height := r.config.Size, r.config.Size
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("png: invalid dimensions %dx%d", width, height)
	}

	img, err := rasterizeSVGData(svg, width, height)
	if err != nil {
		return nil, err
	}

	if r.hasLogo() {
		if err := r.compositeLogoOnto(img); err != nil {
			return nil, err
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("png: encode: %w", err)
	}
	return buf.Bytes(), nil
}

// rasterizeSVGData rasterizes raw SVG bytes into an RGBA bitmap.
func rasterizeSVGData(svg []byte, width, height int) (*image.RGBA, error) {
	icon, err := oksvg.ReadIconStream(bytes.NewReader(svg))
	if err != nil {
		return nil, fmt.Errorf("png: parse svg: %w", err)
	}
	icon.SetTarget(0, 0, float64(width), float64(height))
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(width, height, img, img.Bounds())
	raster := rasterx.NewDasher(width, height, scanner)
	icon.Draw(raster, 1.0)
	return img, nil
}

// compositeLogoOnto resizes and draws the configured logo onto the target.
func (r *renderer) compositeLogoOnto(target *image.RGBA) error {
	logoImg, err := r.loadLogoForRaster()
	if err != nil {
		return fmt.Errorf("png: load logo: %w", err)
	}
	if logoImg == nil {
		return nil
	}
	logoW, logoH, _, err := r.calculateLogoDimensions()
	if err != nil {
		return err
	}
	qrSize := float64(r.config.Size)
	logoX := int((qrSize - logoW) / 2)
	logoY := int((qrSize - logoH) / 2)
	dst := image.Rect(logoX, logoY, logoX+int(logoW), logoY+int(logoH))
	xdraw.CatmullRom.Scale(target, dst, logoImg, logoImg.Bounds(), xdraw.Over, nil)
	return nil
}

// loadLogoForRaster returns the logo as an image.Image, decoding files
// (PNG/JPG) or rasterizing SVG logos via the same oksvg pipeline.
func (r *renderer) loadLogoForRaster() (image.Image, error) {
	logo := r.config.Logo
	if logo.Image != nil {
		return logo.Image, nil
	}
	if logo.Path == "" {
		return nil, nil
	}
	if strings.HasSuffix(strings.ToLower(logo.Path), ".svg") {
		data, err := os.ReadFile(logo.Path)
		if err != nil {
			return nil, err
		}
		// Rasterize at logo dimensions so scaling artifacts are minimal.
		w, h, _, err := r.calculateLogoDimensions()
		if err != nil || w <= 0 || h <= 0 {
			w, h = 256, 256
		}
		return rasterizeSVGData(data, int(w), int(h))
	}
	f, err := os.Open(logo.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

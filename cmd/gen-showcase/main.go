// Generates the README showcase SVGs into assets/showcase/.
// Run from the repo root: go run ./cmd/gen-showcase
package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"

	qrgode "github.com/ahmedtahas/qr-gode"
)

type sample struct {
	name    string
	caption string
	build   func() *qrgode.QRCode
}

// makeLogo returns a 192×192 placeholder logo: a purple ring on transparent
// background. Used by the "with logo" samples so the showcase doesn't depend
// on committing a binary logo asset.
func makeLogo() image.Image {
	const size = 192
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	cx, cy := size/2, size/2
	outerR2 := (size/2 - 4) * (size/2 - 4)
	innerR2 := (size / 3) * (size / 3)
	ring := color.RGBA{108, 92, 231, 255} // #6c5ce7
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx, dy := x-cx, y-cy
			d2 := dx*dx + dy*dy
			if d2 <= outerR2 && d2 >= innerR2 {
				img.Set(x, y, ring)
			}
		}
	}
	return img
}

func main() {
	const data = "https://github.com/ahmedtahas/qr-gode"
	const size = 320

	logo := makeLogo()

	samples := []sample{
		{"01-square", "Classic square", func() *qrgode.QRCode {
			return qrgode.New(data).Size(size).Shape(qrgode.ShapeSquare)
		}},
		{"02-circle", "Circle modules", func() *qrgode.QRCode {
			return qrgode.New(data).Size(size).Shape(qrgode.ShapeCircle).Foreground("#3498db")
		}},
		{"03-rounded-linear", "Rounded · linear gradient", func() *qrgode.QRCode {
			return qrgode.New(data).Size(size).Shape(qrgode.ShapeRounded).
				LinearGradient(45, "#ff6b6b", "#feca57")
		}},
		{"04-dot-radial", "Dot · radial gradient", func() *qrgode.QRCode {
			return qrgode.New(data).Size(size).Shape(qrgode.ShapeDot).
				RadialGradient(0.5, 0.5, "#4ecdc4", "#1a535c")
		}},
		{"05-heart", "Heart modules", func() *qrgode.QRCode {
			return qrgode.New(data).Size(size).Shape(qrgode.ShapeHeart).Foreground("#e84393")
		}},
		{"06-star", "Star · gradient", func() *qrgode.QRCode {
			return qrgode.New(data).Size(size).Shape(qrgode.ShapeStar).
				LinearGradient(90, "#a29bfe", "#6c5ce7")
		}},
		{"07-diamond", "Diamond modules", func() *qrgode.QRCode {
			return qrgode.New(data).Size(size).Shape(qrgode.ShapeDiamond).Foreground("#27ae60")
		}},
		{"08-darkmode", "Dark mode", func() *qrgode.QRCode {
			return qrgode.New(data).Size(size).Shape(qrgode.ShapeRounded).
				Background("#0f1419").Foreground("#e6e6e6")
		}},
		{"09-logo-classic", "Logo · classic", func() *qrgode.QRCode {
			return qrgode.New(data).Size(size).Shape(qrgode.ShapeRounded).
				LinearGradient(45, "#4ecdc4", "#6c5ce7").
				ErrorCorrection(qrgode.LevelH).LogoImage(logo)
		}},
		{"10-logo-overlay", "Logo · overlay mode", func() *qrgode.QRCode {
			return qrgode.New(data).Size(size).Shape(qrgode.ShapeCircle).
				Foreground("#2d3436").
				ErrorCorrection(qrgode.LevelH).LogoImage(logo).
				LogoMode(qrgode.LogoOverlay).LogoBackground("transparent")
		}},
	}

	outDir := filepath.Join("assets", "showcase")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("mkdir %s: %v", outDir, err)
	}

	for _, s := range samples {
		path := filepath.Join(outDir, s.name+".svg")
		if err := s.build().SaveAs(path); err != nil {
			log.Fatalf("render %s: %v", s.name, err)
		}
		fmt.Printf("wrote %s — %s\n", path, s.caption)
	}
}

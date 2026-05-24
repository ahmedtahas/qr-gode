// Generates the README showcase SVGs into assets/showcase/.
// Run from the repo root: go run ./cmd/gen-showcase
package main

import (
	"fmt"
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

func main() {
	const data = "https://github.com/ahmedtahas/qr-gode"
	const size = 320

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

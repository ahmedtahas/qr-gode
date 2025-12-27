package qrcode_test

import (
	"fmt"
	"log"

	"github.com/ahmedtahas/qr-gode/pkg/qrcode"
)

func Example_basic() {
	// Simple QR code generation
	svg, err := qrcode.Generate("https://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated %d bytes of SVG\n", len(svg))
}

func Example_builder() {
	// Using the builder pattern for customization
	qr := qrcode.New("https://example.com").
		Size(512).
		ErrorCorrection(qrcode.LevelH).
		Shape("circle").
		Foreground("#3498db").
		Background("#ffffff")

	svg, err := qr.SVG()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated %d bytes of SVG\n", len(svg))
}

func Example_gradient() {
	// Linear gradient
	qr := qrcode.New("https://example.com").
		Size(512).
		LinearGradient(45, "#ff0000", "#00ff00", "#0000ff").
		Shape("rounded")

	svg, err := qr.SVG()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated %d bytes of SVG\n", len(svg))
}

func Example_radialGradient() {
	// Radial gradient from center
	qr := qrcode.New("https://example.com").
		Size(512).
		RadialGradient(0.5, 0.5, "#ff6b6b", "#4ecdc4").
		Shape("dot")

	svg, err := qr.SVG()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated %d bytes of SVG\n", len(svg))
}

func Example_customImages() {
	// Using custom images for QR elements
	qr := qrcode.New("https://example.com").
		Size(512).
		ModuleImage("dot.png").
		FinderImage("finder.png").
		AlignmentImage("align.png")

	// Note: This would fail without actual image files
	_, _ = qr.SVG()
}

func Example_saveToFile() {
	// Save directly to file
	err := qrcode.New("https://example.com").
		Size(1024).
		Shape("diamond").
		Foreground("#2c3e50").
		SaveAs("my_qr.svg")

	if err != nil {
		log.Fatal(err)
	}
}

func Example_advancedConfig() {
	// For advanced usage, access the underlying config
	qr := qrcode.New("https://example.com")
	cfg := qr.GetConfig()

	// Modify config directly
	cfg.Size = 800
	cfg.QuietZone = 2
	cfg.Modules.Shape = "star"
	cfg.Modules.Color = qrcode.NewLinearGradientColor(90, []string{"#e74c3c", "#9b59b6", "#3498db"})

	svg, err := qr.SVG()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated %d bytes of SVG\n", len(svg))
}

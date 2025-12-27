package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ahmedtahas/qr-gode/internal/colors"
	"github.com/ahmedtahas/qr-gode/pkg/qrcode"
)

func main() {
	// Flags
	output := flag.String("o", "qrcode.svg", "Output file path")
	size := flag.Int("size", 512, "Output size in pixels")
	shape := flag.String("shape", "square", "Module shape: square, circle, rounded, diamond, dot, star, heart")
	fgColor := flag.String("fg", "#000000", "Foreground color (hex)")
	bgColor := flag.String("bg", "#FFFFFF", "Background color (hex)")
	gradient := flag.String("gradient", "", "Gradient colors (comma-separated, e.g. '#ff0000,#0000ff')")
	gradientAngle := flag.Float64("gradient-angle", 45, "Gradient angle in degrees")
	radial := flag.Bool("radial", false, "Use radial gradient instead of linear")
	ecl := flag.String("ecl", "M", "Error correction level: L, M, Q, H")

	// Custom image flags
	moduleImg := flag.String("module-img", "", "Custom PNG/JPG for data modules")
	finderImg := flag.String("finder-img", "", "Custom PNG/JPG for finder pattern modules")
	alignImg := flag.String("align-img", "", "Custom PNG/JPG for alignment pattern modules")

	// Logo flags
	logoImg := flag.String("logo", "", "Logo image to place in center (PNG/JPG/SVG)")
	logoWidth := flag.Int("logo-width", 0, "Optional: logo width in pixels (0 = auto)")
	logoHeight := flag.Int("logo-height", 0, "Optional: logo height in pixels (0 = auto)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: qr-gode [options] <data>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  qr-gode 'https://google.com'\n")
		fmt.Fprintf(os.Stderr, "  qr-gode -shape circle -fg '#3498db' 'Hello World'\n")
		fmt.Fprintf(os.Stderr, "  qr-gode -gradient '#ff6b6b,#4ecdc4' -shape rounded 'Gradient QR'\n")
		fmt.Fprintf(os.Stderr, "  qr-gode -module-img dot.png -finder-img finder.png 'Custom Images'\n")
		fmt.Fprintf(os.Stderr, "  qr-gode -logo logo.png 'QR with Logo'\n")
		fmt.Fprintf(os.Stderr, "  qr-gode -logo logo.png -logo-width 100 'QR with custom logo size'\n")
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	data := flag.Arg(0)

	// Build config
	cfg := qrcode.DefaultConfig()
	cfg.Size = *size
	cfg.Modules.Shape = *shape
	cfg.Background = colors.NewSolid(*bgColor)

	// Set error correction level
	switch strings.ToUpper(*ecl) {
	case "L":
		cfg.ErrorCorrection = qrcode.LevelL
	case "M":
		cfg.ErrorCorrection = qrcode.LevelM
	case "Q":
		cfg.ErrorCorrection = qrcode.LevelQ
	case "H":
		cfg.ErrorCorrection = qrcode.LevelH
	}

	// Set color (gradient or solid)
	if *gradient != "" {
		stops := strings.Split(*gradient, ",")
		for i := range stops {
			stops[i] = strings.TrimSpace(stops[i])
		}
		if *radial {
			cfg.Modules.Color = colors.NewRadialGradient(0.5, 0.5, stops)
		} else {
			cfg.Modules.Color = colors.NewLinearGradient(*gradientAngle, stops)
		}
	} else {
		cfg.Modules.Color = colors.NewSolid(*fgColor)
	}

	// Set custom images if provided
	if *moduleImg != "" || *finderImg != "" || *alignImg != "" {
		cfg.Images = &qrcode.CustomImages{
			Module:    *moduleImg,
			Finder:    *finderImg,
			Alignment: *alignImg,
		}
	}

	// Set logo if provided
	if *logoImg != "" {
		cfg.Logo = &qrcode.LogoConfig{
			Path:   *logoImg,
			Width:  *logoWidth,
			Height: *logoHeight,
		}
	}

	// Generate
	err := qrcode.GenerateToFile(data, cfg, *output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated QR code for %q -> %s\n", data, *output)
}

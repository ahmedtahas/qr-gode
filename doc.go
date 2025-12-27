// Package qrgode provides QR code generation with extensive customization.
//
// # Quick Start
//
// The simplest way to generate a QR code:
//
//	svg, err := qrgode.Generate("https://example.com", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	os.WriteFile("qr.svg", svg, 0644)
//
// Or save directly to a file:
//
//	err := qrgode.GenerateToFile("https://example.com", nil, "qr.svg")
//
// # Builder Pattern
//
// For customization, use the fluent builder API:
//
//	qr := qrgode.New("https://example.com").
//		Size(512).
//		ErrorCorrection(qrgode.LevelH).
//		Shape("circle").
//		Foreground("#3498db").
//		Background("#ffffff")
//
//	svg, err := qr.SVG()
//	// or
//	err := qr.SaveAs("qr.svg")
//
// # Gradients
//
// Apply linear or radial gradients to modules:
//
//	// Linear gradient at 45 degrees
//	qr := qrgode.New("https://example.com").
//		LinearGradient(45, "#ff0000", "#00ff00", "#0000ff")
//
//	// Radial gradient from center
//	qr := qrgode.New("https://example.com").
//		RadialGradient(0.5, 0.5, "#ff0000", "#0000ff")
//
// # Custom Images
//
// Use custom PNG/JPG images for QR elements:
//
//	qr := qrgode.New("https://example.com").
//		ModuleImage("dot.png").        // Custom image for data modules
//		FinderImage("finder.png").     // Custom image for finder patterns (7x7)
//		AlignmentImage("align.png")    // Custom image for alignment patterns (5x5)
//
// # Logo
//
// Add a logo in the center of the QR code:
//
//	// From file path
//	qr := qrgode.New("https://example.com").
//		Logo("logo.png")
//
//	// From image.Image in memory
//	qr := qrgode.New("https://example.com").
//		LogoImage(myImage)
//
// # Shapes
//
// Available module shapes:
//   - "square" (default)
//   - "circle"
//   - "rounded" (rounded corners)
//   - "diamond"
//   - "dot" (smaller circle)
//   - "star"
//   - "heart"
//
// # Error Correction Levels
//
// Higher levels provide more redundancy but result in denser QR codes:
//   - LevelL: ~7% recovery (smallest QR)
//   - LevelM: ~15% recovery (default)
//   - LevelQ: ~25% recovery
//   - LevelH: ~30% recovery (densest QR)
package qrgode

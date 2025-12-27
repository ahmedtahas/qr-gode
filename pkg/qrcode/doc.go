// Package qrcode provides QR code generation with extensive customization.
//
// # Quick Start
//
// The simplest way to generate a QR code:
//
//	svg, err := qrcode.Generate("https://example.com", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	os.WriteFile("qr.svg", svg, 0644)
//
// Or save directly to a file:
//
//	err := qrcode.GenerateToFile("https://example.com", nil, "qr.svg")
//
// # Builder Pattern
//
// For customization, use the fluent builder API:
//
//	qr := qrcode.New("https://example.com").
//		Size(512).
//		ErrorCorrection(qrcode.LevelH).
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
//	qr := qrcode.New("https://example.com").
//		LinearGradient(45, "#ff0000", "#00ff00", "#0000ff")
//
//	// Radial gradient from center
//	qr := qrcode.New("https://example.com").
//		RadialGradient(0.5, 0.5, "#ff0000", "#0000ff")
//
// # Custom Images
//
// Use custom PNG/JPG images for QR elements:
//
//	qr := qrcode.New("https://example.com").
//		ModuleImage("dot.png").        // Custom image for data modules
//		FinderImage("finder.png").     // Custom image for finder patterns (7x7)
//		AlignmentImage("align.png")    // Custom image for alignment patterns (5x5)
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
package qrcode

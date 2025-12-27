package elements

// LogoRenderer handles logo embedding in the QR code.
type LogoRenderer struct {
	Path       string  // Path to logo image
	Size       float64 // Size as fraction of QR (0.0-1.0)
	Padding    float64 // Padding around logo
	Background string  // Background color (hex or "transparent")
	QRSize     float64 // Total QR code size in SVG units
}

// Render renders the logo overlay.
func (r *LogoRenderer) Render() string {
	// TODO:
	// 1. Calculate logo dimensions based on Size * QRSize
	// 2. Calculate center position
	// 3. Render background rectangle if not transparent
	// 4. Render logo image (embedded or linked)
	return ""
}

// GetClearZone returns the rectangle that should not have modules.
// Returns (x, y, width, height) in module coordinates.
func (r *LogoRenderer) GetClearZone(matrixSize int) (int, int, int, int) {
	// TODO:
	// Calculate which modules should be skipped
	// to make room for the logo
	return 0, 0, 0, 0
}

// embedImage converts an image file to base64 data URI.
func embedImage(path string) (string, error) {
	// TODO:
	// 1. Read file
	// 2. Detect MIME type
	// 3. Base64 encode
	// 4. Return data:image/png;base64,... URI
	return "", nil
}

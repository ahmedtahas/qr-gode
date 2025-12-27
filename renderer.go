package qrgode

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ahmedtahas/qr-gode/internal/encoder"
	"github.com/ahmedtahas/qr-gode/internal/shapes"
)

// Logo size constraints (as fraction of QR size)
const (
	logoMinSize = 0.15 // Minimum 15% of QR
	logoMaxSize = 0.30 // Maximum 30% of QR
)

// hasLogo returns true if a logo is configured (either as image or path)
func (r *Renderer) hasLogo() bool {
	logo := r.config.Logo
	return logo != nil && (logo.Image != nil || logo.Path != "")
}

// getLogoSourceDimensions returns the original dimensions of the logo source
func (r *Renderer) getLogoSourceDimensions() (int, int, error) {
	logo := r.config.Logo
	if logo.Image != nil {
		bounds := logo.Image.Bounds()
		return bounds.Dx(), bounds.Dy(), nil
	}
	return getLogoDimensions(logo.Path)
}

// calculateLogoDimensions returns logo width, height, and padding in pixels
func (r *Renderer) calculateLogoDimensions() (float64, float64, float64, error) {
	if !r.hasLogo() {
		return 0, 0, 0, nil
	}

	logo := r.config.Logo
	qrSize := float64(r.config.Size)
	var logoWidth, logoHeight float64

	if logo.Width > 0 && logo.Height > 0 {
		logoWidth = float64(logo.Width)
		logoHeight = float64(logo.Height)
	} else if logo.Width > 0 {
		imgW, imgH, err := r.getLogoSourceDimensions()
		if err != nil {
			return 0, 0, 0, err
		}
		logoWidth = float64(logo.Width)
		logoHeight = logoWidth * float64(imgH) / float64(imgW)
	} else if logo.Height > 0 {
		imgW, imgH, err := r.getLogoSourceDimensions()
		if err != nil {
			return 0, 0, 0, err
		}
		logoHeight = float64(logo.Height)
		logoWidth = logoHeight * float64(imgW) / float64(imgH)
	} else {
		imgW, imgH, err := r.getLogoSourceDimensions()
		if err != nil {
			return 0, 0, 0, err
		}

		targetFraction := (logoMinSize + logoMaxSize) / 2
		maxDimension := qrSize * targetFraction
		aspectRatio := float64(imgW) / float64(imgH)

		if aspectRatio >= 1 {
			logoWidth = maxDimension
			logoHeight = maxDimension / aspectRatio
		} else {
			logoHeight = maxDimension
			logoWidth = maxDimension * aspectRatio
		}
	}

	padding := logoWidth
	if logoHeight > logoWidth {
		padding = logoHeight
	}
	padding *= 0.1

	return logoWidth, logoHeight, padding, nil
}

// Renderer converts a QR matrix to SVG output.
type Renderer struct {
	config *Config
	matrix *encoder.Matrix
}

// NewRenderer creates a renderer for the given matrix and config.
func NewRenderer(matrix *encoder.Matrix, config *Config) *Renderer {
	return &Renderer{
		config: config,
		matrix: matrix,
	}
}

// RenderSVG generates the SVG representation of the QR code.
func (r *Renderer) RenderSVG() ([]byte, error) {
	// Check if using custom images
	if r.config.Images != nil && (r.config.Images.Module != "" || r.config.Images.Finder != "" || r.config.Images.Alignment != "") {
		return r.renderWithImages()
	}
	return r.renderWithShapes()
}

// renderWithShapes renders QR code using vector shapes
func (r *Renderer) renderWithShapes() ([]byte, error) {
	matrixSize := r.matrix.Size()
	quietZone := r.config.QuietZone
	totalModules := matrixSize + 2*quietZone

	// Calculate module size based on desired output size
	moduleSize := float64(r.config.Size) / float64(totalModules)

	// Calculate logo exclusion zone (in module coordinates)
	var logoMinX, logoMinY, logoMaxX, logoMaxY int
	hasLogoZone := false
	if r.hasLogo() {
		logoWidth, logoHeight, padding, err := r.calculateLogoDimensions()
		if err != nil {
			return nil, fmt.Errorf("failed to calculate logo dimensions: %w", err)
		}

		// Total logo area including padding
		totalWidth := logoWidth + 2*padding
		totalHeight := logoHeight + 2*padding

		// Convert pixel dimensions to module counts
		excludeHalfX := int(totalWidth/moduleSize/2) + 1
		excludeHalfY := int(totalHeight/moduleSize/2) + 1

		// Center of SVG in matrix coordinates
		svgCenterModule := float64(totalModules) / 2
		matrixCenterX := int(svgCenterModule) - quietZone
		matrixCenterY := int(svgCenterModule) - quietZone

		logoMinX = matrixCenterX - excludeHalfX
		logoMinY = matrixCenterY - excludeHalfY
		logoMaxX = matrixCenterX + excludeHalfX
		logoMaxY = matrixCenterY + excludeHalfY
		hasLogoZone = true
	}

	var buf bytes.Buffer

	// SVG header
	fmt.Fprintf(&buf, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`,
		r.config.Size, r.config.Size, r.config.Size, r.config.Size)
	buf.WriteString("\n")

	// Defs section for gradients
	if r.config.Modules.Color != nil {
		defs := r.config.Modules.Color.SVGDefs("module-fill")
		if defs != "" {
			fmt.Fprintf(&buf, "<defs>%s</defs>\n", defs)
		}
	}

	// Background
	bgColor := "#FFFFFF"
	if r.config.Background != nil {
		bgColor = r.config.Background.SVGFill("")
	}
	fmt.Fprintf(&buf, `<rect width="100%%" height="100%%" fill="%s"/>`, bgColor)
	buf.WriteString("\n")

	// Get module color/fill
	moduleFill := "#000000"
	if r.config.Modules.Color != nil {
		moduleFill = r.config.Modules.Color.SVGFill("module-fill")
	}

	// Get shape
	shapeName := r.config.Modules.Shape
	if shapeName == "" {
		shapeName = "square"
	}
	shape := shapes.Get(shapeName)
	if shape == nil {
		shape = shapes.Get("square")
	}

	// Render all dark modules as a single path for efficiency
	fmt.Fprintf(&buf, `<path fill="%s" d="`, moduleFill)

	shapePath := shape.SVGPath()
	for y := 0; y < matrixSize; y++ {
		for x := 0; x < matrixSize; x++ {
			// Skip modules in the logo zone
			if hasLogoZone && x >= logoMinX && x <= logoMaxX && y >= logoMinY && y <= logoMaxY {
				continue
			}

			if r.matrix.Get(x, y).Dark {
				// Calculate position with quiet zone offset
				px := float64(quietZone+x) * moduleSize
				py := float64(quietZone+y) * moduleSize

				// Transform and add shape path
				transformed := transformPath(shapePath, px, py, moduleSize)
				buf.WriteString(transformed)
				buf.WriteString(" ")
			}
		}
	}

	buf.WriteString(`"/>`)
	buf.WriteString("\n")

	// Render logo if configured
	if r.hasLogo() {
		logoSVG, err := r.renderLogo()
		if err != nil {
			return nil, err
		}
		buf.WriteString(logoSVG)
	}

	// Close SVG
	buf.WriteString("</svg>")

	return buf.Bytes(), nil
}

// renderWithImages renders QR code using custom PNG images
func (r *Renderer) renderWithImages() ([]byte, error) {
	matrixSize := r.matrix.Size()
	quietZone := r.config.QuietZone
	totalModules := matrixSize + 2*quietZone

	moduleSize := float64(r.config.Size) / float64(totalModules)

	// Calculate logo exclusion zone (in module coordinates)
	var logoMinX, logoMinY, logoMaxX, logoMaxY int
	hasLogoZone := false
	if r.hasLogo() {
		logoWidth, logoHeight, padding, err := r.calculateLogoDimensions()
		if err != nil {
			return nil, fmt.Errorf("failed to calculate logo dimensions: %w", err)
		}

		// Total logo area including padding
		totalWidth := logoWidth + 2*padding
		totalHeight := logoHeight + 2*padding

		// Convert pixel dimensions to module counts
		excludeHalfX := int(totalWidth/moduleSize/2) + 1
		excludeHalfY := int(totalHeight/moduleSize/2) + 1

		// Center of SVG in matrix coordinates
		svgCenterModule := float64(totalModules) / 2
		matrixCenterX := int(svgCenterModule) - quietZone
		matrixCenterY := int(svgCenterModule) - quietZone

		logoMinX = matrixCenterX - excludeHalfX
		logoMinY = matrixCenterY - excludeHalfY
		logoMaxX = matrixCenterX + excludeHalfX
		logoMaxY = matrixCenterY + excludeHalfY
		hasLogoZone = true
	}

	var buf bytes.Buffer

	// SVG header with xlink namespace for images
	fmt.Fprintf(&buf, `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 %d %d" width="%d" height="%d">`,
		r.config.Size, r.config.Size, r.config.Size, r.config.Size)
	buf.WriteString("\n")

	// Background
	bgColor := "#FFFFFF"
	if r.config.Background != nil {
		bgColor = r.config.Background.SVGFill("")
	}
	fmt.Fprintf(&buf, `<rect width="100%%" height="100%%" fill="%s"/>`, bgColor)
	buf.WriteString("\n")

	// Load and embed images as base64
	var moduleImg, finderImg, alignImg string
	var err error

	if r.config.Images.Module != "" {
		moduleImg, err = loadImageAsDataURI(r.config.Images.Module)
		if err != nil {
			return nil, fmt.Errorf("failed to load module image: %w", err)
		}
	}
	if r.config.Images.Finder != "" {
		finderImg, err = loadImageAsDataURI(r.config.Images.Finder)
		if err != nil {
			return nil, fmt.Errorf("failed to load finder image: %w", err)
		}
	}
	if r.config.Images.Alignment != "" {
		alignImg, err = loadImageAsDataURI(r.config.Images.Alignment)
		if err != nil {
			return nil, fmt.Errorf("failed to load alignment image: %w", err)
		}
	}

	// Render finder patterns as unified 7x7 images
	if finderImg != "" {
		finderSize := 7 * moduleSize

		// Top-left finder
		px := float64(quietZone) * moduleSize
		py := float64(quietZone) * moduleSize
		fmt.Fprintf(&buf, `<image x="%.2f" y="%.2f" width="%.2f" height="%.2f" href="%s"/>`,
			px, py, finderSize, finderSize, finderImg)
		buf.WriteString("\n")

		// Top-right finder
		px = float64(quietZone+matrixSize-7) * moduleSize
		py = float64(quietZone) * moduleSize
		fmt.Fprintf(&buf, `<image x="%.2f" y="%.2f" width="%.2f" height="%.2f" href="%s" transform="scale(-1,1) translate(%.2f,0)"/>`,
			px, py, finderSize, finderSize, finderImg, -(2*px + finderSize))
		buf.WriteString("\n")

		// Bottom-left finder
		px = float64(quietZone) * moduleSize
		py = float64(quietZone+matrixSize-7) * moduleSize
		fmt.Fprintf(&buf, `<image x="%.2f" y="%.2f" width="%.2f" height="%.2f" href="%s" transform="scale(1,-1) translate(0,%.2f)"/>`,
			px, py, finderSize, finderSize, finderImg, -(2*py + finderSize))
		buf.WriteString("\n")
	}

	// Render alignment patterns as unified 5x5 images
	if alignImg != "" {
		alignSize := 5 * moduleSize
		alignPositions := getAlignmentPositions(matrixSize)

		for _, ay := range alignPositions {
			for _, ax := range alignPositions {
				// Skip if overlapping with finder patterns
				if isFinderArea(ax, ay, matrixSize) {
					continue
				}
				// Alignment pattern is centered, so offset by 2
				px := float64(quietZone+ax-2) * moduleSize
				py := float64(quietZone+ay-2) * moduleSize
				fmt.Fprintf(&buf, `<image x="%.2f" y="%.2f" width="%.2f" height="%.2f" href="%s"/>`,
					px, py, alignSize, alignSize, alignImg)
				buf.WriteString("\n")
			}
		}
	}

	// Render regular modules (skip finder and alignment areas if custom images provided)
	for y := 0; y < matrixSize; y++ {
		for x := 0; x < matrixSize; x++ {
			// Skip modules in the logo zone
			if hasLogoZone && x >= logoMinX && x <= logoMaxX && y >= logoMinY && y <= logoMaxY {
				continue
			}

			mod := r.matrix.Get(x, y)
			if !mod.Dark {
				continue
			}

			// Skip finder pattern modules if we rendered them as unified images
			if finderImg != "" && mod.Type == encoder.ModuleFinder {
				continue
			}

			// Skip alignment pattern modules if we rendered them as unified images
			if alignImg != "" && mod.Type == encoder.ModuleAlignment {
				continue
			}

			// Skip separator modules around finders when using custom finder images
			if finderImg != "" && mod.Type == encoder.ModuleFinderSeparator {
				continue
			}

			px := float64(quietZone+x) * moduleSize
			py := float64(quietZone+y) * moduleSize

			if moduleImg != "" {
				fmt.Fprintf(&buf, `<image x="%.2f" y="%.2f" width="%.2f" height="%.2f" href="%s"/>`,
					px, py, moduleSize, moduleSize, moduleImg)
				buf.WriteString("\n")
			}
		}
	}

	// Render logo if configured
	if r.hasLogo() {
		logoSVG, err := r.renderLogo()
		if err != nil {
			return nil, err
		}
		buf.WriteString(logoSVG)
	}

	buf.WriteString("</svg>")
	return buf.Bytes(), nil
}

// getAlignmentPositions returns center positions of alignment patterns
func getAlignmentPositions(matrixSize int) []int {
	version := (matrixSize - 17) / 4
	if version < 2 {
		return nil
	}

	count := version/7 + 2
	first := 6
	last := matrixSize - 7

	if count == 2 {
		return []int{first, last}
	}

	step := (last - first) / (count - 1)
	if step%2 != 0 {
		step++
	}

	positions := make([]int, count)
	positions[0] = first
	positions[count-1] = last
	for i := count - 2; i >= 1; i-- {
		positions[i] = positions[i+1] - step
	}

	return positions
}

// isFinderArea checks if position overlaps with finder patterns
func isFinderArea(x, y, size int) bool {
	// Top-left finder (0-6, 0-6)
	if x <= 8 && y <= 8 {
		return true
	}
	// Top-right finder (size-7 to size-1, 0-6)
	if x >= size-9 && y <= 8 {
		return true
	}
	// Bottom-left finder (0-6, size-7 to size-1)
	if x <= 8 && y >= size-9 {
		return true
	}
	return false
}

// getLogoDimensions reads image dimensions from file
// For SVG files, returns 1:1 aspect ratio
func getLogoDimensions(path string) (int, int, error) {
	if strings.HasSuffix(strings.ToLower(path), ".svg") {
		return 1, 1, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}

	return cfg.Width, cfg.Height, nil
}

// renderLogo renders the logo in the center of the QR code
func (r *Renderer) renderLogo() (string, error) {
	if !r.hasLogo() {
		return "", nil
	}

	logo := r.config.Logo

	// Get logo as data URI (from image.Image or file path)
	logoURI, err := r.getLogoDataURI()
	if err != nil {
		return "", fmt.Errorf("failed to load logo: %w", err)
	}

	// Get calculated dimensions
	logoWidth, logoHeight, padding, err := r.calculateLogoDimensions()
	if err != nil {
		return "", err
	}

	qrSize := float64(r.config.Size)

	// Center position
	logoX := (qrSize - logoWidth) / 2
	logoY := (qrSize - logoHeight) / 2

	var buf strings.Builder

	// Draw background rectangle behind logo
	bgColor := logo.Background
	if bgColor == "" {
		bgColor = "#FFFFFF"
	}
	if bgColor != "transparent" {
		bgX := logoX - padding
		bgY := logoY - padding
		bgWidth := logoWidth + 2*padding
		bgHeight := logoHeight + 2*padding
		cornerRadius := padding / 2
		fmt.Fprintf(&buf, `<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s" rx="%.2f"/>`,
			bgX, bgY, bgWidth, bgHeight, bgColor, cornerRadius)
		buf.WriteString("\n")
	}

	// Draw the logo image
	fmt.Fprintf(&buf, `<image x="%.2f" y="%.2f" width="%.2f" height="%.2f" href="%s"/>`,
		logoX, logoY, logoWidth, logoHeight, logoURI)
	buf.WriteString("\n")

	return buf.String(), nil
}

// getLogoDataURI returns the logo as a data URI, from either image.Image or file path
func (r *Renderer) getLogoDataURI() (string, error) {
	logo := r.config.Logo

	// If we have an in-memory image, encode it to PNG
	if logo.Image != nil {
		return encodeImageToDataURI(logo.Image)
	}

	// Otherwise load from file path
	return loadImageAsDataURI(logo.Path)
}

// encodeImageToDataURI encodes an image.Image to a PNG data URI
func encodeImageToDataURI(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/png;base64," + encoded, nil
}

// loadImageAsDataURI reads a PNG file and returns a data URI
func loadImageAsDataURI(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// Detect image type from extension
	mimeType := "image/png"
	if strings.HasSuffix(strings.ToLower(path), ".jpg") || strings.HasSuffix(strings.ToLower(path), ".jpeg") {
		mimeType = "image/jpeg"
	} else if strings.HasSuffix(strings.ToLower(path), ".svg") {
		mimeType = "image/svg+xml"
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

// transformPath scales and translates an SVG path from unit coordinates
func transformPath(path string, tx, ty, scale float64) string {
	// Parse path and transform coordinates
	// Handle both absolute (uppercase) and relative (lowercase) commands
	re := regexp.MustCompile(`([MLHVACQZmlhvacqz])([^MLHVACQZmlhvacqz]*)`)
	matches := re.FindAllStringSubmatch(path, -1)

	var result strings.Builder
	var curX, curY float64 // Track current position for relative commands

	for _, match := range matches {
		cmd := match[1]
		args := strings.TrimSpace(match[2])
		isRelative := cmd >= "a" && cmd <= "z"
		cmdUpper := strings.ToUpper(cmd)

		switch cmdUpper {
		case "M", "L":
			parts := splitNumbers(args)
			if len(parts) >= 2 {
				var x, y float64
				if isRelative {
					x = curX + parts[0]*scale
					y = curY + parts[1]*scale
				} else {
					x = parts[0]*scale + tx
					y = parts[1]*scale + ty
				}
				curX, curY = x, y
				fmt.Fprintf(&result, "%s%.2f %.2f", cmdUpper, x, y)
			}
		case "H":
			parts := splitNumbers(args)
			if len(parts) >= 1 {
				var x float64
				if isRelative {
					x = curX + parts[0]*scale
				} else {
					x = parts[0]*scale + tx
				}
				curX = x
				fmt.Fprintf(&result, "L%.2f %.2f", x, curY)
			}
		case "V":
			parts := splitNumbers(args)
			if len(parts) >= 1 {
				var y float64
				if isRelative {
					y = curY + parts[0]*scale
				} else {
					y = parts[0]*scale + ty
				}
				curY = y
				fmt.Fprintf(&result, "L%.2f %.2f", curX, y)
			}
		case "A":
			// Arc: rx ry rotation large-arc sweep x y
			parts := splitNumbers(args)
			if len(parts) >= 7 {
				rx := parts[0] * scale
				ry := parts[1] * scale
				rot := parts[2]
				large := parts[3]
				sweep := parts[4]
				var x, y float64
				if isRelative {
					x = curX + parts[5]*scale
					y = curY + parts[6]*scale
				} else {
					x = parts[5]*scale + tx
					y = parts[6]*scale + ty
				}
				curX, curY = x, y
				fmt.Fprintf(&result, "A%.2f %.2f %.0f %.0f %.0f %.2f %.2f", rx, ry, rot, large, sweep, x, y)
			}
		case "C":
			// Cubic bezier: x1 y1 x2 y2 x y
			parts := splitNumbers(args)
			if len(parts) >= 6 {
				var x1, y1, x2, y2, x, y float64
				if isRelative {
					x1 = curX + parts[0]*scale
					y1 = curY + parts[1]*scale
					x2 = curX + parts[2]*scale
					y2 = curY + parts[3]*scale
					x = curX + parts[4]*scale
					y = curY + parts[5]*scale
				} else {
					x1 = parts[0]*scale + tx
					y1 = parts[1]*scale + ty
					x2 = parts[2]*scale + tx
					y2 = parts[3]*scale + ty
					x = parts[4]*scale + tx
					y = parts[5]*scale + ty
				}
				curX, curY = x, y
				fmt.Fprintf(&result, "C%.2f %.2f %.2f %.2f %.2f %.2f", x1, y1, x2, y2, x, y)
			}
		case "Q":
			// Quadratic bezier: x1 y1 x y
			parts := splitNumbers(args)
			if len(parts) >= 4 {
				var x1, y1, x, y float64
				if isRelative {
					x1 = curX + parts[0]*scale
					y1 = curY + parts[1]*scale
					x = curX + parts[2]*scale
					y = curY + parts[3]*scale
				} else {
					x1 = parts[0]*scale + tx
					y1 = parts[1]*scale + ty
					x = parts[2]*scale + tx
					y = parts[3]*scale + ty
				}
				curX, curY = x, y
				fmt.Fprintf(&result, "Q%.2f %.2f %.2f %.2f", x1, y1, x, y)
			}
		case "Z":
			result.WriteString("Z")
		}
	}
	return result.String()
}

// splitNumbers parses space/comma separated numbers
func splitNumbers(s string) []float64 {
	s = strings.ReplaceAll(s, ",", " ")
	parts := strings.Fields(s)
	nums := make([]float64, 0, len(parts))
	for _, p := range parts {
		if n, err := strconv.ParseFloat(p, 64); err == nil {
			nums = append(nums, n)
		}
	}
	return nums
}

// RenderPNG generates a PNG by rasterizing the SVG.
func (r *Renderer) RenderPNG() ([]byte, error) {
	return nil, fmt.Errorf("PNG rendering not supported - use SVG output")
}

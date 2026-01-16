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
func (r *renderer) hasLogo() bool {
	logo := r.config.Logo
	return logo != nil && (logo.Image != nil || logo.Path != "")
}

// getLogoSourceDimensions returns the original dimensions of the logo source
func (r *renderer) getLogoSourceDimensions() (int, int, error) {
	logo := r.config.Logo
	if logo.Image != nil {
		bounds := logo.Image.Bounds()
		return bounds.Dx(), bounds.Dy(), nil
	}
	return getLogoDimensions(logo.Path)
}

// calculateLogoDimensions returns logo width, height, and padding in pixels
func (r *renderer) calculateLogoDimensions() (float64, float64, float64, error) {
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

// renderer converts a QR matrix to SVG output.
type renderer struct {
	config *Config
	matrix *encoder.Matrix
}

// newRenderer creates a renderer for the given matrix and config.
func newRenderer(matrix *encoder.Matrix, config *Config) *renderer {
	return &renderer{
		config: config,
		matrix: matrix,
	}
}

// renderSVG generates the SVG representation of the QR code.
func (r *renderer) renderSVG() ([]byte, error) {
	// Check if using custom images
	if r.config.Images != nil && (r.config.Images.Module != "" || r.config.Images.Finder != "" || r.config.Images.Alignment != "") {
		return r.renderWithImages()
	}
	return r.renderWithShapes()
}

// renderWithShapes renders QR code using vector shapes
func (r *renderer) renderWithShapes() ([]byte, error) {
	// Calculate logo exclusion zone
	logoMinX, logoMinY, logoMaxX, logoMaxY, hasLogoZone, err := r.calculateExclusionZone()
	if err != nil {
		return nil, err
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
	r.writeBackground(&buf)

	// Get module color/fill
	moduleFill := "#000000"
	if r.config.Modules.Color != nil {
		moduleFill = r.config.Modules.Color.SVGFill("module-fill")
	}

	// Get shape - using "square" as safe default if nil or unknown
	shapeName := r.config.Modules.Shape
	if shapeName == "" {
		shapeName = "square"
	}
	shape := shapes.Get(shapeName)
	if shape == nil {
		shape = shapes.Get("square")
	}

	// Draw all valid modules
	r.drawModulesShapes(&buf, shape, moduleFill, hasLogoZone, logoMinX, logoMinY, logoMaxX, logoMaxY)

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

func (r *renderer) drawModulesShapes(buf *bytes.Buffer, shape shapes.Shape, moduleFill string, hasLogoZone bool, logoMinX, logoMinY, logoMaxX, logoMaxY int) {
	matrixSize := r.matrix.Size()
	quietZone := r.config.QuietZone
	totalModules := matrixSize + 2*quietZone
	moduleSize := float64(r.config.Size) / float64(totalModules)

	// Render all dark modules as a single path for efficiency
	fmt.Fprintf(buf, `<path fill="%s" d="`, moduleFill)

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
}

// renderWithImages renders QR code using custom PNG images
func (r *renderer) renderWithImages() ([]byte, error) {
	matrixSize := r.matrix.Size()

	// Calculate logo exclusion zone
	logoMinX, logoMinY, logoMaxX, logoMaxY, hasLogoZone, err := r.calculateExclusionZone()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	r.writeSVGHeader(&buf)

	// Background
	r.writeBackground(&buf)

	// Load images
	moduleImg, finderImg, alignImg, err := r.loadCustomImages()
	if err != nil {
		return nil, err
	}

	// Render finder patterns
	if finderImg != "" {
		r.renderFinderImages(&buf, finderImg, matrixSize)
	}

	// Render alignment patterns
	if alignImg != "" {
		r.renderAlignmentImages(&buf, alignImg, finderImg != "", matrixSize)
	}

	// Render custom image modules
	// Skip modules in the logo zone or those covered by custom finders/alignments
	r.renderImageModules(&buf, moduleImg, finderImg, alignImg, matrixSize, hasLogoZone, logoMinX, logoMinY, logoMaxX, logoMaxY)

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

func (r *renderer) writeSVGHeader(buf *bytes.Buffer) {
	// SVG header with xlink namespace for images
	fmt.Fprintf(buf, `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 %d %d" width="%d" height="%d">`,
		r.config.Size, r.config.Size, r.config.Size, r.config.Size)
	buf.WriteString("\n")
}

func (r *renderer) writeBackground(buf *bytes.Buffer) {
	bgColor := "#FFFFFF"
	if r.config.Background != nil {
		bgColor = r.config.Background.SVGFill("")
	}
	fmt.Fprintf(buf, `<rect width="100%%" height="100%%" fill="%s"/>`, bgColor)
	buf.WriteString("\n")
}

func (r *renderer) calculateExclusionZone() (minX, minY, maxX, maxY int, active bool, err error) {
	if !r.hasLogo() {
		return 0, 0, 0, 0, false, nil
	}

	logoWidth, logoHeight, padding, err := r.calculateLogoDimensions()
	if err != nil {
		return 0, 0, 0, 0, false, fmt.Errorf("failed to calculate logo dimensions: %w", err)
	}

	matrixSize := r.matrix.Size()
	quietZone := r.config.QuietZone
	totalModules := matrixSize + 2*quietZone
	moduleSize := float64(r.config.Size) / float64(totalModules)

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

	minX = matrixCenterX - excludeHalfX
	minY = matrixCenterY - excludeHalfY
	maxX = matrixCenterX + excludeHalfX
	maxY = matrixCenterY + excludeHalfY

	return minX, minY, maxX, maxY, true, nil
}

func (r *renderer) loadCustomImages() (modImg, findImg, alignImg string, err error) {
	if r.config.Images.Module != "" {
		modImg, err = loadImageAsDataURI(r.config.Images.Module)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to load module image: %w", err)
		}
	}
	if r.config.Images.Finder != "" {
		findImg, err = loadImageAsDataURI(r.config.Images.Finder)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to load finder image: %w", err)
		}
	}
	if r.config.Images.Alignment != "" {
		alignImg, err = loadImageAsDataURI(r.config.Images.Alignment)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to load alignment image: %w", err)
		}
	}
	return
}

func (r *renderer) renderFinderImages(buf *bytes.Buffer, finderImg string, matrixSize int) {
	quietZone := r.config.QuietZone
	totalModules := matrixSize + 2*quietZone
	moduleSize := float64(r.config.Size) / float64(totalModules)
	finderSize := 7 * moduleSize

	// Top-left finder
	px := float64(quietZone) * moduleSize
	py := float64(quietZone) * moduleSize
	fmt.Fprintf(buf, `<image x="%.2f" y="%.2f" width="%.2f" height="%.2f" href="%s"/>`,
		px, py, finderSize, finderSize, finderImg)
	buf.WriteString("\n")

	// Top-right finder
	px = float64(quietZone+matrixSize-7) * moduleSize
	py = float64(quietZone) * moduleSize
	fmt.Fprintf(buf, `<image x="%.2f" y="%.2f" width="%.2f" height="%.2f" href="%s" transform="scale(-1,1) translate(%.2f,0)"/>`,
		px, py, finderSize, finderSize, finderImg, -(2*px + finderSize))
	buf.WriteString("\n")

	// Bottom-left finder
	px = float64(quietZone) * moduleSize
	py = float64(quietZone+matrixSize-7) * moduleSize
	fmt.Fprintf(buf, `<image x="%.2f" y="%.2f" width="%.2f" height="%.2f" href="%s" transform="scale(1,-1) translate(0,%.2f)"/>`,
		px, py, finderSize, finderSize, finderImg, -(2*py + finderSize))
	buf.WriteString("\n")
}

func (r *renderer) renderAlignmentImages(buf *bytes.Buffer, alignImg string, hasFinderImg bool, matrixSize int) {
	quietZone := r.config.QuietZone
	totalModules := matrixSize + 2*quietZone
	moduleSize := float64(r.config.Size) / float64(totalModules)
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
			fmt.Fprintf(buf, `<image x="%.2f" y="%.2f" width="%.2f" height="%.2f" href="%s"/>`,
				px, py, alignSize, alignSize, alignImg)
			buf.WriteString("\n")
		}
	}
}

func (r *renderer) renderImageModules(buf *bytes.Buffer, moduleImg, finderImg, alignImg string, matrixSize int, hasLogoZone bool, logoMinX, logoMinY, logoMaxX, logoMaxY int) {
	if moduleImg == "" {
		return
	}

	quietZone := r.config.QuietZone
	totalModules := matrixSize + 2*quietZone
	moduleSize := float64(r.config.Size) / float64(totalModules)

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

			fmt.Fprintf(buf, `<image x="%.2f" y="%.2f" width="%.2f" height="%.2f" href="%s"/>`,
				px, py, moduleSize, moduleSize, moduleImg)
			buf.WriteString("\n")
		}
	}
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
func (r *renderer) renderLogo() (string, error) {
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
func (r *renderer) getLogoDataURI() (string, error) {
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

// pathCommandRe matches SVG path commands and their arguments
var pathCommandRe = regexp.MustCompile(`([MLHVACQZmlhvacqz])([^MLHVACQZmlhvacqz]*)`)

// transformPath scales and translates an SVG path from unit coordinates
func transformPath(path string, tx, ty, scale float64) string {
	matches := pathCommandRe.FindAllStringSubmatch(path, -1)

	var result strings.Builder
	var curX, curY float64 // Track current position for relative commands

	for _, match := range matches {
		cmd := match[1]
		args := strings.TrimSpace(match[2])

		segment, newX, newY := processPathCommand(cmd, args, curX, curY, tx, ty, scale)
		result.WriteString(segment)
		curX, curY = newX, newY
	}

	return result.String()
}

// processPathCommand handles a single SVG path command transformation
func processPathCommand(cmd, args string, curX, curY, tx, ty, scale float64) (string, float64, float64) {
	isRelative := cmd >= "a" && cmd <= "z"
	cmdUpper := strings.ToUpper(cmd)

	switch cmdUpper {
	case "M", "L":
		return handleMoveLine(cmdUpper, args, curX, curY, tx, ty, scale, isRelative)
	case "H":
		return handleHorizontal(args, curX, curY, tx, ty, scale, isRelative)
	case "V":
		return handleVertical(args, curX, curY, tx, ty, scale, isRelative)
	case "A":
		return handleArc(args, curX, curY, tx, ty, scale, isRelative)
	case "C":
		return handleCubic(args, curX, curY, tx, ty, scale, isRelative)
	case "Q":
		return handleQuad(args, curX, curY, tx, ty, scale, isRelative)
	case "Z":
		return "Z ", curX, curY
	}

	return "", curX, curY
}

func handleMoveLine(cmd, args string, curX, curY, tx, ty, scale float64, isRelative bool) (string, float64, float64) {
	parts := splitNumbers(args)
	if len(parts) < 2 {
		return "", curX, curY
	}
	var x, y float64
	if isRelative {
		x = curX + parts[0]*scale
		y = curY + parts[1]*scale
	} else {
		x = parts[0]*scale + tx
		y = parts[1]*scale + ty
	}
	return fmt.Sprintf("%s%.2f %.2f ", cmd, x, y), x, y
}

func handleHorizontal(args string, curX, curY, tx, ty, scale float64, isRelative bool) (string, float64, float64) {
	parts := splitNumbers(args)
	if len(parts) < 1 {
		return "", curX, curY
	}
	var x float64
	if isRelative {
		x = curX + parts[0]*scale
	} else {
		x = parts[0]*scale + tx
	}
	return fmt.Sprintf("L%.2f %.2f ", x, curY), x, curY
}

func handleVertical(args string, curX, curY, tx, ty, scale float64, isRelative bool) (string, float64, float64) {
	parts := splitNumbers(args)
	if len(parts) < 1 {
		return "", curX, curY
	}
	var y float64
	if isRelative {
		y = curY + parts[0]*scale
	} else {
		y = parts[0]*scale + ty
	}
	return fmt.Sprintf("L%.2f %.2f ", curX, y), curX, y
}

func handleArc(args string, curX, curY, tx, ty, scale float64, isRelative bool) (string, float64, float64) {
	parts := splitNumbers(args)
	if len(parts) < 7 {
		return "", curX, curY
	}
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
	return fmt.Sprintf("A%.2f %.2f %.0f %.0f %.0f %.2f %.2f ", rx, ry, rot, large, sweep, x, y), x, y
}

func handleCubic(args string, curX, curY, tx, ty, scale float64, isRelative bool) (string, float64, float64) {
	parts := splitNumbers(args)
	if len(parts) < 6 {
		return "", curX, curY
	}
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
	return fmt.Sprintf("C%.2f %.2f %.2f %.2f %.2f %.2f ", x1, y1, x2, y2, x, y), x, y
}

func handleQuad(args string, curX, curY, tx, ty, scale float64, isRelative bool) (string, float64, float64) {
	parts := splitNumbers(args)
	if len(parts) < 4 {
		return "", curX, curY
	}
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
	return fmt.Sprintf("Q%.2f %.2f %.2f %.2f ", x1, y1, x, y), x, y
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
func (r *renderer) RenderPNG() ([]byte, error) {
	return nil, fmt.Errorf("PNG rendering not supported - use SVG output")
}

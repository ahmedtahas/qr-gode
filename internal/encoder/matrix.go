package encoder

// Module represents a single module (cell) in the QR code.
type Module struct {
	Dark     bool       // true = dark module
	Type     ModuleType // What kind of module this is
	Reserved bool       // true = function pattern, can't be masked
}

// ModuleType identifies the type of module for styling purposes.
type ModuleType int

const (
	ModuleData ModuleType = iota
	ModuleFinder
	ModuleFinderSeparator
	ModuleAlignment
	ModuleTiming
	ModuleFormatInfo
	ModuleVersionInfo
	ModuleDarkModule // The single always-dark module
)

// Matrix represents the QR code module grid.
type Matrix struct {
	size    int
	modules [][]Module
}

// NewMatrix creates a matrix for the given version.
func NewMatrix(version Version) *Matrix {
	size := version.Size()
	modules := make([][]Module, size)
	for i := range modules {
		modules[i] = make([]Module, size)
	}
	return &Matrix{
		size:    size,
		modules: modules,
	}
}

// Size returns the dimension of the matrix.
func (m *Matrix) Size() int {
	return m.size
}

// Get returns the module at (x, y).
func (m *Matrix) Get(x, y int) Module {
	return m.modules[y][x]
}

// Set sets the module at (x, y).
func (m *Matrix) Set(x, y int, mod Module) {
	m.modules[y][x] = mod
}

// PlaceFunctionPatterns places all non-data patterns on the matrix.
func (m *Matrix) PlaceFunctionPatterns(version Version) {
	// 1. Place finder patterns (3 corners)
	m.placeFinder(0, 0)        // top-left
	m.placeFinder(m.size-7, 0) // top-right
	m.placeFinder(0, m.size-7) // bottom-left

	// 2. Place finder separators
	m.placeSeparators()

	// 3. Place timing patterns
	m.placeTiming()

	// 4. Place alignment patterns (version 2+)
	m.placeAlignment(version)

	// 5. Place dark module (always at this position)
	m.Set(8, 4*int(version)+9, Module{Dark: true, Type: ModuleDarkModule, Reserved: true})

	// 6. Reserve format info areas
	m.reserveFormatInfo()

	// 7. Reserve version info areas (version 7+)
	if version >= 7 {
		m.reserveVersionInfo()
	}
}

// PlaceData places the data codewords onto the matrix.
func (m *Matrix) PlaceData(data []byte) {
	// Convert bytes to bits
	bits := make([]bool, len(data)*8)
	for i, b := range data {
		for j := 0; j < 8; j++ {
			bits[i*8+j] = (b>>(7-j))&1 == 1
		}
	}

	bitIndex := 0
	// Start from right side, move left in 2-column strips
	// Skip column 6 (timing pattern)
	for col := m.size - 1; col >= 0; col -= 2 {
		if col == 6 {
			col = 5 // Skip timing column
		}

		// Zigzag: go up on even strips, down on odd
		// We determine direction based on which strip we're on
		upward := ((m.size-1-col)/2)%2 == 0

		for row := 0; row < m.size; row++ {
			actualRow := row
			if upward {
				actualRow = m.size - 1 - row
			}

			// Try right column, then left column
			for dx := 0; dx <= 1; dx++ {
				x := col - dx
				if x < 0 {
					continue
				}

				if !m.modules[actualRow][x].Reserved {
					dark := false
					if bitIndex < len(bits) {
						dark = bits[bitIndex]
						bitIndex++
					}
					m.modules[actualRow][x].Dark = dark
					m.modules[actualRow][x].Type = ModuleData
				}
			}
		}
	}
}

// Clone creates a deep copy of the matrix.
func (m *Matrix) Clone() *Matrix {
	clone := &Matrix{
		size:    m.size,
		modules: make([][]Module, m.size),
	}
	for i := range m.modules {
		clone.modules[i] = make([]Module, m.size)
		copy(clone.modules[i], m.modules[i])
	}
	return clone
}

// placeTiming places the horizontal and vertical timing patterns.
// These are alternating dark/light modules between finder patterns.
func (m *Matrix) placeTiming() {
	// Horizontal timing: row 6, from col 8 to size-9
	for i := 8; i < m.size-8; i++ {
		dark := i%2 == 0
		m.Set(i, 6, Module{Dark: dark, Type: ModuleTiming, Reserved: true})
	}

	// Vertical timing: col 6, from row 8 to size-9
	for i := 8; i < m.size-8; i++ {
		dark := i%2 == 0
		m.Set(6, i, Module{Dark: dark, Type: ModuleTiming, Reserved: true})
	}
}

// placeAlignment places alignment patterns for version 2+.
func (m *Matrix) placeAlignment(version Version) {
	positions := alignmentPatternPositions(version)
	if len(positions) == 0 {
		return
	}

	for _, row := range positions {
		for _, col := range positions {
			// Skip if overlapping with finder patterns
			if m.Get(col, row).Type == ModuleFinder || m.Get(col, row).Type == ModuleFinderSeparator {
				continue
			}
			m.placeAlignmentPattern(col, row)
		}
	}
}

// placeAlignmentPattern places a 5x5 alignment pattern centered at (cx, cy).
func (m *Matrix) placeAlignmentPattern(cx, cy int) {
	for row := -2; row <= 2; row++ {
		for col := -2; col <= 2; col++ {
			// Dark if on border OR center
			onBorder := row == -2 || row == 2 || col == -2 || col == 2
			isCenter := row == 0 && col == 0
			dark := onBorder || isCenter

			m.Set(cx+col, cy+row, Module{Dark: dark, Type: ModuleAlignment, Reserved: true})
		}
	}
}

// reserveFormatInfo reserves the format information areas (filled in later).
func (m *Matrix) reserveFormatInfo() {
	// Around top-left finder
	for i := 0; i < 9; i++ {
		if m.Get(i, 8).Type == 0 { // not already set
			m.Set(i, 8, Module{Type: ModuleFormatInfo, Reserved: true})
		}
		if m.Get(8, i).Type == 0 {
			m.Set(8, i, Module{Type: ModuleFormatInfo, Reserved: true})
		}
	}

	// Below top-right finder
	for i := 0; i < 8; i++ {
		m.Set(m.size-1-i, 8, Module{Type: ModuleFormatInfo, Reserved: true})
	}

	// Right of bottom-left finder
	for i := 0; i < 7; i++ {
		m.Set(8, m.size-1-i, Module{Type: ModuleFormatInfo, Reserved: true})
	}
}

// reserveVersionInfo reserves version information areas (version 7+).
func (m *Matrix) reserveVersionInfo() {
	// Bottom-left of top-right finder (6x3 block)
	for i := 0; i < 6; i++ {
		for j := 0; j < 3; j++ {
			m.Set(m.size-11+j, i, Module{Type: ModuleVersionInfo, Reserved: true})
		}
	}

	// Top-right of bottom-left finder (3x6 block)
	for i := 0; i < 6; i++ {
		for j := 0; j < 3; j++ {
			m.Set(i, m.size-11+j, Module{Type: ModuleVersionInfo, Reserved: true})
		}
	}
}

// alignmentPatternPositions calculates center positions for alignment patterns.
func alignmentPatternPositions(version Version) []int {
	if version == 1 {
		return nil
	}

	// Number of alignment coordinates = floor(version/7) + 2
	count := int(version)/7 + 2

	// First is always 6, last is always size-7
	first := 6
	last := int(version)*4 + 10 // same as size-7

	if count == 2 {
		return []int{first, last}
	}

	// Calculate step: must be even, spread evenly between first and last
	totalDistance := last - first
	steps := count - 1
	step := totalDistance / steps

	// Step must be even (QR spec requirement)
	if step%2 != 0 {
		step++
	}

	// Build positions from last going backwards
	positions := make([]int, count)
	positions[0] = first
	positions[count-1] = last
	for i := count - 2; i >= 1; i-- {
		positions[i] = positions[i+1] - step
	}

	return positions
}

// placeSeparators places white separators around the 3 finder patterns.
func (m *Matrix) placeSeparators() {
	// Top-left: right edge (col 7) and bottom edge (row 7)
	for i := 0; i < 8; i++ {
		m.Set(7, i, Module{Type: ModuleFinderSeparator, Reserved: true}) // right edge
		m.Set(i, 7, Module{Type: ModuleFinderSeparator, Reserved: true}) // bottom edge
	}

	// Top-right: left edge (col size-8) and bottom edge (row 7)
	for i := 0; i < 8; i++ {
		m.Set(m.size-8, i, Module{Type: ModuleFinderSeparator, Reserved: true})   // left edge
		m.Set(m.size-8+i, 7, Module{Type: ModuleFinderSeparator, Reserved: true}) // bottom edge
	}

	// Bottom-left: right edge (col 7) and top edge (row size-8)
	for i := 0; i < 8; i++ {
		m.Set(7, m.size-8+i, Module{Type: ModuleFinderSeparator, Reserved: true}) // right edge
		m.Set(i, m.size-8, Module{Type: ModuleFinderSeparator, Reserved: true})   // top edge
	}
}

// placeFinder places a 7x7 finder pattern at (x, y).
func (m *Matrix) placeFinder(x, y int) {
	for row := range 7 {
		for col := range 7 {
			// Dark if on border OR in 3x3 center
			onBorder := row == 0 || row == 6 || col == 0 || col == 6
			inCenter := row >= 2 && row <= 4 && col >= 2 && col <= 4
			dark := onBorder || inCenter

			m.Set(x+col, y+row, Module{
				Dark:     dark,
				Type:     ModuleFinder,
				Reserved: true,
			})
		}
	}
}

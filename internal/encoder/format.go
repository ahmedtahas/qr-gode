package encoder

// Format info mask pattern for XOR
const formatInfoMask = 0x5412 // 101010000010010

// ECL bits for format string (L=01, M=00, Q=11, H=10)
var eclBits = [4]uint16{0b01, 0b00, 0b11, 0b10}

// FormatInfo encodes error correction level and mask pattern.
// 15 bits total: 5 data bits + 10 ECC bits (BCH code).
func FormatInfo(ecl ErrorCorrectionLevel, mask MaskPattern) uint16 {
	// 1. Combine ECL (2 bits) + mask (3 bits) = 5 data bits
	data := (eclBits[ecl] << 3) | uint16(mask)

	// 2. Generate BCH(15,5) error correction
	// Generator polynomial: x^10 + x^8 + x^5 + x^4 + x^2 + x + 1 = 0x537
	info := data << 10 // Make room for 10 ECC bits

	// Polynomial division to get remainder
	for i := 4; i >= 0; i-- {
		if info&(1<<(i+10)) != 0 {
			info ^= 0x537 << i
		}
	}

	// Combine data and ECC
	info = (data << 10) | info

	// 3. XOR with mask pattern
	return info ^ formatInfoMask
}

// PlaceFormatInfo places the format info bits on the matrix.
// Format info appears in two locations for redundancy.
func PlaceFormatInfo(matrix *Matrix, info uint16) {
	size := matrix.Size()

	// Format info bit positions around top-left finder (bits 14-0, MSB first)
	// Vertical strip (column 8): rows 0-5, skip 6, rows 7-8
	// Horizontal strip (row 8): cols 0-5, skip 6, cols 7-8

	// Location 1: Around top-left finder
	// Bits 0-5 go in column 8, rows 0-5
	for i := 0; i <= 5; i++ {
		bit := (info >> (14 - i)) & 1
		matrix.Set(8, i, Module{Dark: bit == 1, Type: ModuleFormatInfo, Reserved: true})
	}
	// Bit 6 goes in column 8, row 7 (skip row 6 - timing)
	bit6 := (info >> 8) & 1
	matrix.Set(8, 7, Module{Dark: bit6 == 1, Type: ModuleFormatInfo, Reserved: true})
	// Bit 7 goes in column 8, row 8
	bit7 := (info >> 7) & 1
	matrix.Set(8, 8, Module{Dark: bit7 == 1, Type: ModuleFormatInfo, Reserved: true})
	// Bit 8 goes in row 8, column 7
	bit8 := (info >> 6) & 1
	matrix.Set(7, 8, Module{Dark: bit8 == 1, Type: ModuleFormatInfo, Reserved: true})
	// Bits 9-14 go in row 8, columns 5-0
	for i := 9; i <= 14; i++ {
		bit := (info >> (14 - i)) & 1
		matrix.Set(14-i, 8, Module{Dark: bit == 1, Type: ModuleFormatInfo, Reserved: true})
	}

	// Location 2: Split between top-right and bottom-left
	// Top-right: row 8, columns size-1 to size-8 (bits 0-7)
	for i := 0; i <= 7; i++ {
		bit := (info >> i) & 1
		matrix.Set(size-1-i, 8, Module{Dark: bit == 1, Type: ModuleFormatInfo, Reserved: true})
	}
	// Bottom-left: column 8, rows size-7 to size-1 (bits 8-14)
	for i := 0; i <= 6; i++ {
		bit := (info >> (8 + i)) & 1
		matrix.Set(8, size-7+i, Module{Dark: bit == 1, Type: ModuleFormatInfo, Reserved: true})
	}
}

// Pre-computed version info for versions 7-40 (already includes BCH ECC)
var versionInfoTable = []uint32{
	0x07C94, 0x085BC, 0x09A99, 0x0A4D3, 0x0BBF6, 0x0C762, 0x0D847, 0x0E60D,
	0x0F928, 0x10B78, 0x1145D, 0x12A17, 0x13532, 0x149A6, 0x15683, 0x168C9,
	0x177EC, 0x18EC4, 0x191E1, 0x1AFAB, 0x1B08E, 0x1CC1A, 0x1D33F, 0x1ED75,
	0x1F250, 0x209D5, 0x216F0, 0x228BA, 0x2379F, 0x24B0B, 0x2542E, 0x26A64,
	0x27541, 0x28C69,
}

// VersionInfo encodes the version number for versions 7-40.
// 18 bits total: 6 data bits + 12 ECC bits (BCH code).
func VersionInfo(version Version) uint32 {
	if version < 7 || version > 40 {
		return 0
	}
	return versionInfoTable[version-7]
}

// PlaceVersionInfo places version info on the matrix.
// Only for version 7+. Appears in two 6x3 blocks.
func PlaceVersionInfo(matrix *Matrix, version Version, info uint32) {
	if version < 7 {
		return
	}

	size := matrix.Size()

	// 18 bits: 6 rows x 3 columns
	// Location 1: Bottom-left of top-right finder (columns size-11 to size-9, rows 0-5)
	// Location 2: Top-right of bottom-left finder (rows size-11 to size-9, columns 0-5)
	for i := 0; i < 6; i++ {
		for j := 0; j < 3; j++ {
			bitIndex := i*3 + j
			bit := (info >> bitIndex) & 1

			// Location 1: near top-right finder
			matrix.Set(size-11+j, i, Module{Dark: bit == 1, Type: ModuleVersionInfo, Reserved: true})
			// Location 2: near bottom-left finder
			matrix.Set(i, size-11+j, Module{Dark: bit == 1, Type: ModuleVersionInfo, Reserved: true})
		}
	}
}

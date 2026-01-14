package encoder

// Encoder orchestrates the QR code encoding process.
type Encoder struct {
	data            string
	errorCorrection ErrorCorrectionLevel
	version         int
	mode            Mode
}

// ErrorCorrectionLevel defines redundancy level.
type ErrorCorrectionLevel int

const (
	LevelL ErrorCorrectionLevel = iota
	LevelM
	LevelQ
	LevelH
)

// New creates a new encoder for the given data.
func New(data string, ecl ErrorCorrectionLevel) *Encoder {
	return &Encoder{
		data:            data,
		errorCorrection: ecl,
	}
}

// Encode performs the full encoding process and returns the module matrix.
func (e *Encoder) Encode() (*Matrix, error) {
	// 1. Analyze data to determine best mode
	mode := AnalyzeData(e.data)
	e.mode = mode

	// 2. Determine minimum version that fits data + error correction
	version, err := DetermineVersion(len(e.data), mode, e.errorCorrection)
	if err != nil {
		return nil, err
	}
	e.version = int(version)

	// 3. Encode data to bit stream
	eccInfo := GetECCInfo(version, e.errorCorrection)
	dataCapacity := eccInfo.DataCapacity()
	bs := EncodeDataWithPadding(e.data, mode, version, dataCapacity)
	dataBytes := bs.Bytes()

	// 4. Generate error correction codewords
	dataBlocks, eccBlocks := GenerateECC(dataBytes, eccInfo)

	// 5. Structure data (interleave blocks)
	finalData := InterleaveBlocks(dataBlocks, eccBlocks)

	// 6. Create matrix and place function patterns
	matrix := NewMatrix(version)
	matrix.PlaceFunctionPatterns(version)

	// 7. Place data modules
	matrix.PlaceData(finalData)

	// 8. Apply masking and select best mask
	bestMask := SelectBestMask(matrix)
	ApplyMask(matrix, bestMask)

	// 9. Add format and version information
	formatInfo := FormatInfo(e.errorCorrection, bestMask)
	PlaceFormatInfo(matrix, formatInfo)

	if version >= 7 {
		versionInfo := VersionInfo(version)
		PlaceVersionInfo(matrix, version, versionInfo)
	}

	return matrix, nil
}

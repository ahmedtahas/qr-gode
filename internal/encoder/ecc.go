package encoder

// BlockInfo describes how data is split into blocks for error correction.
type BlockInfo struct {
	Count         int // Number of blocks with this configuration
	TotalCodewords int // Total codewords in each block
	DataCodewords  int // Data codewords per block (rest are ECC)
}

// ECCInfo holds error correction configuration for a version/level combo.
type ECCInfo struct {
	TotalCodewords  int         // Total codewords for this version
	ECCPerBlock     int         // Error correction codewords per block
	Group1          BlockInfo   // First group of blocks
	Group2          BlockInfo   // Second group (Count=0 if not used)
}

// DataCapacity returns the total data codewords (excluding ECC).
func (e ECCInfo) DataCapacity() int {
	return e.Group1.Count*e.Group1.DataCodewords + e.Group2.Count*e.Group2.DataCodewords
}

// eccTable holds ECC info for all version/level combinations.
// Index: [version-1][ecl] where ecl is L=0, M=1, Q=2, H=3
var eccTable = [40][4]ECCInfo{
	// Version 1
	{
		{26, 7, BlockInfo{1, 26, 19}, BlockInfo{}},  // L
		{26, 10, BlockInfo{1, 26, 16}, BlockInfo{}}, // M
		{26, 13, BlockInfo{1, 26, 13}, BlockInfo{}}, // Q
		{26, 17, BlockInfo{1, 26, 9}, BlockInfo{}},  // H
	},
	// Version 2
	{
		{44, 10, BlockInfo{1, 44, 34}, BlockInfo{}}, // L
		{44, 16, BlockInfo{1, 44, 28}, BlockInfo{}}, // M
		{44, 22, BlockInfo{1, 44, 22}, BlockInfo{}}, // Q
		{44, 28, BlockInfo{1, 44, 16}, BlockInfo{}}, // H
	},
	// Version 3
	{
		{70, 15, BlockInfo{1, 70, 55}, BlockInfo{}},  // L
		{70, 26, BlockInfo{1, 70, 44}, BlockInfo{}},  // M
		{70, 18, BlockInfo{2, 35, 17}, BlockInfo{}},  // Q
		{70, 22, BlockInfo{2, 35, 13}, BlockInfo{}},  // H
	},
	// Version 4
	{
		{100, 20, BlockInfo{1, 100, 80}, BlockInfo{}}, // L
		{100, 18, BlockInfo{2, 50, 32}, BlockInfo{}},  // M
		{100, 26, BlockInfo{2, 50, 24}, BlockInfo{}},  // Q
		{100, 16, BlockInfo{4, 25, 9}, BlockInfo{}},   // H
	},
	// Version 5
	{
		{134, 26, BlockInfo{1, 134, 108}, BlockInfo{}},        // L
		{134, 24, BlockInfo{2, 67, 43}, BlockInfo{}},          // M
		{134, 18, BlockInfo{2, 33, 15}, BlockInfo{2, 34, 16}}, // Q
		{134, 22, BlockInfo{2, 33, 11}, BlockInfo{2, 34, 12}}, // H
	},
	// Version 6
	{
		{172, 18, BlockInfo{2, 86, 68}, BlockInfo{}},          // L
		{172, 16, BlockInfo{4, 43, 27}, BlockInfo{}},          // M
		{172, 24, BlockInfo{4, 43, 19}, BlockInfo{}},          // Q
		{172, 28, BlockInfo{4, 43, 15}, BlockInfo{}},          // H
	},
	// Version 7
	{
		{196, 20, BlockInfo{2, 98, 78}, BlockInfo{}},          // L
		{196, 18, BlockInfo{4, 49, 31}, BlockInfo{}},          // M
		{196, 18, BlockInfo{2, 32, 14}, BlockInfo{4, 33, 15}}, // Q
		{196, 26, BlockInfo{4, 39, 13}, BlockInfo{1, 40, 14}}, // H
	},
	// Version 8
	{
		{242, 24, BlockInfo{2, 121, 97}, BlockInfo{}},          // L
		{242, 22, BlockInfo{2, 60, 38}, BlockInfo{2, 61, 39}},  // M
		{242, 22, BlockInfo{4, 40, 18}, BlockInfo{2, 41, 19}},  // Q
		{242, 26, BlockInfo{4, 40, 14}, BlockInfo{2, 41, 15}},  // H
	},
	// Version 9
	{
		{292, 30, BlockInfo{2, 146, 116}, BlockInfo{}},         // L
		{292, 22, BlockInfo{3, 58, 36}, BlockInfo{2, 59, 37}},  // M
		{292, 20, BlockInfo{4, 36, 16}, BlockInfo{4, 37, 17}},  // Q
		{292, 24, BlockInfo{4, 36, 12}, BlockInfo{4, 37, 13}},  // H
	},
	// Version 10
	{
		{346, 18, BlockInfo{2, 86, 68}, BlockInfo{2, 87, 69}},  // L
		{346, 26, BlockInfo{4, 69, 43}, BlockInfo{1, 70, 44}},  // M
		{346, 24, BlockInfo{6, 43, 19}, BlockInfo{2, 44, 20}},  // Q
		{346, 28, BlockInfo{6, 43, 15}, BlockInfo{2, 44, 16}},  // H
	},
	// Version 11
	{
		{404, 20, BlockInfo{4, 101, 81}, BlockInfo{}},          // L
		{404, 30, BlockInfo{1, 80, 50}, BlockInfo{4, 81, 51}},  // M
		{404, 28, BlockInfo{4, 50, 22}, BlockInfo{4, 51, 23}},  // Q
		{404, 24, BlockInfo{3, 36, 12}, BlockInfo{8, 37, 13}},  // H
	},
	// Version 12
	{
		{466, 24, BlockInfo{2, 116, 92}, BlockInfo{2, 117, 93}}, // L
		{466, 22, BlockInfo{6, 58, 36}, BlockInfo{2, 59, 37}},   // M
		{466, 26, BlockInfo{4, 46, 20}, BlockInfo{6, 47, 21}},   // Q
		{466, 28, BlockInfo{7, 42, 14}, BlockInfo{4, 43, 15}},   // H
	},
	// Version 13
	{
		{532, 26, BlockInfo{4, 133, 107}, BlockInfo{}},          // L
		{532, 22, BlockInfo{8, 59, 37}, BlockInfo{1, 60, 38}},   // M
		{532, 24, BlockInfo{8, 44, 20}, BlockInfo{4, 45, 21}},   // Q
		{532, 22, BlockInfo{12, 33, 11}, BlockInfo{4, 34, 12}},  // H
	},
	// Version 14
	{
		{581, 30, BlockInfo{3, 145, 115}, BlockInfo{1, 146, 116}}, // L
		{581, 24, BlockInfo{4, 64, 40}, BlockInfo{5, 65, 41}},     // M
		{581, 20, BlockInfo{11, 36, 16}, BlockInfo{5, 37, 17}},    // Q
		{581, 24, BlockInfo{11, 36, 12}, BlockInfo{5, 37, 13}},    // H
	},
	// Version 15
	{
		{655, 22, BlockInfo{5, 109, 87}, BlockInfo{1, 110, 88}},  // L
		{655, 24, BlockInfo{5, 65, 41}, BlockInfo{5, 66, 42}},    // M
		{655, 30, BlockInfo{5, 54, 24}, BlockInfo{7, 55, 25}},    // Q
		{655, 24, BlockInfo{11, 36, 12}, BlockInfo{7, 37, 13}},   // H
	},
	// Version 16
	{
		{733, 24, BlockInfo{5, 122, 98}, BlockInfo{1, 123, 99}},  // L
		{733, 28, BlockInfo{7, 73, 45}, BlockInfo{3, 74, 46}},    // M
		{733, 24, BlockInfo{15, 43, 19}, BlockInfo{2, 44, 20}},   // Q
		{733, 30, BlockInfo{3, 45, 15}, BlockInfo{13, 46, 16}},   // H
	},
	// Version 17
	{
		{815, 28, BlockInfo{1, 135, 107}, BlockInfo{5, 136, 108}}, // L
		{815, 28, BlockInfo{10, 74, 46}, BlockInfo{1, 75, 47}},    // M
		{815, 28, BlockInfo{1, 50, 22}, BlockInfo{15, 51, 23}},    // Q
		{815, 28, BlockInfo{2, 42, 14}, BlockInfo{17, 43, 15}},    // H
	},
	// Version 18
	{
		{901, 30, BlockInfo{5, 150, 120}, BlockInfo{1, 151, 121}}, // L
		{901, 26, BlockInfo{9, 69, 43}, BlockInfo{4, 70, 44}},     // M
		{901, 28, BlockInfo{17, 50, 22}, BlockInfo{1, 51, 23}},    // Q
		{901, 28, BlockInfo{2, 42, 14}, BlockInfo{19, 43, 15}},    // H
	},
	// Version 19
	{
		{991, 28, BlockInfo{3, 141, 113}, BlockInfo{4, 142, 114}}, // L
		{991, 26, BlockInfo{3, 70, 44}, BlockInfo{11, 71, 45}},    // M
		{991, 26, BlockInfo{17, 47, 21}, BlockInfo{4, 48, 22}},    // Q
		{991, 26, BlockInfo{9, 39, 13}, BlockInfo{16, 40, 14}},    // H
	},
	// Version 20
	{
		{1085, 28, BlockInfo{3, 135, 107}, BlockInfo{5, 136, 108}}, // L
		{1085, 26, BlockInfo{3, 67, 41}, BlockInfo{13, 68, 42}},    // M
		{1085, 30, BlockInfo{15, 54, 24}, BlockInfo{5, 55, 25}},    // Q
		{1085, 28, BlockInfo{15, 43, 15}, BlockInfo{10, 44, 16}},   // H
	},
	// Version 21
	{
		{1156, 28, BlockInfo{4, 144, 116}, BlockInfo{4, 145, 117}}, // L
		{1156, 26, BlockInfo{17, 68, 42}, BlockInfo{}},             // M
		{1156, 28, BlockInfo{17, 50, 22}, BlockInfo{6, 51, 23}},    // Q
		{1156, 30, BlockInfo{19, 46, 16}, BlockInfo{6, 47, 17}},    // H
	},
	// Version 22
	{
		{1258, 28, BlockInfo{2, 139, 111}, BlockInfo{7, 140, 112}}, // L
		{1258, 28, BlockInfo{17, 74, 46}, BlockInfo{}},             // M
		{1258, 30, BlockInfo{7, 54, 24}, BlockInfo{16, 55, 25}},    // Q
		{1258, 24, BlockInfo{34, 37, 13}, BlockInfo{}},             // H
	},
	// Version 23
	{
		{1364, 30, BlockInfo{4, 151, 121}, BlockInfo{5, 152, 122}}, // L
		{1364, 28, BlockInfo{4, 75, 47}, BlockInfo{14, 76, 48}},    // M
		{1364, 30, BlockInfo{11, 54, 24}, BlockInfo{14, 55, 25}},   // Q
		{1364, 30, BlockInfo{16, 45, 15}, BlockInfo{14, 46, 16}},   // H
	},
	// Version 24
	{
		{1474, 30, BlockInfo{6, 147, 117}, BlockInfo{4, 148, 118}}, // L
		{1474, 28, BlockInfo{6, 73, 45}, BlockInfo{14, 74, 46}},    // M
		{1474, 30, BlockInfo{11, 54, 24}, BlockInfo{16, 55, 25}},   // Q
		{1474, 30, BlockInfo{30, 46, 16}, BlockInfo{2, 47, 17}},    // H
	},
	// Version 25
	{
		{1588, 26, BlockInfo{8, 132, 106}, BlockInfo{4, 133, 107}}, // L
		{1588, 28, BlockInfo{8, 75, 47}, BlockInfo{13, 76, 48}},    // M
		{1588, 30, BlockInfo{7, 54, 24}, BlockInfo{22, 55, 25}},    // Q
		{1588, 30, BlockInfo{22, 45, 15}, BlockInfo{13, 46, 16}},   // H
	},
	// Version 26
	{
		{1706, 28, BlockInfo{10, 142, 114}, BlockInfo{2, 143, 115}}, // L
		{1706, 28, BlockInfo{19, 74, 46}, BlockInfo{4, 75, 47}},     // M
		{1706, 28, BlockInfo{28, 50, 22}, BlockInfo{6, 51, 23}},     // Q
		{1706, 30, BlockInfo{33, 46, 16}, BlockInfo{4, 47, 17}},     // H
	},
	// Version 27
	{
		{1828, 30, BlockInfo{8, 152, 122}, BlockInfo{4, 153, 123}}, // L
		{1828, 28, BlockInfo{22, 73, 45}, BlockInfo{3, 74, 46}},    // M
		{1828, 30, BlockInfo{8, 53, 23}, BlockInfo{26, 54, 24}},    // Q
		{1828, 30, BlockInfo{12, 45, 15}, BlockInfo{28, 46, 16}},   // H
	},
	// Version 28
	{
		{1921, 30, BlockInfo{3, 147, 117}, BlockInfo{10, 148, 118}}, // L
		{1921, 28, BlockInfo{3, 73, 45}, BlockInfo{23, 74, 46}},     // M
		{1921, 30, BlockInfo{4, 54, 24}, BlockInfo{31, 55, 25}},     // Q
		{1921, 30, BlockInfo{11, 45, 15}, BlockInfo{31, 46, 16}},    // H
	},
	// Version 29
	{
		{2051, 30, BlockInfo{7, 146, 116}, BlockInfo{7, 147, 117}}, // L
		{2051, 28, BlockInfo{21, 73, 45}, BlockInfo{7, 74, 46}},    // M
		{2051, 30, BlockInfo{1, 53, 23}, BlockInfo{37, 54, 24}},    // Q
		{2051, 30, BlockInfo{19, 45, 15}, BlockInfo{26, 46, 16}},   // H
	},
	// Version 30
	{
		{2185, 30, BlockInfo{5, 145, 115}, BlockInfo{10, 146, 116}}, // L
		{2185, 28, BlockInfo{19, 75, 47}, BlockInfo{10, 76, 48}},    // M
		{2185, 30, BlockInfo{15, 54, 24}, BlockInfo{25, 55, 25}},    // Q
		{2185, 30, BlockInfo{23, 45, 15}, BlockInfo{25, 46, 16}},    // H
	},
	// Version 31
	{
		{2323, 30, BlockInfo{13, 145, 115}, BlockInfo{3, 146, 116}}, // L
		{2323, 28, BlockInfo{2, 74, 46}, BlockInfo{29, 75, 47}},     // M
		{2323, 30, BlockInfo{42, 54, 24}, BlockInfo{1, 55, 25}},     // Q
		{2323, 30, BlockInfo{23, 45, 15}, BlockInfo{28, 46, 16}},    // H
	},
	// Version 32
	{
		{2465, 30, BlockInfo{17, 145, 115}, BlockInfo{}},            // L
		{2465, 28, BlockInfo{10, 74, 46}, BlockInfo{23, 75, 47}},    // M
		{2465, 30, BlockInfo{10, 54, 24}, BlockInfo{35, 55, 25}},    // Q
		{2465, 30, BlockInfo{19, 45, 15}, BlockInfo{35, 46, 16}},    // H
	},
	// Version 33
	{
		{2611, 30, BlockInfo{17, 145, 115}, BlockInfo{1, 146, 116}}, // L
		{2611, 28, BlockInfo{14, 74, 46}, BlockInfo{21, 75, 47}},    // M
		{2611, 30, BlockInfo{29, 54, 24}, BlockInfo{19, 55, 25}},    // Q
		{2611, 30, BlockInfo{11, 45, 15}, BlockInfo{46, 46, 16}},    // H
	},
	// Version 34
	{
		{2761, 30, BlockInfo{13, 145, 115}, BlockInfo{6, 146, 116}}, // L
		{2761, 28, BlockInfo{14, 74, 46}, BlockInfo{23, 75, 47}},    // M
		{2761, 30, BlockInfo{44, 54, 24}, BlockInfo{7, 55, 25}},     // Q
		{2761, 30, BlockInfo{59, 46, 16}, BlockInfo{1, 47, 17}},     // H
	},
	// Version 35
	{
		{2876, 30, BlockInfo{12, 151, 121}, BlockInfo{7, 152, 122}}, // L
		{2876, 28, BlockInfo{12, 75, 47}, BlockInfo{26, 76, 48}},    // M
		{2876, 30, BlockInfo{39, 54, 24}, BlockInfo{14, 55, 25}},    // Q
		{2876, 30, BlockInfo{22, 45, 15}, BlockInfo{41, 46, 16}},    // H
	},
	// Version 36
	{
		{3034, 30, BlockInfo{6, 151, 121}, BlockInfo{14, 152, 122}}, // L
		{3034, 28, BlockInfo{6, 75, 47}, BlockInfo{34, 76, 48}},     // M
		{3034, 30, BlockInfo{46, 54, 24}, BlockInfo{10, 55, 25}},    // Q
		{3034, 30, BlockInfo{2, 45, 15}, BlockInfo{64, 46, 16}},     // H
	},
	// Version 37
	{
		{3196, 30, BlockInfo{17, 152, 122}, BlockInfo{4, 153, 123}}, // L
		{3196, 28, BlockInfo{29, 74, 46}, BlockInfo{14, 75, 47}},    // M
		{3196, 30, BlockInfo{49, 54, 24}, BlockInfo{10, 55, 25}},    // Q
		{3196, 30, BlockInfo{24, 45, 15}, BlockInfo{46, 46, 16}},    // H
	},
	// Version 38
	{
		{3362, 30, BlockInfo{4, 152, 122}, BlockInfo{18, 153, 123}}, // L
		{3362, 28, BlockInfo{13, 74, 46}, BlockInfo{32, 75, 47}},    // M
		{3362, 30, BlockInfo{48, 54, 24}, BlockInfo{14, 55, 25}},    // Q
		{3362, 30, BlockInfo{42, 45, 15}, BlockInfo{32, 46, 16}},    // H
	},
	// Version 39
	{
		{3532, 30, BlockInfo{20, 147, 117}, BlockInfo{4, 148, 118}}, // L
		{3532, 28, BlockInfo{40, 75, 47}, BlockInfo{7, 76, 48}},     // M
		{3532, 30, BlockInfo{43, 54, 24}, BlockInfo{22, 55, 25}},    // Q
		{3532, 30, BlockInfo{10, 45, 15}, BlockInfo{67, 46, 16}},    // H
	},
	// Version 40
	{
		{3706, 30, BlockInfo{19, 148, 118}, BlockInfo{6, 149, 119}}, // L
		{3706, 28, BlockInfo{18, 75, 47}, BlockInfo{31, 76, 48}},    // M
		{3706, 30, BlockInfo{34, 54, 24}, BlockInfo{34, 55, 25}},    // Q
		{3706, 30, BlockInfo{20, 45, 15}, BlockInfo{61, 46, 16}},    // H
	},
}

// GetECCInfo returns error correction info for version and level.
func GetECCInfo(version Version, ecl ErrorCorrectionLevel) ECCInfo {
	if version < 1 || version > 40 {
		return ECCInfo{}
	}
	return eccTable[version-1][ecl]
}

// GenerateECC splits data into blocks and generates ECC for each.
// Returns (dataBlocks, eccBlocks).
func GenerateECC(data []byte, eccInfo ECCInfo) ([][]byte, [][]byte) {
	totalBlocks := eccInfo.Group1.Count + eccInfo.Group2.Count
	dataBlocks := make([][]byte, totalBlocks)
	eccBlocks := make([][]byte, totalBlocks)

	offset := 0
	blockIndex := 0

	// Process Group 1 blocks
	for i := 0; i < eccInfo.Group1.Count; i++ {
		size := eccInfo.Group1.DataCodewords
		dataBlocks[blockIndex] = make([]byte, size)
		copy(dataBlocks[blockIndex], data[offset:offset+size])
		eccBlocks[blockIndex] = ReedSolomonEncode(dataBlocks[blockIndex], eccInfo.ECCPerBlock)
		offset += size
		blockIndex++
	}

	// Process Group 2 blocks
	for i := 0; i < eccInfo.Group2.Count; i++ {
		size := eccInfo.Group2.DataCodewords
		dataBlocks[blockIndex] = make([]byte, size)
		copy(dataBlocks[blockIndex], data[offset:offset+size])
		eccBlocks[blockIndex] = ReedSolomonEncode(dataBlocks[blockIndex], eccInfo.ECCPerBlock)
		offset += size
		blockIndex++
	}

	return dataBlocks, eccBlocks
}

// InterleaveBlocks interleaves data and ECC blocks for final placement.
func InterleaveBlocks(dataBlocks, eccBlocks [][]byte) []byte {
	var result []byte

	// Find max lengths
	maxDataLen := 0
	maxEccLen := 0
	for _, block := range dataBlocks {
		if len(block) > maxDataLen {
			maxDataLen = len(block)
		}
	}
	for _, block := range eccBlocks {
		if len(block) > maxEccLen {
			maxEccLen = len(block)
		}
	}

	// Interleave data: take byte i from each block, then byte i+1, etc.
	for i := 0; i < maxDataLen; i++ {
		for _, block := range dataBlocks {
			if i < len(block) {
				result = append(result, block[i])
			}
		}
	}

	// Interleave ECC: same pattern
	for i := 0; i < maxEccLen; i++ {
		for _, block := range eccBlocks {
			if i < len(block) {
				result = append(result, block[i])
			}
		}
	}

	return result
}

package encoder

import "testing"

func TestGetECCInfo(t *testing.T) {
	tests := []struct {
		version      Version
		ecl          ErrorCorrectionLevel
		wantTotal    int
		wantECCPer   int
		wantDataCap  int
	}{
		// Version 1
		{1, LevelL, 26, 7, 19},
		{1, LevelM, 26, 10, 16},
		{1, LevelQ, 26, 13, 13},
		{1, LevelH, 26, 17, 9},
		// Version 5 (has two groups)
		{5, LevelL, 134, 26, 108},
		{5, LevelQ, 134, 18, 62}, // 2*15 + 2*16 = 62
		// Version 40
		{40, LevelL, 3706, 30, 2956}, // 19*118 + 6*119 = 2956
	}

	for _, tt := range tests {
		info := GetECCInfo(tt.version, tt.ecl)

		if info.TotalCodewords != tt.wantTotal {
			t.Errorf("Version %d ECL %d: TotalCodewords = %d, want %d",
				tt.version, tt.ecl, info.TotalCodewords, tt.wantTotal)
		}

		if info.ECCPerBlock != tt.wantECCPer {
			t.Errorf("Version %d ECL %d: ECCPerBlock = %d, want %d",
				tt.version, tt.ecl, info.ECCPerBlock, tt.wantECCPer)
		}

		if info.DataCapacity() != tt.wantDataCap {
			t.Errorf("Version %d ECL %d: DataCapacity = %d, want %d",
				tt.version, tt.ecl, info.DataCapacity(), tt.wantDataCap)
		}
	}
}

func TestGenerateECC(t *testing.T) {
	// Version 1-L: 19 data bytes, 1 block, 7 ECC per block
	eccInfo := GetECCInfo(1, LevelL)
	data := make([]byte, 19)
	for i := range data {
		data[i] = byte(i + 1)
	}

	dataBlocks, eccBlocks := GenerateECC(data, eccInfo)

	// Should have 1 block
	if len(dataBlocks) != 1 {
		t.Fatalf("GenerateECC returned %d data blocks, want 1", len(dataBlocks))
	}
	if len(eccBlocks) != 1 {
		t.Fatalf("GenerateECC returned %d ECC blocks, want 1", len(eccBlocks))
	}

	// Data block should be 19 bytes
	if len(dataBlocks[0]) != 19 {
		t.Errorf("Data block has %d bytes, want 19", len(dataBlocks[0]))
	}

	// ECC block should be 7 bytes
	if len(eccBlocks[0]) != 7 {
		t.Errorf("ECC block has %d bytes, want 7", len(eccBlocks[0]))
	}
}

func TestGenerateECCMultipleBlocks(t *testing.T) {
	// Version 5-Q: 2 blocks of 15 data + 2 blocks of 16 data, 18 ECC each
	eccInfo := GetECCInfo(5, LevelQ)
	data := make([]byte, 62) // 2*15 + 2*16 = 62
	for i := range data {
		data[i] = byte(i)
	}

	dataBlocks, eccBlocks := GenerateECC(data, eccInfo)

	// Should have 4 blocks total
	expectedBlocks := eccInfo.Group1.Count + eccInfo.Group2.Count
	if len(dataBlocks) != expectedBlocks {
		t.Fatalf("GenerateECC returned %d data blocks, want %d", len(dataBlocks), expectedBlocks)
	}
	if len(eccBlocks) != expectedBlocks {
		t.Fatalf("GenerateECC returned %d ECC blocks, want %d", len(eccBlocks), expectedBlocks)
	}

	// Group 1 blocks should be 15 bytes
	for i := 0; i < eccInfo.Group1.Count; i++ {
		if len(dataBlocks[i]) != 15 {
			t.Errorf("Group1 block %d has %d bytes, want 15", i, len(dataBlocks[i]))
		}
	}

	// Group 2 blocks should be 16 bytes
	for i := eccInfo.Group1.Count; i < expectedBlocks; i++ {
		if len(dataBlocks[i]) != 16 {
			t.Errorf("Group2 block %d has %d bytes, want 16", i, len(dataBlocks[i]))
		}
	}

	// All ECC blocks should be 18 bytes
	for i, block := range eccBlocks {
		if len(block) != 18 {
			t.Errorf("ECC block %d has %d bytes, want 18", i, len(block))
		}
	}
}

func TestInterleaveBlocks(t *testing.T) {
	// Simple case: 2 blocks of equal length
	dataBlocks := [][]byte{
		{1, 2, 3},
		{4, 5, 6},
	}
	eccBlocks := [][]byte{
		{10, 20},
		{30, 40},
	}

	result := InterleaveBlocks(dataBlocks, eccBlocks)

	// Data interleaved: 1,4, 2,5, 3,6
	// ECC interleaved: 10,30, 20,40
	expected := []byte{1, 4, 2, 5, 3, 6, 10, 30, 20, 40}

	if len(result) != len(expected) {
		t.Fatalf("InterleaveBlocks returned %d bytes, want %d", len(result), len(expected))
	}

	for i, b := range expected {
		if result[i] != b {
			t.Errorf("result[%d] = %d, want %d", i, result[i], b)
		}
	}
}

func TestInterleaveBlocksUnequalLength(t *testing.T) {
	// Unequal blocks (like Group1 vs Group2)
	dataBlocks := [][]byte{
		{1, 2},    // Group 1: shorter
		{3, 4, 5}, // Group 2: longer
	}
	eccBlocks := [][]byte{
		{10, 20},
		{30, 40},
	}

	result := InterleaveBlocks(dataBlocks, eccBlocks)

	// Data: 1,3, 2,4, 5 (only block 2 has byte 3)
	// ECC: 10,30, 20,40
	expected := []byte{1, 3, 2, 4, 5, 10, 30, 20, 40}

	if len(result) != len(expected) {
		t.Fatalf("InterleaveBlocks returned %d bytes, want %d", len(result), len(expected))
	}

	for i, b := range expected {
		if result[i] != b {
			t.Errorf("result[%d] = %d, want %d", i, result[i], b)
		}
	}
}

func TestECCInfoDataCapacity(t *testing.T) {
	// Verify DataCapacity calculation for various versions
	tests := []struct {
		version Version
		ecl     ErrorCorrectionLevel
		want    int
	}{
		{1, LevelL, 19},
		{1, LevelH, 9},
		{5, LevelQ, 62},   // 2*15 + 2*16
		{10, LevelM, 216}, // 4*43 + 1*44 = 216
	}

	for _, tt := range tests {
		info := GetECCInfo(tt.version, tt.ecl)
		got := info.DataCapacity()
		if got != tt.want {
			t.Errorf("Version %d ECL %d: DataCapacity() = %d, want %d",
				tt.version, tt.ecl, got, tt.want)
		}
	}
}

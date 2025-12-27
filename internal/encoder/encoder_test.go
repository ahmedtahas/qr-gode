package encoder

import "testing"

func TestEncode(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		ecl     ErrorCorrectionLevel
		wantErr bool
	}{
		{"simple numeric", "12345", LevelL, false},
		{"simple alpha", "HELLO", LevelM, false},
		{"simple byte", "hello", LevelQ, false},
		{"url", "https://example.com", LevelH, false},
		{"empty", "", LevelL, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := New(tt.data, tt.ecl)
			matrix, err := enc.Encode()

			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && matrix == nil {
				t.Error("Encode() returned nil matrix")
				return
			}

			if matrix != nil {
				size := matrix.Size()
				if size < 21 {
					t.Errorf("Matrix size %d is less than minimum 21", size)
				}
				// Version 1 = 21, each version adds 4
				if (size-21)%4 != 0 {
					t.Errorf("Matrix size %d is not valid (should be 21 + 4*n)", size)
				}
			}
		})
	}
}

func TestEncodeMatrixStructure(t *testing.T) {
	// Test that a simple encode produces a valid matrix structure
	enc := New("HELLO", LevelM)
	matrix, err := enc.Encode()
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	size := matrix.Size()

	// Check finder patterns exist in corners
	// Top-left finder: (0,0) to (6,6) should have the pattern
	if !matrix.Get(0, 0).Dark {
		t.Error("Top-left finder corner should be dark")
	}
	if !matrix.Get(3, 3).Dark {
		t.Error("Top-left finder center should be dark")
	}

	// Top-right finder
	if !matrix.Get(size-1, 0).Dark {
		t.Error("Top-right finder corner should be dark")
	}

	// Bottom-left finder
	if !matrix.Get(0, size-1).Dark {
		t.Error("Bottom-left finder corner should be dark")
	}

	// Check timing patterns (row 6 and column 6 should alternate)
	// Starting from position 8, should be dark (even position)
	if !matrix.Get(8, 6).Dark {
		t.Error("Horizontal timing at (8,6) should be dark")
	}
	if matrix.Get(9, 6).Dark {
		t.Error("Horizontal timing at (9,6) should be light")
	}
}

func TestEncodeDifferentModes(t *testing.T) {
	// Numeric mode
	enc := New("1234567890", LevelL)
	matrix, err := enc.Encode()
	if err != nil {
		t.Errorf("Numeric encode failed: %v", err)
	}
	if matrix == nil {
		t.Error("Numeric encode returned nil matrix")
	}

	// Alphanumeric mode
	enc = New("HELLO WORLD", LevelL)
	matrix, err = enc.Encode()
	if err != nil {
		t.Errorf("Alphanumeric encode failed: %v", err)
	}
	if matrix == nil {
		t.Error("Alphanumeric encode returned nil matrix")
	}

	// Byte mode
	enc = New("Hello, World!", LevelL)
	matrix, err = enc.Encode()
	if err != nil {
		t.Errorf("Byte encode failed: %v", err)
	}
	if matrix == nil {
		t.Error("Byte encode returned nil matrix")
	}
}

func TestEncodeTooLong(t *testing.T) {
	// Create data that's too long for any QR version
	longData := make([]byte, 3000)
	for i := range longData {
		longData[i] = 'A'
	}

	enc := New(string(longData), LevelH)
	_, err := enc.Encode()

	if err == nil {
		t.Error("Expected error for data too long, got nil")
	}
}

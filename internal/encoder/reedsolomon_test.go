package encoder

import (
	"bytes"
	"testing"
)

func TestGFTables(t *testing.T) {
	// Verify exp and log tables are inverses
	for i := 1; i < 256; i++ {
		exp := expTable[logTable[i]]
		if exp != byte(i) {
			t.Errorf("exp[log[%d]] = %d, want %d", i, exp, i)
		}
	}

	// Verify expTable wraps correctly (alpha^255 = 1)
	if expTable[0] != 1 {
		t.Errorf("expTable[0] = %d, want 1", expTable[0])
	}

	// Known values for QR code's GF(256)
	knownValues := []struct {
		exp   int
		value byte
	}{
		{0, 1},
		{1, 2},
		{2, 4},
		{3, 8},
		{4, 16},
		{5, 32},
		{6, 64},
		{7, 128},
		{8, 29}, // After reduction by primitive polynomial
	}

	for _, kv := range knownValues {
		if expTable[kv.exp] != kv.value {
			t.Errorf("expTable[%d] = %d, want %d", kv.exp, expTable[kv.exp], kv.value)
		}
	}
}

func TestGFMul(t *testing.T) {
	tests := []struct {
		a, b, want byte
	}{
		{0, 5, 0},    // Anything times 0 is 0
		{5, 0, 0},
		{1, 7, 7},    // 1 is multiplicative identity
		{7, 1, 7},
		{2, 4, 8},    // Simple cases
		{2, 128, 29}, // 2 * 128 = 256 -> reduced by primitive poly
	}

	for _, tt := range tests {
		got := gfMul(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("gfMul(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}

	// Verify commutativity: a*b = b*a
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			ab := gfMul(byte(a), byte(b))
			ba := gfMul(byte(b), byte(a))
			if ab != ba {
				t.Errorf("gfMul not commutative: %d*%d=%d, %d*%d=%d", a, b, ab, b, a, ba)
			}
		}
	}
}

func TestGFDiv(t *testing.T) {
	// a / b * b = a
	for a := 1; a < 256; a++ {
		for b := 1; b < 256; b++ {
			quotient := gfDiv(byte(a), byte(b))
			product := gfMul(quotient, byte(b))
			if product != byte(a) {
				t.Errorf("gfDiv(%d, %d) * %d = %d, want %d", a, b, b, product, a)
			}
		}
	}
}

func TestGeneratorPolynomial(t *testing.T) {
	// Known generator polynomial for 2 ECC codewords:
	// g(x) = (x + α^0)(x + α^1) = x^2 + (α^0 + α^1)x + α^0*α^1
	//      = x^2 + 3x + 2
	// Coefficients stored high-to-low: [x^2, x^1, x^0]
	gen2 := GeneratorPolynomial(2)
	if len(gen2) != 3 {
		t.Fatalf("GeneratorPolynomial(2) has length %d, want 3", len(gen2))
	}
	expected2 := []byte{1, 3, 2} // x^2 + 3x + 2
	if !bytes.Equal(gen2, expected2) {
		t.Errorf("GeneratorPolynomial(2) = %v, want %v", gen2, expected2)
	}

	// Generator for 7 ECC codewords (used in Version 1-L)
	// Coefficients stored high-to-low
	gen7 := GeneratorPolynomial(7)
	if len(gen7) != 8 {
		t.Errorf("GeneratorPolynomial(7) has length %d, want 8", len(gen7))
	}
	// The leading coefficient is always 1, last is product of all α^i
	if gen7[0] != 1 {
		t.Errorf("GeneratorPolynomial(7)[0] = %d, want 1 (leading coefficient)", gen7[0])
	}
}

func TestReedSolomonEncode(t *testing.T) {
	// Test that RS encoding produces correct length output
	tests := []struct {
		name     string
		dataLen  int
		eccCount int
	}{
		{"Version 1-L", 19, 7},
		{"Version 1-M", 16, 10},
		{"Version 1-Q", 13, 13},
		{"Version 1-H", 9, 17},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, tt.dataLen)
			for i := range data {
				data[i] = byte(i + 1)
			}

			ecc := ReedSolomonEncode(data, tt.eccCount)

			if len(ecc) != tt.eccCount {
				t.Errorf("ReedSolomonEncode produced %d bytes, want %d", len(ecc), tt.eccCount)
			}

			// Verify ECC is not all zeros (basic sanity check)
			allZero := true
			for _, b := range ecc {
				if b != 0 {
					allZero = false
					break
				}
			}
			if allZero {
				t.Error("ReedSolomonEncode produced all zeros")
			}
		})
	}
}

func TestReedSolomonEncodeConsistency(t *testing.T) {
	// Same input should always produce same output
	data := []byte{0x40, 0x11, 0x22, 0x33, 0x44, 0x55}
	eccCount := 4

	ecc1 := ReedSolomonEncode(data, eccCount)
	ecc2 := ReedSolomonEncode(data, eccCount)

	if !bytes.Equal(ecc1, ecc2) {
		t.Errorf("ReedSolomonEncode not deterministic: %X vs %X", ecc1, ecc2)
	}
}

func TestReedSolomonEncodeVersion1L(t *testing.T) {
	// Version 1-L: 19 data codewords, 7 ECC codewords
	data := []byte{
		0x10, 0x20, 0x0C, 0x56, 0x61, 0x80, 0xEC, 0x11,
		0xEC, 0x11, 0xEC, 0x11, 0xEC, 0x11, 0xEC, 0x11,
		0xEC, 0x11, 0xEC,
	}

	ecc := ReedSolomonEncode(data, 7)

	if len(ecc) != 7 {
		t.Fatalf("ReedSolomonEncode produced %d bytes, want 7", len(ecc))
	}

	// Just verify we get the right length and non-zero output
	allZero := true
	for _, b := range ecc {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		t.Error("ReedSolomonEncode produced all zeros")
	}
}

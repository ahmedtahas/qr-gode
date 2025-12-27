package encoder

// GF256 implements Galois Field arithmetic for Reed-Solomon encoding.
// QR codes use GF(2^8) with primitive polynomial x^8 + x^4 + x^3 + x^2 + 1.

var (
	// expTable: alpha^i for i in 0..255
	expTable [256]byte
	// logTable: log_alpha(i) for i in 1..255
	logTable [256]byte
)

func init() {
	// Initialize exp and log tables for GF(2^8)
	// Primitive polynomial: x^8 + x^4 + x^3 + x^2 + 1 = 0x11D
	x := 1
	for i := 0; i < 256; i++ {
		expTable[i] = byte(x)
		if i < 255 {
			logTable[x] = byte(i)
		}
		// Multiply by alpha (which is 2 in this field)
		x <<= 1
		if x >= 256 {
			x ^= 0x11D // Reduce by primitive polynomial
		}
	}
}

// gfMul multiplies two numbers in GF(256).
func gfMul(a, b byte) byte {
	if a == 0 || b == 0 {
		return 0
	}
	// a * b = exp[(log[a] + log[b]) mod 255]
	sum := int(logTable[a]) + int(logTable[b])
	if sum >= 255 {
		sum -= 255
	}
	return expTable[sum]
}

// gfDiv divides two numbers in GF(256).
func gfDiv(a, b byte) byte {
	if a == 0 {
		return 0
	}
	if b == 0 {
		panic("division by zero in GF(256)")
	}
	// a / b = exp[(log[a] - log[b]) mod 255]
	diff := int(logTable[a]) - int(logTable[b])
	if diff < 0 {
		diff += 255
	}
	return expTable[diff]
}

// GeneratorPolynomial creates the generator polynomial for n ECC codewords.
// g(x) = (x - alpha^0)(x - alpha^1)...(x - alpha^(n-1))
// In GF(256), subtraction is the same as addition (XOR), so (x - α^i) = (x + α^i)
// Coefficients are stored high-to-low: [x^n, x^(n-1), ..., x^1, x^0]
func GeneratorPolynomial(eccCount int) []byte {
	// Start with g(x) = 1 (coefficient of x^0)
	gen := make([]byte, eccCount+1)
	gen[eccCount] = 1 // Constant term at the end

	// Multiply by (x + α^i) for i = 0 to eccCount-1
	for i := 0; i < eccCount; i++ {
		alphaI := expTable[i]
		// Multiply polynomial by (x + α^i)
		// Work from left to right (high degree to low)
		for j := 0; j < eccCount; j++ {
			gen[j] = gfMul(gen[j], alphaI) ^ gen[j+1]
		}
		gen[eccCount] = gfMul(gen[eccCount], alphaI)
	}

	return gen
}

// ReedSolomonEncode generates ECC codewords for the given data.
func ReedSolomonEncode(data []byte, eccCount int) []byte {
	gen := GeneratorPolynomial(eccCount)

	// Work buffer initialized with data, will hold remainder
	result := make([]byte, eccCount)

	// Polynomial division - process each data byte
	for _, b := range data {
		// XOR data byte with current first byte of result
		coef := b ^ result[0]

		// Shift result left by one position
		copy(result, result[1:])
		result[eccCount-1] = 0

		// XOR with generator polynomial scaled by coef
		if coef != 0 {
			for j := 0; j < eccCount; j++ {
				result[j] ^= gfMul(gen[j+1], coef)
			}
		}
	}

	return result
}

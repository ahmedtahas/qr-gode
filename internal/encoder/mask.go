package encoder

// MaskPattern represents one of the 8 QR mask patterns.
type MaskPattern int

const (
	Mask0 MaskPattern = iota // (row + col) mod 2 == 0
	Mask1                    // row mod 2 == 0
	Mask2                    // col mod 3 == 0
	Mask3                    // (row + col) mod 3 == 0
	Mask4                    // (row/2 + col/3) mod 2 == 0
	Mask5                    // (row*col) mod 2 + (row*col) mod 3 == 0
	Mask6                    // ((row*col) mod 2 + (row*col) mod 3) mod 2 == 0
	Mask7                    // ((row+col) mod 2 + (row*col) mod 3) mod 2 == 0
)

// ShouldFlip returns true if the module at (x, y) should be flipped.
func (mp MaskPattern) ShouldFlip(x, y int) bool {
	switch mp {
	case Mask0:
		return (x+y)%2 == 0
	case Mask1:
		return y%2 == 0
	case Mask2:
		return x%3 == 0
	case Mask3:
		return (x+y)%3 == 0
	case Mask4:
		return (y/2+x/3)%2 == 0
	case Mask5:
		return (x*y)%2+(x*y)%3 == 0
	case Mask6:
		return ((x*y)%2+(x*y)%3)%2 == 0
	case Mask7:
		return ((x+y)%2+(x*y)%3)%2 == 0
	}
	return false
}

// ApplyMask applies the mask pattern to a matrix.
// Only data modules are affected, not function patterns.
func ApplyMask(matrix *Matrix, pattern MaskPattern) {
	size := matrix.Size()
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			mod := matrix.Get(x, y)
			if !mod.Reserved && pattern.ShouldFlip(x, y) {
				mod.Dark = !mod.Dark
				matrix.Set(x, y, mod)
			}
		}
	}
}

// EvaluateMask scores a masked matrix using the 4 penalty rules.
// Lower score is better.
func EvaluateMask(matrix *Matrix) int {
	return penaltyRule1(matrix) + penaltyRule2(matrix) + penaltyRule3(matrix) + penaltyRule4(matrix)
}

// penaltyRule1: 5+ consecutive same-color modules in row/column
func penaltyRule1(matrix *Matrix) int {
	penalty := 0
	size := matrix.Size()

	// Check rows
	for y := 0; y < size; y++ {
		count := 1
		for x := 1; x < size; x++ {
			if matrix.Get(x, y).Dark == matrix.Get(x-1, y).Dark {
				count++
			} else {
				if count >= 5 {
					penalty += 3 + (count - 5)
				}
				count = 1
			}
		}
		if count >= 5 {
			penalty += 3 + (count - 5)
		}
	}

	// Check columns
	for x := 0; x < size; x++ {
		count := 1
		for y := 1; y < size; y++ {
			if matrix.Get(x, y).Dark == matrix.Get(x, y-1).Dark {
				count++
			} else {
				if count >= 5 {
					penalty += 3 + (count - 5)
				}
				count = 1
			}
		}
		if count >= 5 {
			penalty += 3 + (count - 5)
		}
	}

	return penalty
}

// penaltyRule2: 2x2 blocks of same color
func penaltyRule2(matrix *Matrix) int {
	penalty := 0
	size := matrix.Size()

	for y := 0; y < size-1; y++ {
		for x := 0; x < size-1; x++ {
			dark := matrix.Get(x, y).Dark
			if matrix.Get(x+1, y).Dark == dark &&
				matrix.Get(x, y+1).Dark == dark &&
				matrix.Get(x+1, y+1).Dark == dark {
				penalty += 3
			}
		}
	}

	return penalty
}

// penaltyRule3: Finder-like pattern (1011101 with 4 white on either side)
func penaltyRule3(matrix *Matrix) int {
	penalty := 0
	size := matrix.Size()

	// Pattern: dark-light-dark-dark-dark-light-dark (1011101)
	// With 4 light modules on one side: 00001011101 or 10111010000

	for y := 0; y < size; y++ {
		for x := 0; x < size-10; x++ {
			if matchesFinderPattern(matrix, x, y, true) {
				penalty += 40
			}
		}
	}

	for x := 0; x < size; x++ {
		for y := 0; y < size-10; y++ {
			if matchesFinderPattern(matrix, x, y, false) {
				penalty += 40
			}
		}
	}

	return penalty
}

func matchesFinderPattern(matrix *Matrix, x, y int, horizontal bool) bool {
	// Patterns to match (0=light, 1=dark):
	// 10111010000 or 00001011101
	pattern1 := []bool{true, false, true, true, true, false, true, false, false, false, false}
	pattern2 := []bool{false, false, false, false, true, false, true, true, true, false, true}

	matches1 := true
	matches2 := true

	for i := 0; i < 11; i++ {
		var dark bool
		if horizontal {
			dark = matrix.Get(x+i, y).Dark
		} else {
			dark = matrix.Get(x, y+i).Dark
		}
		if dark != pattern1[i] {
			matches1 = false
		}
		if dark != pattern2[i] {
			matches2 = false
		}
	}

	return matches1 || matches2
}

// penaltyRule4: Proportion of dark modules
func penaltyRule4(matrix *Matrix) int {
	size := matrix.Size()
	darkCount := 0
	total := size * size

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if matrix.Get(x, y).Dark {
				darkCount++
			}
		}
	}

	// Calculate percentage (0-100)
	percent := (darkCount * 100) / total

	// Find deviation from 50%
	// Penalty = 10 * floor(|percent - 50| / 5)
	deviation := percent - 50
	if deviation < 0 {
		deviation = -deviation
	}

	return (deviation / 5) * 10
}

// SelectBestMask tries all 8 masks and returns the best one.
func SelectBestMask(matrix *Matrix) MaskPattern {
	bestMask := Mask0
	bestScore := -1

	for mask := Mask0; mask <= Mask7; mask++ {
		// Clone and apply mask
		testMatrix := matrix.Clone()
		ApplyMask(testMatrix, mask)

		// Evaluate
		score := EvaluateMask(testMatrix)

		if bestScore == -1 || score < bestScore {
			bestScore = score
			bestMask = mask
		}
	}

	return bestMask
}

package encoder

import "strings"

// Mode represents the encoding mode for QR data.
type Mode int

const (
	ModeNumeric      Mode = iota // 0-9 only
	ModeAlphanumeric             // 0-9, A-Z, space, $%*+-./:
	ModeByte                     // Any 8-bit data (UTF-8)
	ModeKanji                    // Kanji characters
)

func (m Mode) ModeIndicator() uint8 {
	switch m {
	case ModeNumeric:
		return 1
	case ModeAlphanumeric:
		return 2
	case ModeByte:
		return 4
	case ModeKanji:
		return 8
	}
	return 0
}

func (m Mode) CharCountBits(version Version) int {
	if version < 1 || version > 40 {
		return 0
	}

	// 3 ranges: 1-9, 10-26, 27-40
	var i int
	if version <= 9 {
		i = 0
	} else if version <= 26 {
		i = 1
	} else {
		i = 2
	}

	switch m {
	case ModeNumeric:
		return []int{10, 12, 14}[i]
	case ModeAlphanumeric:
		return []int{9, 11, 13}[i]
	case ModeByte:
		return []int{8, 16, 16}[i]
	case ModeKanji:
		return []int{8, 10, 12}[i]
	}
	return 0
}

func (m Mode) BitsPerChar() float64 {
	switch m {
	case ModeNumeric:
		return 10.0 / 3.0 // 3 digits in 10 bits
	case ModeAlphanumeric:
		return 11.0 / 2.0 // 2 chars in 11 bits
	case ModeByte:
		return 8.0 // 1 byte = 8 bits
	case ModeKanji:
		return 13.0 // 1 character = 13 bits
	}
	return 0
}

func AnalyzeData(data string) Mode {
	isNum := true
	isAlph := true
	for _, r := range data {
		if isAlph && (r < '0' || r > '9') {
			isNum = false
			if (r < 'A' || r > 'Z') && !strings.ContainsRune(" $%*+-./:", r) {
				isAlph = false
			}
		}
	}
	if isNum {
		return ModeNumeric
	} else if isAlph {
		return ModeAlphanumeric
	}
	return ModeByte
}

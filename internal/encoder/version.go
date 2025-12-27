package encoder

import "errors"

// Version represents a QR code version (1-40).
// Version determines the size: (version * 4) + 17 modules per side.
type Version int

// Size returns the number of modules per side for this version.
func (v Version) Size() int {
	return int(v)*4 + 17
}

// DetermineVersion finds the minimum version that can encode the data
// with the given mode and error correction level.
func DetermineVersion(dataLen int, mode Mode, ecl ErrorCorrectionLevel) (Version, error) {
	for v := 1; v <= 40; v++ {
		version := Version(v)
		capacity := GetECCInfo(version, ecl).DataCapacity() * 8 // bits

		// Calculate bits needed: mode indicator + char count + data
		bitsNeeded := 4 + mode.CharCountBits(version)
		bitsNeeded += int(float64(dataLen) * mode.BitsPerChar())

		if bitsNeeded <= capacity {
			return version, nil
		}
	}
	return 0, errors.New("data too long for any QR version")
}

# QR-Gode Progress

## Completed

### `internal/encoder/mode.go`
- [x] `Mode` type and constants (Numeric, Alphanumeric, Byte, Kanji)
- [x] `ModeIndicator()` - returns 4-bit mode indicator
- [x] `CharCountBits(version)` - returns character count bits based on mode and version
- [x] `BitsPerChar()` - returns bits per character for capacity estimation
- [x] `AnalyzeData(data)` - determines best encoding mode for input

### `internal/encoder/data.go`
- [x] `BitStream` struct and methods
  - [x] `AppendBits(value, n)` - adds n bits MSB first
  - [x] `AppendByte(b)` - adds 8 bits
  - [x] `Len()` - returns bit count
  - [x] `Bytes()` - converts to byte slice with padding
- [x] `encodeNumeric()` - 3 digits → 10 bits, 2 → 7, 1 → 4
- [x] `encodeAlphanumeric()` - pairs → 11 bits (first*45 + second), odd → 6 bits
- [x] `encodeByte()` - 8 bits per byte
- [x] `alphanumericValue()` - maps char to 0-44 value (note: typo `alhanumericValue`)

### `internal/encoder/encoder.go`
- [x] `Encoder` struct
- [x] `ErrorCorrectionLevel` constants (L, M, Q, H)
- [x] `New()` constructor

### `internal/encoder/data.go`
- [x] `EncodeData(data, mode, version)` - main encoding function
  1. Add mode indicator (4 bits)
  2. Add character count indicator (varies by mode/version)
  3. Encode data using appropriate helper
  4. Add terminator (4 zero bits)
  5. Pad to byte boundary
  6. TODO: Add pad bytes to fill capacity (needs capacity from caller)

### `internal/encoder/ecc.go`
- [x] `BlockInfo` struct (Count, TotalCodewords, DataCodewords)
- [x] `ECCInfo` struct (TotalCodewords, ECCPerBlock, Group1, Group2)
- [x] `ECCInfo.DataCapacity()` - returns total data codewords
- [x] `eccTable` - full capacity table for versions 1-40, all ECL levels
- [x] `GetECCInfo(version, ecl)` - lookup from table

## Next Up

### `internal/encoder/encoder.go`
- [ ] `Encode()` - full encoding pipeline
  1. Analyze data to determine best mode
  2. Determine minimum version that fits data + error correction
  3. Encode data to bit stream
  4. Generate error correction codewords
  5. Structure data (interleave blocks)
  6. Create matrix and place function patterns
  7. Place data modules
  8. Apply masking and select best mask
  9. Add format and version information

## Files to Review
- `version.go` - version capacity tables
- `ecc.go` - error correction
- `reedsolomon.go` - Reed-Solomon encoding
- `matrix.go` - QR matrix operations
- `mask.go` - masking patterns
- `format.go` - format/version info encoding

## Notes
- Alphanumeric encoding uses base-45: `first * 45 + second` fits in 11 bits (max 2024)
- Numeric encoding: 3 digits fit in 10 bits (max 999)
- Kanji is 13 bits per character (no pairing)

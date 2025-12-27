package encoder

import "strconv"

type BitStream struct {
	bits []bool
}

func NewBitStream() *BitStream {
	return &BitStream{
		bits: make([]bool, 0),
	}
}

func (bs *BitStream) AppendBits(value uint, n int) {
	for i := n - 1; i >= 0; i-- {
		bit := value & (1 << i)
		bs.bits = append(bs.bits, bit != 0)
	}

}

func (bs *BitStream) AppendByte(b byte) {
	bs.AppendBits(uint(b), 8)
}

func (bs *BitStream) Len() int {
	return len(bs.bits)
}

func (bs *BitStream) Bytes() []byte {
	temp := byte(0)
	bytes := make([]byte, 0, (len(bs.bits)+7)/8)
	for i, value := range bs.bits {
		temp = temp << 1
		if value {
			temp |= 1
		}
		if i%8 == 7 {
			bytes = append(bytes, temp)
			temp = 0
		}
	}
	if len(bs.bits)%8 != 0 {
		leftover := len(bs.bits) % 8
		temp = temp << (8 - leftover)
		bytes = append(bytes, temp)
	}
	return bytes
}

func EncodeData(data string, mode Mode, version Version) *BitStream {
	bs := NewBitStream()

	// 1. Mode indicator (4 bits)
	bs.AppendBits(uint(mode.ModeIndicator()), 4)

	// 2. Character count indicator
	charCountBits := mode.CharCountBits(version)
	bs.AppendBits(uint(len(data)), charCountBits)

	// 3. Encode the actual data
	switch mode {
	case ModeNumeric:
		encodeNumeric(bs, data)
	case ModeAlphanumeric:
		encodeAlphanumeric(bs, data)
	case ModeByte:
		encodeByte(bs, data)
	default:
		// Kanji encoding not implemented
	}

	// 4. Terminator (up to 4 zero bits)
	bs.AppendBits(0, 4)

	// 5. Pad to byte boundary
	if bs.Len()%8 != 0 {
		bs.AppendBits(0, 8-bs.Len()%8)
	}

	return bs
}

// EncodeDataWithPadding encodes data and pads to the required capacity.
func EncodeDataWithPadding(data string, mode Mode, version Version, capacityBytes int) *BitStream {
	bs := EncodeData(data, mode, version)

	// Add padding bytes to fill capacity (alternating 0xEC and 0x11)
	padBytes := []byte{0xEC, 0x11}
	padIndex := 0
	for bs.Len() < capacityBytes*8 {
		bs.AppendByte(padBytes[padIndex])
		padIndex = (padIndex + 1) % 2
	}

	return bs
}

func encodeNumeric(bs *BitStream, data string) {
	i := 0
	for i+3 <= len(data) {
		num, _ := strconv.Atoi(data[i : i+3])
		bs.AppendBits(uint(num), 10)
		i += 3
	}
	remaining := len(data) - i
	num, _ := strconv.Atoi(data[i : i+remaining])
	switch remaining {
	case 2:
		bs.AppendBits(uint(num), 7)
	case 1:
		bs.AppendBits(uint(num), 4)
	}
}

func encodeAlphanumeric(bs *BitStream, data string) {
	i := 0
	for i+2 <= len(data) {
		first := alphanumericValue(data[i])
		second := alphanumericValue(data[i+1])
		combined := first*45 + second
		bs.AppendBits(uint(combined), 11)
		i += 2
	}
	if len(data)%2 == 1 {
		last := alphanumericValue(data[len(data)-1])
		bs.AppendBits(uint(last), 6)
	}
}

func encodeByte(bs *BitStream, data string) {
	for _, b := range []byte(data) {
		bs.AppendByte(b)
	}
}

func alphanumericValue(char byte) int {
	switch {
	case char >= '0' && char <= '9':
		return int(char - '0')
	case char >= 'A' && char <= 'Z':
		return int(char - 'A' + 10)
	default:
		specials := " $%*+-./:"
		for i, c := range specials {
			if byte(c) == char {
				return 36 + i
			}
		}
	}
	return 0
}

package encoder

import (
	"testing"
)

func TestBitStream_AppendBits(t *testing.T) {
	tests := []struct {
		name     string
		value    uint
		n        int
		wantBits []bool
	}{
		{"5 in 4 bits", 5, 4, []bool{false, true, false, true}},                             // 0101
		{"5 in 8 bits", 5, 8, []bool{false, false, false, false, false, true, false, true}}, // 00000101
		{"255 in 8 bits", 255, 8, []bool{true, true, true, true, true, true, true, true}},   // 11111111
		{"1 in 1 bit", 1, 1, []bool{true}},
		{"0 in 4 bits", 0, 4, []bool{false, false, false, false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitStream()
			bs.AppendBits(tt.value, tt.n)
			if len(bs.bits) != len(tt.wantBits) {
				t.Errorf("got %d bits, want %d", len(bs.bits), len(tt.wantBits))
				return
			}
			for i, want := range tt.wantBits {
				if bs.bits[i] != want {
					t.Errorf("bit %d: got %v, want %v", i, bs.bits[i], want)
				}
			}
		})
	}
}

func TestBitStream_AppendByte(t *testing.T) {
	bs := NewBitStream()
	bs.AppendByte(0xA5) // 10100101

	want := []bool{true, false, true, false, false, true, false, true}
	if len(bs.bits) != 8 {
		t.Fatalf("got %d bits, want 8", len(bs.bits))
	}
	for i, w := range want {
		if bs.bits[i] != w {
			t.Errorf("bit %d: got %v, want %v", i, bs.bits[i], w)
		}
	}
}

func TestBitStream_Bytes(t *testing.T) {
	tests := []struct {
		name      string
		bits      []bool
		wantBytes []byte
	}{
		{
			"8 bits exact",
			[]bool{true, false, true, false, false, true, false, true}, // 10100101 = 165
			[]byte{165},
		},
		{
			"16 bits exact",
			[]bool{true, true, true, true, false, false, false, false, // 11110000 = 240
				false, false, false, false, true, true, true, true}, // 00001111 = 15
			[]byte{240, 15},
		},
		{
			"3 bits with padding",
			[]bool{true, false, true}, // 101 -> 10100000 = 160
			[]byte{160},
		},
		{
			"10 bits with padding",
			[]bool{true, false, true, false, true, false, true, false, // 10101010 = 170
				true, true}, // 11 -> 11000000 = 192
			[]byte{170, 192},
		},
		{
			"empty",
			[]bool{},
			[]byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := &BitStream{bits: tt.bits}
			got := bs.Bytes()
			if len(got) != len(tt.wantBytes) {
				t.Errorf("got %d bytes, want %d", len(got), len(tt.wantBytes))
				return
			}
			for i, want := range tt.wantBytes {
				if got[i] != want {
					t.Errorf("byte %d: got %d, want %d", i, got[i], want)
				}
			}
		})
	}
}

func TestBitStream_Len(t *testing.T) {
	bs := NewBitStream()
	if bs.Len() != 0 {
		t.Errorf("empty stream: got %d, want 0", bs.Len())
	}

	bs.AppendBits(5, 4)
	if bs.Len() != 4 {
		t.Errorf("after 4 bits: got %d, want 4", bs.Len())
	}

	bs.AppendByte(0xFF)
	if bs.Len() != 12 {
		t.Errorf("after 4+8 bits: got %d, want 12", bs.Len())
	}
}

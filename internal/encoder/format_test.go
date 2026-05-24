package encoder

import "testing"

// Canonical format-info values from the QR spec (Table C.1, ISO/IEC 18004).
// Each entry: (ECL, mask) -> 15-bit pattern.
func TestFormatInfo_KnownValues(t *testing.T) {
	cases := []struct {
		ecl  ErrorCorrectionLevel
		mask MaskPattern
		want uint16
	}{
		// L (eclBits=01)
		{LevelL, 0, 0x77c4},
		{LevelL, 1, 0x72f3},
		{LevelL, 7, 0x6976},
		// M (eclBits=00)
		{LevelM, 0, 0x5412},
		{LevelM, 4, 0x45f9},
		// Q (eclBits=11)
		{LevelQ, 0, 0x355f},
		{LevelQ, 7, 0x2bed},
		// H (eclBits=10)
		{LevelH, 0, 0x1689},
		{LevelH, 4, 0x0762},
	}
	for _, tc := range cases {
		got := FormatInfo(tc.ecl, tc.mask)
		if got != tc.want {
			t.Errorf("FormatInfo(ecl=%d, mask=%d) = 0x%04x, want 0x%04x",
				tc.ecl, tc.mask, got, tc.want)
		}
	}
}

func TestPlaceFormatInfo_PlacesAtKnownPositions(t *testing.T) {
	m := NewMatrix(1) // 21x21 matrix
	info := FormatInfo(LevelM, 0)
	PlaceFormatInfo(m, info)

	// Format info must be placed at column 8 around the top-left finder.
	// Check that at least one cell got marked as ModuleFormatInfo.
	found := false
	for i := range 9 {
		if m.Get(8, i).Type == ModuleFormatInfo && m.Get(8, i).Reserved {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected format-info modules around top-left finder")
	}

	// The mirrored location near top-right finder must also be populated.
	size := m.Size()
	if m.Get(size-1, 8).Type != ModuleFormatInfo {
		t.Error("expected format info at top-right strip")
	}
}

func TestVersionInfo_OutOfRange(t *testing.T) {
	// Versions <7 don't carry version-info blocks.
	for _, v := range []Version{0, 1, 6} {
		if got := VersionInfo(v); got != 0 {
			t.Errorf("VersionInfo(%d) = %#x, want 0", v, got)
		}
	}
	// Above 40 should also return 0.
	if got := VersionInfo(41); got != 0 {
		t.Errorf("VersionInfo(41) = %#x, want 0", got)
	}
}

func TestVersionInfo_TableLookup(t *testing.T) {
	// Canonical values from the QR spec for the lowest and a middle version.
	cases := map[Version]uint32{
		7:  0x07C94,
		10: 0x0A4D3,
		40: 0x28C69,
	}
	for v, want := range cases {
		if got := VersionInfo(v); got != want {
			t.Errorf("VersionInfo(%d) = %#x, want %#x", v, got, want)
		}
	}
}

func TestPlaceVersionInfo_SkipsLowVersions(t *testing.T) {
	m := NewMatrix(1)
	PlaceVersionInfo(m, 1, 0xdeadbeef) // no-op for version <7
	// All cells should still be the zero-value Module (no Reserved set by us).
	if m.Get(0, 0).Reserved {
		t.Error("version <7 should not have written version info")
	}
}

func TestPlaceVersionInfo_PlacesTwoBlocks(t *testing.T) {
	v := Version(7)
	m := NewMatrix(v)
	info := VersionInfo(v)
	PlaceVersionInfo(m, v, info)

	size := m.Size()
	// Location 1: columns size-11..size-9, rows 0..5
	// Location 2: rows size-11..size-9, columns 0..5
	loc1Set, loc2Set := false, false
	for i := range 6 {
		for j := range 3 {
			if m.Get(size-11+j, i).Type == ModuleVersionInfo {
				loc1Set = true
			}
			if m.Get(i, size-11+j).Type == ModuleVersionInfo {
				loc2Set = true
			}
		}
	}
	if !loc1Set {
		t.Error("expected version info near top-right finder")
	}
	if !loc2Set {
		t.Error("expected version info near bottom-left finder")
	}
}

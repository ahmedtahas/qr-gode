package main

import (
	"fmt"

	"github.com/ahmedtahas/qr-gode/internal/encoder"
)

func main() {
	tests := []struct {
		name string
		data string
	}{
		{"numeric (20 digits)", "12345678901234567890"},
		{"alphanumeric", "HELLO WORLD 123"},
		{"google.com", "https://google.com"},
		{"long text", "The quick brown fox jumps over the lazy dog. This is a longer message to test byte mode encoding with more data!"},
		{"lorem ipsum", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur."},
		{"huge lorem", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Curabitur pretium tincidunt lacus. Nulla gravida orci a odio."},
	}

	fmt.Println("QR Code Version Analysis (Error Correction: M)")
	fmt.Println("=" + string(make([]byte, 60)))

	for _, tt := range tests {
		mode := encoder.AnalyzeData(tt.data)
		version, _ := encoder.DetermineVersion(len(tt.data), mode, encoder.LevelM)
		modeName := []string{"Numeric", "Alphanumeric", "Byte", "Kanji"}[mode]
		fmt.Printf("%-20s: %3d chars, %-12s -> Version %2d (%dx%d)\n",
			tt.name, len(tt.data), modeName, version, version.Size(), version.Size())
	}

	// Test version 7+ (needs version info blocks)
	fmt.Println("\nTesting Version 7+ (requires version info blocks):")
	bigData := make([]byte, 200)
	for i := range bigData {
		bigData[i] = 'A' + byte(i%26)
	}
	mode := encoder.AnalyzeData(string(bigData))
	version, _ := encoder.DetermineVersion(len(bigData), mode, encoder.LevelM)
	fmt.Printf("200 chars alphanumeric: Version %d (%dx%d)\n", version, version.Size(), version.Size())

	// Force a version 7+
	bigData = make([]byte, 370)
	for i := range bigData {
		bigData[i] = 'a' + byte(i%26) // lowercase = byte mode
	}
	mode = encoder.AnalyzeData(string(bigData))
	version, _ = encoder.DetermineVersion(len(bigData), mode, encoder.LevelM)
	fmt.Printf("370 chars byte mode:    Version %d (%dx%d)\n", version, version.Size(), version.Size())
}

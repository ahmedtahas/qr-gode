package encoder

import "testing"

func TestAnalyzeData(t *testing.T) {
	tests := []struct {
		name string
		data string
		want Mode
	}{
		// Numeric mode - only digits
		{"digits only", "12345", ModeNumeric},
		{"phone number", "5551234567", ModeNumeric},
		{"single digit", "0", ModeNumeric},

		// Alphanumeric mode - uppercase, digits, special chars
		{"uppercase", "HELLO", ModeAlphanumeric},
		{"uppercase with digits", "ABC123", ModeAlphanumeric},
		{"with space", "HELLO WORLD", ModeAlphanumeric},
		{"with special chars", "TOTAL: $99", ModeAlphanumeric},
		{"url uppercase", "HTTPS://EXAMPLE.COM", ModeAlphanumeric},
		{"all special", "$%*+-./:", ModeAlphanumeric},

		// Byte mode - anything else
		{"lowercase", "hello", ModeByte},
		{"mixed case", "Hello", ModeByte},
		{"lowercase url", "https://example.com", ModeByte},
		{"emoji", "QR📱CODE", ModeByte},
		{"special not in set", "test@email.com", ModeByte},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AnalyzeData(tt.data)
			if got != tt.want {
				t.Errorf("AnalyzeData(%q) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

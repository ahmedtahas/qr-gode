package qrgode

import (
	"strings"
	"testing"
)

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		fn          func() ([]byte, error)
		shouldError bool
	}{
		{
			name: "empty data",
			fn: func() ([]byte, error) {
				return New("").SVG()
			},
			shouldError: true, // TODO: should error but currently doesn't
		},
		{
			name: "very long data",
			fn: func() ([]byte, error) {
				data := strings.Repeat("a", 3000)
				return New(data).SVG()
			},
			shouldError: true, // exceeds QR capacity
		},
		{
			name: "invalid shape falls back to square",
			fn: func() ([]byte, error) {
				return New("test").Shape(Shape("invalid_shape")).SVG()
			},
			shouldError: false,
		},
		{
			name: "zero size",
			fn: func() ([]byte, error) {
				return New("test").Size(0).SVG()
			},
			shouldError: true, // validation rejects non-positive size
		},
		{
			name: "negative size",
			fn: func() ([]byte, error) {
				return New("test").Size(-100).SVG()
			},
			shouldError: true, // validation rejects non-positive size
		},
		{
			name: "non-existent logo",
			fn: func() ([]byte, error) {
				return New("test").Logo("/nonexistent/logo.png").SVG()
			},
			shouldError: true,
		},
		{
			name: "non-existent finder image",
			fn: func() ([]byte, error) {
				return New("test").FinderImage("/nonexistent/finder.png").SVG()
			},
			shouldError: true,
		},
		{
			name: "special characters in data",
			fn: func() ([]byte, error) {
				return New("日本語テスト").SVG()
			},
			shouldError: false,
		},
		{
			name: "URL with special chars",
			fn: func() ([]byte, error) {
				return New("https://example.com/path?q=test&foo=bar#anchor").SVG()
			},
			shouldError: false,
		},
		{
			name: "all error correction levels",
			fn: func() ([]byte, error) {
				for _, ecl := range []ErrorCorrectionLevel{LevelL, LevelM, LevelQ, LevelH} {
					_, err := New("test").ErrorCorrection(ecl).SVG()
					if err != nil {
						return nil, err
					}
				}
				return []byte("ok"), nil
			},
			shouldError: false,
		},
		{
			name: "all shapes",
			fn: func() ([]byte, error) {
				shapes := []Shape{ShapeSquare, ShapeCircle, ShapeRounded, ShapeDiamond, ShapeDot, ShapeStar, ShapeHeart}
				for _, s := range shapes {
					_, err := New("test").Shape(s).SVG()
					if err != nil {
						return nil, err
					}
				}
				return []byte("ok"), nil
			},
			shouldError: false,
		},
		{
			name: "gradient with single color",
			fn: func() ([]byte, error) {
				return New("test").LinearGradient(45, "#ff0000").SVG()
			},
			shouldError: false,
		},
		{
			name: "gradient with many colors",
			fn: func() ([]byte, error) {
				return New("test").LinearGradient(45, "#ff0000", "#00ff00", "#0000ff", "#ffff00", "#ff00ff").SVG()
			},
			shouldError: false,
		},
		{
			name: "radial gradient",
			fn: func() ([]byte, error) {
				return New("test").RadialGradient(0.5, 0.5, "#ff0000", "#0000ff").SVG()
			},
			shouldError: false,
		},
		{
			name: "very small QR",
			fn: func() ([]byte, error) {
				return New("a").Size(10).SVG()
			},
			shouldError: false,
		},
		{
			name: "very large QR",
			fn: func() ([]byte, error) {
				return New("test").Size(10000).SVG()
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PANIC: %v", r)
				}
			}()

			result, err := tt.fn()

			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(result) == 0 {
					t.Error("expected non-empty result")
				}
			}
		})
	}
}

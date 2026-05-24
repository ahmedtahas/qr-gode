package shapes

import (
	"strings"
	"testing"
)

func TestBuiltinRegistry(t *testing.T) {
	for _, name := range []string{"square", "circle", "rounded", "diamond", "dot", "star", "heart"} {
		s := Get(name)
		if s == nil {
			t.Errorf("Get(%q) returned nil — built-in shape missing", name)
			continue
		}
		if s.Name() != name {
			t.Errorf("Name() = %q, want %q", s.Name(), name)
		}
		path := s.SVGPath()
		if path == "" {
			t.Errorf("%s: SVGPath is empty", name)
		}
		if path[0] != 'M' && path[0] != 'm' {
			t.Errorf("%s: SVGPath doesn't start with moveto: %q", name, path)
		}
	}
}

func TestGet_Unknown(t *testing.T) {
	if got := Get("nonexistent-shape"); got != nil {
		t.Errorf("Get unknown = %v, want nil", got)
	}
}

func TestRegister(t *testing.T) {
	defer func(prev Shape) { Registry["test-shape"] = prev }(Registry["test-shape"])

	Register(FromPath("test-shape", "M0 0L1 1"))
	got := Get("test-shape")
	if got == nil {
		t.Fatal("Get after Register returned nil")
	}
	if got.Name() != "test-shape" {
		t.Errorf("Name = %q", got.Name())
	}
	delete(Registry, "test-shape")
}

func TestFromPath(t *testing.T) {
	s := FromPath("foo", "M0 0h1v1z")
	if s.Name() != "foo" {
		t.Errorf("Name = %s", s.Name())
	}
	if s.SVGPath() != "M0 0h1v1z" {
		t.Errorf("SVGPath = %s", s.SVGPath())
	}
}

func TestParsePath(t *testing.T) {
	cases := []struct {
		in      string
		wantErr bool
	}{
		{"M0 0h1v1z", false},
		{"m0 0l1 1", false},
		{"", true},
		{"L0 0", true},  // doesn't start with M
		{"0 0L1 1", true},
	}
	for _, tc := range cases {
		got, err := ParsePath(tc.in)
		if tc.wantErr && err == nil {
			t.Errorf("ParsePath(%q) expected error", tc.in)
		}
		if !tc.wantErr {
			if err != nil {
				t.Errorf("ParsePath(%q) unexpected error: %v", tc.in, err)
			}
			if got != tc.in {
				t.Errorf("ParsePath(%q) = %q", tc.in, got)
			}
		}
	}
}

func TestIsValidPath(t *testing.T) {
	cases := map[string]bool{
		"M0 0h1z":   true,
		"m0 0l1 1z": true,
		"":          false,
		"L0 0":      false,
		"x":         false,
	}
	for in, want := range cases {
		if got := IsValidPath(in); got != want {
			t.Errorf("IsValidPath(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestResolvePath(t *testing.T) {
	// Registered name resolves to the built-in shape.
	s, err := ResolvePath("square")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name() != "square" {
		t.Errorf("Name = %s, want square", s.Name())
	}

	// Raw SVG path becomes a custom shape.
	s, err = ResolvePath("M0 0h1v1z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name() != "custom" {
		t.Errorf("custom Name = %s, want custom", s.Name())
	}

	// Invalid input bubbles up the parse error.
	if _, err := ResolvePath(""); err == nil {
		t.Error("expected error for empty input")
	}
	if _, err := ResolvePath("L0 0"); err == nil {
		t.Error("expected error for path without M")
	}
}

func TestSVGPath_ContainsCoordinates(t *testing.T) {
	// Sanity-check that built-in shapes generate paths with numeric content.
	for _, name := range []string{"square", "circle", "rounded", "diamond", "dot", "star", "heart"} {
		path := Get(name).SVGPath()
		if !strings.ContainsAny(path, "0123456789") {
			t.Errorf("%s path has no numbers: %q", name, path)
		}
	}
}

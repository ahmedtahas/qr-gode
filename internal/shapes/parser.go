package shapes

import "errors"

// ParsePath validates and normalizes an SVG path string.
func ParsePath(path string) (string, error) {
	// TODO:
	// 1. Validate path syntax
	// 2. Normalize commands (e.g., relative to absolute)
	// 3. Return cleaned path
	if path == "" {
		return "", errors.New("empty path")
	}
	return path, nil
}

// IsValidPath checks if a string is a valid SVG path.
func IsValidPath(path string) bool {
	// TODO: Implement path validation
	// Check for valid SVG path commands: M, L, H, V, C, S, Q, T, A, Z
	return true
}

// ResolvePath returns a shape from either a registered name or raw SVG path.
func ResolvePath(nameOrPath string) (Shape, error) {
	// First check if it's a registered shape name
	if s := Get(nameOrPath); s != nil {
		return s, nil
	}

	// Otherwise treat as SVG path
	path, err := ParsePath(nameOrPath)
	if err != nil {
		return nil, err
	}

	return FromPath("custom", path), nil
}

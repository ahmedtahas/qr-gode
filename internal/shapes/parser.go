package shapes

import "errors"

// ParsePath validates and normalizes an SVG path string.
func ParsePath(path string) (string, error) {
	if path == "" {
		return "", errors.New("empty path")
	}
	// Basic validation - path should start with M or m (moveto)
	if len(path) > 0 && path[0] != 'M' && path[0] != 'm' {
		return "", errors.New("path must start with M or m command")
	}
	return path, nil
}

// IsValidPath checks if a string is a valid SVG path.
func IsValidPath(path string) bool {
	if path == "" || (path[0] != 'M' && path[0] != 'm') {
		return false
	}
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

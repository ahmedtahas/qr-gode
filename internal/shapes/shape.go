package shapes

// Shape defines how a module is rendered in SVG.
type Shape interface {
	// SVGPath returns the SVG path for a module at position (0,0) with size 1.
	// The path should be centered in the unit square.
	// It will be scaled and translated during rendering.
	SVGPath() string

	// Name returns the shape identifier.
	Name() string
}

// Registry holds registered shapes by name.
var Registry = make(map[string]Shape)

// Register adds a shape to the registry.
func Register(s Shape) {
	Registry[s.Name()] = s
}

// Get retrieves a shape by name, or nil if not found.
func Get(name string) Shape {
	return Registry[name]
}

// FromPath creates a custom shape from an SVG path string.
func FromPath(name, path string) Shape {
	return &customShape{name: name, path: path}
}

type customShape struct {
	name string
	path string
}

func (s *customShape) SVGPath() string { return s.path }
func (s *customShape) Name() string    { return s.name }

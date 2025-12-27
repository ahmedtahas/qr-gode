package qrcode

// Option is a functional option for configuring QR generation.
type Option func(*Config)

// WithErrorCorrection sets the error correction level.
func WithErrorCorrection(level ErrorCorrectionLevel) Option {
	return func(c *Config) {
		c.ErrorCorrection = level
	}
}

// WithSize sets the output size in pixels.
func WithSize(size int) Option {
	return func(c *Config) {
		c.Size = size
	}
}

// WithQuietZone sets the margin around the QR code.
func WithQuietZone(modules int) Option {
	return func(c *Config) {
		c.QuietZone = modules
	}
}

// WithModuleShape sets the shape for data modules.
func WithModuleShape(shape string) Option {
	return func(c *Config) {
		c.Modules.Shape = shape
	}
}

// WithLogo sets the logo path. Size is auto-calculated.
func WithLogo(path string) Option {
	return func(c *Config) {
		c.Logo = &LogoConfig{
			Path: path,
		}
	}
}

// WithLogoSize sets logo path with custom dimensions.
// Pass 0 for width or height to auto-calculate that dimension.
func WithLogoSize(path string, width, height int) Option {
	return func(c *Config) {
		c.Logo = &LogoConfig{
			Path:   path,
			Width:  width,
			Height: height,
		}
	}
}

// GenerateWithOptions creates a QR code using functional options.
func GenerateWithOptions(data string, opts ...Option) ([]byte, error) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return Generate(data, cfg)
}

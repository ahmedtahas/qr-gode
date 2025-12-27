# qr-gode

A feature-rich QR code generator library for Go with extensive customization options.

## Features

- Multiple module shapes (square, circle, rounded, diamond, dot, star, heart)
- Solid colors and gradients (linear & radial)
- Custom images for finder patterns, alignment patterns, and modules
- Logo support with automatic sizing and aspect ratio preservation
- SVG output with clean, optimized markup
- Configurable error correction levels

## Installation

```bash
go get github.com/ahmedtahas/qr-gode
```

## Quick Start

### CLI Usage

```bash
# Basic QR code
qr-gode "https://example.com"

# Custom styling
qr-gode -shape circle -fg "#3498db" "Hello World"

# Gradient colors
qr-gode -gradient "#ff6b6b,#4ecdc4" -shape rounded "Gradient QR"

# With a logo
qr-gode -logo logo.png "https://example.com"

# Custom images for patterns
qr-gode -finder-img finder.png -module-img dot.png "Custom QR"
```

### CLI Options

| Flag | Description | Default |
|------|-------------|---------|
| `-o` | Output file path | `qrcode.svg` |
| `-size` | Output size in pixels | `512` |
| `-shape` | Module shape | `square` |
| `-fg` | Foreground color (hex) | `#000000` |
| `-bg` | Background color (hex) | `#FFFFFF` |
| `-gradient` | Gradient colors (comma-separated) | - |
| `-gradient-angle` | Gradient angle in degrees | `45` |
| `-radial` | Use radial gradient | `false` |
| `-ecl` | Error correction level (L, M, Q, H) | `M` |
| `-logo` | Logo image path (PNG/JPG/SVG) | - |
| `-logo-width` | Logo width in pixels (0 = auto) | `0` |
| `-logo-height` | Logo height in pixels (0 = auto) | `0` |
| `-finder-img` | Custom finder pattern image | - |
| `-align-img` | Custom alignment pattern image | - |
| `-module-img` | Custom module image | - |

## Library Usage

### Builder API (Recommended)

The fluent builder API is the simplest way to generate QR codes:

```go
package main

import (
    "os"
    "github.com/ahmedtahas/qr-gode/pkg/qrcode"
)

func main() {
    // Simple QR code
    svg, err := qrcode.New("https://example.com").SVG()
    if err != nil {
        panic(err)
    }
    os.WriteFile("qr.svg", svg, 0644)
}
```

#### Styling Options

```go
// Custom colors and shape
svg, _ := qrcode.New("Hello World").
    Size(512).
    Shape("circle").
    Foreground("#3498db").
    Background("#ffffff").
    SVG()

// Linear gradient
svg, _ := qrcode.New("Gradient QR").
    Shape("rounded").
    LinearGradient(45, "#ff6b6b", "#4ecdc4").
    SVG()

// Radial gradient
svg, _ := qrcode.New("Radial QR").
    Shape("dot").
    RadialGradient(0.5, 0.5, "#ff6b6b", "#4ecdc4", "#45b7d1").
    SVG()
```

#### Adding a Logo

```go
// Auto-sized logo (15-30% of QR size, preserves aspect ratio)
svg, _ := qrcode.New("https://example.com").
    Logo("logo.png").
    SVG()

// Custom logo dimensions
svg, _ := qrcode.New("https://example.com").
    Logo("logo.png").
    LogoWidth(100).
    SVG()

// Transparent logo background
svg, _ := qrcode.New("https://example.com").
    Logo("logo.png").
    LogoBackground("transparent").
    SVG()
```

#### Custom Pattern Images

```go
// Use custom images for QR elements
svg, _ := qrcode.New("Custom QR").
    FinderImage("finder.png").
    AlignmentImage("alignment.png").
    ModuleImage("module.png").
    SVG()
```

#### Error Correction

```go
// Higher error correction for logos or damaged codes
svg, _ := qrcode.New("https://example.com").
    ErrorCorrection(qrcode.LevelH). // 30% recovery
    Logo("logo.png").
    SVG()
```

### Functional Options API

Alternative API using functional options:

```go
svg, err := qrcode.GenerateWithOptions("https://example.com",
    qrcode.WithSize(512),
    qrcode.WithShape("circle"),
    qrcode.WithLogo("logo.png"),
)
```

### Direct Config API

For full control, use the Config struct directly:

```go
cfg := qrcode.DefaultConfig()
cfg.Size = 512
cfg.Modules.Shape = "circle"
cfg.Modules.Color = qrcode.NewLinearGradientColor(45, []string{"#ff6b6b", "#4ecdc4"})
cfg.Logo = &qrcode.LogoConfig{
    Path: "logo.png",
}

svg, err := qrcode.Generate("https://example.com", cfg)
```

## Available Shapes

| Shape | Description |
|-------|-------------|
| `square` | Standard square modules (default) |
| `circle` | Circular modules |
| `rounded` | Rounded square modules |
| `diamond` | Diamond/rotated square modules |
| `dot` | Small centered dots |
| `star` | Star-shaped modules |
| `heart` | Heart-shaped modules |

## Error Correction Levels

| Level | Recovery Capacity | Use Case |
|-------|-------------------|----------|
| `L` | ~7% | Maximum data density |
| `M` | ~15% | Default, balanced |
| `Q` | ~25% | Better recovery |
| `H` | ~30% | Best recovery, use with logos |

When adding a logo, consider using `LevelQ` or `LevelH` to ensure the QR code remains scannable.

## Custom Images

### Finder Pattern Image

The finder pattern image should be a 7x7 module pattern:
- Outer ring (dark)
- White gap
- Inner 3x3 square (dark)

### Alignment Pattern Image

The alignment pattern image should be a 5x5 module pattern:
- Outer ring (dark)
- White gap
- Center dot (dark)

### Module Image

Any square image that represents a single dark module.

## Logo Handling

- **Auto-sizing**: By default, logos are scaled to fit within 15-30% of the QR code size
- **Aspect ratio**: Always preserved - logos are never stretched
- **SVG logos**: Treated as 1:1 aspect ratio, scale perfectly at any size
- **Background**: White rounded rectangle by default, can be set to transparent
- **Exclusion zone**: Modules under the logo area are not rendered (cleaner than overlay)

## Examples

### Generate to File

```go
err := qrcode.New("https://example.com").
    Shape("circle").
    LinearGradient(45, "#667eea", "#764ba2").
    Logo("logo.png").
    SaveSVG("output.svg")
```

### Generate to Bytes

```go
svg, err := qrcode.New("https://example.com").SVG()
// Use svg bytes directly (e.g., HTTP response)
```

### Validate Custom Images

```go
// Validate before using
if err := qrcode.ValidateFinderImage("finder.png"); err != nil {
    log.Fatal(err)
}

if err := qrcode.ValidateLogoImage("logo.png"); err != nil {
    log.Fatal(err)
}
```

## License

MIT License

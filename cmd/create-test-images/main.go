package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func main() {
	// Create blue circle for modules (single dot)
	createCircle("test_module.png", color.RGBA{52, 152, 219, 255}, 32)

	// Create proper finder pattern (7x7 structure)
	createFinderPattern("test_finder.png", color.RGBA{231, 76, 60, 255}, 70) // 70px = 10px per module

	// Create proper alignment pattern (5x5 structure)
	createAlignmentPattern("test_align.png", color.RGBA{39, 174, 96, 255}, 50) // 50px = 10px per module

	println("Created valid test images!")
	println("- test_module.png: 32x32 blue circle for data modules")
	println("- test_finder.png: 70x70 proper finder pattern (7x7 modules)")
	println("- test_align.png: 50x50 proper alignment pattern (5x5 modules)")
}

func createCircle(path string, c color.RGBA, size int) {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	center := float64(size) / 2
	radius := center - 2

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x) - center + 0.5
			dy := float64(y) - center + 0.5
			if math.Sqrt(dx*dx+dy*dy) <= radius {
				img.Set(x, y, c)
			}
		}
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

type darkFunc func(mx, my int) bool

func drawPattern(path string, c color.RGBA, size int, modules int, isDark darkFunc) {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	moduleSize := size / modules

	// Fill background with white
	white := color.RGBA{255, 255, 255, 255}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, white)
		}
	}

	// Draw the pattern
	for my := 0; my < modules; my++ {
		for mx := 0; mx < modules; mx++ {
			if isDark(mx, my) {
				// Fill this module
				startX := mx * moduleSize
				startY := my * moduleSize
				endX := (mx + 1) * moduleSize
				endY := (my + 1) * moduleSize

				for py := startY; py < endY; py++ {
					for px := startX; px < endX; px++ {
						if px < size && py < size {
							img.Set(px, py, c)
						}
					}
				}
			}
		}
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

// createFinderPattern creates a proper QR finder pattern (7x7 modules)
// Structure:
// ███████
// █     █
// █ ███ █
// █ ███ █
// █ ███ █
// █     █
// ███████
func createFinderPattern(path string, c color.RGBA, size int) {
	isDark := func(mx, my int) bool {
		// Outer ring
		if my == 0 || my == 6 || mx == 0 || mx == 6 {
			return true
		}
		// Inner 3x3 square
		if mx >= 2 && mx <= 4 && my >= 2 && my <= 4 {
			return true
		}
		return false
	}
	drawPattern(path, c, size, 7, isDark)
}

// createAlignmentPattern creates a proper QR alignment pattern (5x5 modules)
// Structure:
// █████
// █   █
// █ █ █
// █   █
// █████
func createAlignmentPattern(path string, c color.RGBA, size int) {
	isDark := func(mx, my int) bool {
		// Outer ring
		if my == 0 || my == 4 || mx == 0 || mx == 4 {
			return true
		}
		// Center dot
		if mx == 2 && my == 2 {
			return true
		}
		return false
	}
	drawPattern(path, c, size, 5, isDark)
}

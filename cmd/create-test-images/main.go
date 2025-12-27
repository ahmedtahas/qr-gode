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

	f, _ := os.Create(path)
	defer f.Close()
	png.Encode(f, img)
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
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	moduleSize := size / 7

	// Fill background with white first
	white := color.RGBA{255, 255, 255, 255}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, white)
		}
	}

	// Draw the finder pattern
	for my := 0; my < 7; my++ {
		for mx := 0; mx < 7; mx++ {
			// Determine if this module should be dark
			isDark := false

			// Outer ring (row 0, row 6, col 0, col 6)
			if my == 0 || my == 6 || mx == 0 || mx == 6 {
				isDark = true
			}
			// Inner 3x3 square (rows 2-4, cols 2-4)
			if mx >= 2 && mx <= 4 && my >= 2 && my <= 4 {
				isDark = true
			}

			if isDark {
				// Fill this module
				for py := my * moduleSize; py < (my+1)*moduleSize; py++ {
					for px := mx * moduleSize; px < (mx+1)*moduleSize; px++ {
						if px < size && py < size {
							img.Set(px, py, c)
						}
					}
				}
			}
		}
	}

	f, _ := os.Create(path)
	defer f.Close()
	png.Encode(f, img)
}

// createAlignmentPattern creates a proper QR alignment pattern (5x5 modules)
// Structure:
// █████
// █   █
// █ █ █
// █   █
// █████
func createAlignmentPattern(path string, c color.RGBA, size int) {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	moduleSize := size / 5

	// Fill background with white first
	white := color.RGBA{255, 255, 255, 255}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, white)
		}
	}

	// Draw the alignment pattern
	for my := 0; my < 5; my++ {
		for mx := 0; mx < 5; mx++ {
			// Determine if this module should be dark
			isDark := false

			// Outer ring (row 0, row 4, col 0, col 4)
			if my == 0 || my == 4 || mx == 0 || mx == 4 {
				isDark = true
			}
			// Center dot
			if mx == 2 && my == 2 {
				isDark = true
			}

			if isDark {
				// Fill this module
				for py := my * moduleSize; py < (my+1)*moduleSize; py++ {
					for px := mx * moduleSize; px < (mx+1)*moduleSize; px++ {
						if px < size && py < size {
							img.Set(px, py, c)
						}
					}
				}
			}
		}
	}

	f, _ := os.Create(path)
	defer f.Close()
	png.Encode(f, img)
}

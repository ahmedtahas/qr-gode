package main

import (
	"fmt"

	"github.com/ahmedtahas/qr-gode/internal/encoder"
)

func main() {
	fmt.Println("Alignment Pattern Positions by Version:")
	fmt.Println("========================================")

	for v := 1; v <= 20; v++ {
		version := encoder.Version(v)
		enc := encoder.New("test", encoder.LevelM)
		matrix, _ := enc.Encode()
		_ = matrix // just to initialize

		// Use the internal function by checking the matrix
		size := version.Size()

		// Count alignment patterns based on formula
		count := v/7 + 2
		if v == 1 {
			count = 0
		}

		total := count * count
		actual := total
		if count > 0 {
			actual = total - 3 // subtract 3 corners where finders are
		}
		if v == 1 {
			actual = 0
		}

		fmt.Printf("Version %2d (%2dx%2d): %d positions in grid -> %d alignment patterns\n",
			v, size, size, count, actual)
	}
}

package utils

import (
	"image/color"
	"math/rand"

	"github.com/tdewolff/canvas"
)

func DrawBackground(ctx *canvas.Context, color color.Color) {
	ctx.SetFillColor(color)
	bgRect := canvas.Rectangle(ctx.Width(), ctx.Height())
	ctx.DrawPath(0, 0, bgRect)
}

// generateRandomSequence generates a sequence of 'count' random numbers
// using the provided 'seed'.
func GenerateRandomSequence(seed int64, count int) []int64 {
	// Create a new pseudo-random number source with the given seed.
	// IMPORTANT: Do NOT use rand.Seed() here, as that sets the global seed,
	// affecting other parts of your program or concurrent goroutines.
	// Always create a new Source for reproducible sequences.
	source := rand.NewSource(seed)
	r := rand.New(source) // Create a new Rand using the source

	numbers := make([]int64, count)
	for i := 0; i < count; i++ {
		// Generate a random integer
		// For a more specific range (e.g., 0-100), use r.Intn(101)
		numbers[i] = r.Int63n(100000)
	}
	return numbers
}

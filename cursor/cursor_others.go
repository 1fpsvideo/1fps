//go:build !windows
// +build !windows

package cursor

import (
	"image"

	"github.com/go-vgo/robotgo"
)

type ResizedDimensions struct {
	Width  int
	Height int
}

func GetCursorPosition(resizedDimensions ResizedDimensions, bounds image.Rectangle) (int, int) {
	// Get the current cursor position
	x, y := robotgo.Location()

	// Calculate the cursor position relative to the bounds
	relativeX := x - bounds.Min.X
	relativeY := y - bounds.Min.Y

	// Ensure the cursor is within the bounds
	if relativeX < 0 {
		relativeX = 0
	} else if relativeX > bounds.Dx() {
		relativeX = bounds.Dx()
	}

	if relativeY < 0 {
		relativeY = 0
	} else if relativeY > bounds.Dy() {
		relativeY = bounds.Dy()
	}

	// Scale the cursor position to the resized dimensions
	scaledX := int(float64(relativeX) * float64(resizedDimensions.Width) / float64(bounds.Dx()))
	scaledY := int(float64(relativeY) * float64(resizedDimensions.Height) / float64(bounds.Dy()))

	return scaledX, scaledY
}

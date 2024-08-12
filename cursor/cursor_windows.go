//go:build windows
// +build windows

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
	// This position is in global coordinates
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

	// Get the actual screen size (considering the scale factor)
	// This gives us the true pixel dimensions of the screen
	actualWidth, actualHeight := robotgo.GetScaleSize()

	// Calculate the scaling ratio between the actual screen and our resized dimensions
	widthRatio := float64(resizedDimensions.Width) / float64(actualWidth)
	heightRatio := float64(resizedDimensions.Height) / float64(actualHeight)

	// Scale the relative cursor position to our resized dimensions
	scaledX := int(float64(relativeX) * widthRatio)
	scaledY := int(float64(relativeY) * heightRatio)

	// Ensure the scaled position is within the resized dimensions
	if scaledX >= resizedDimensions.Width {
		scaledX = resizedDimensions.Width - 1
	}
	if scaledY >= resizedDimensions.Height {
		scaledY = resizedDimensions.Height - 1
	}

	return scaledX, scaledY
}

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
	// Step 1: Get the current cursor position
	// This position might be scaled due to Windows DPI settings
	rawX, rawY := robotgo.Location()

	// Step 2: Get the actual screen size (considering the scale factor)
	// This gives us the true pixel dimensions of the screen
	actualWidth, _ := robotgo.GetScaleSize()

	// Step 3: Adjust cursor position according to the scale
	// This is necessary because Windows might report cursor positions beyond the actual screen size due to scaling
	scaleFactor := float64(actualWidth) / float64(bounds.Dx())
	adjustedX := int(float64(rawX) / scaleFactor)
	adjustedY := int(float64(rawY) / scaleFactor)

	// Step 4: Calculate the cursor position relative to the bounds
	// This handles multi-monitor setups where bounds might not start at (0,0)
	relativeX := adjustedX - bounds.Min.X
	relativeY := adjustedY - bounds.Min.Y

	// Step 5: Ensure the cursor is within the bounds
	// If it's outside, clamp it to the edge of the bounds
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

	// Step 6: Scale the cursor position to the resized dimensions
	// This maintains the relative position of the cursor in the resized view
	scaledX := int(float64(relativeX) * float64(resizedDimensions.Width) / float64(bounds.Dx()))
	scaledY := int(float64(relativeY) * float64(resizedDimensions.Height) / float64(bounds.Dy()))

	// Step 7: Final bounds check for the scaled position
	// Ensure the scaled position doesn't exceed the resized dimensions
	if scaledX >= resizedDimensions.Width {
		scaledX = resizedDimensions.Width - 1
	}
	if scaledY >= resizedDimensions.Height {
		scaledY = resizedDimensions.Height - 1
	}

	return scaledX, scaledY
}

//go:build windows
// +build windows

package cursor

import (
	// "fmt"

	"github.com/go-vgo/robotgo"
)

type ResizedDimensions struct {
	Width  int
	Height int
}

func GetCursorPosition(resizedDimensions ResizedDimensions) (int, int) {
	// Get the current cursor position
	// This position might be scaled due to Windows DPI settings
	x, y := robotgo.Location()
	// fmt.Printf("Raw cursor position: x=%d, y=%d\n", x, y)

	// Get the actual screen size (considering the scale factor)
	// This gives us the true pixel dimensions of the screen
	actualWidth, actualHeight := robotgo.GetScaleSize()
	// fmt.Printf("Actual screen size: width=%d, height=%d\n", actualWidth, actualHeight)

	// Calculate the scaling ratio between the actual screen and our resized dimensions
	widthRatio := float64(resizedDimensions.Width) / float64(actualWidth)
	heightRatio := float64(resizedDimensions.Height) / float64(actualHeight)
	// fmt.Printf("Scaling ratios: widthRatio=%f, heightRatio=%f\n", widthRatio, heightRatio)

	// Scale the cursor position to our resized dimensions
	scaledX := int(float64(x) * widthRatio)
	scaledY := int(float64(y) * heightRatio)
	// fmt.Printf("Scaled cursor position before bounds check: scaledX=%d, scaledY=%d\n", scaledX, scaledY)

	// Ensure the scaled position is within bounds
	if scaledX >= resizedDimensions.Width {
		// fmt.Printf("scaledX out of bounds, adjusting from %d to %d\n", scaledX, resizedDimensions.Width-1)
		scaledX = resizedDimensions.Width - 1
	}
	if scaledY >= resizedDimensions.Height {
		// fmt.Printf("scaledY out of bounds, adjusting from %d to %d\n", scaledY, resizedDimensions.Height-1)
		scaledY = resizedDimensions.Height - 1
	}

	// fmt.Printf("Final scaled cursor position: scaledX=%d, scaledY=%d\n", scaledX, scaledY)
	// fmt.Println("===")
	return scaledX, scaledY
}

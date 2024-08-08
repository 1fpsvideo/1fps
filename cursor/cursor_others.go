//go:build !windows
// +build !windows

package cursor

import "github.com/go-vgo/robotgo"

type ResizedDimensions struct {
	Width  int
	Height int
}

func GetCursorPosition(resizedDimensions ResizedDimensions) (int, int) {
	x, y := robotgo.Location()
	screenWidth, screenHeight := robotgo.GetScreenSize()
	scaledX := int(float64(x) * float64(resizedDimensions.Width) / float64(screenWidth))
	scaledY := int(float64(y) * float64(resizedDimensions.Height) / float64(screenHeight))
	return scaledX, scaledY
}

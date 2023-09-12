package canvassubjects

import (
	canvas "github.com/rocco-gossmann/go_wasmcanvas"
)

type Pixel struct {
	X, Y uint16

	Color canvas.Color

	Alpha uint8
}

// Implement CanvasSubject =====================================================
// ==============================================================================
func (p Pixel) Draw(w, h uint16, pixels *[]uint32) {
	if index, ok := canvas.IndexFromCoords(p.X, p.Y, w, h); ok {

		if p.Alpha == 0xff || p.Alpha == 0 {
			(*pixels)[index] = uint32(p.Color)

		} else {
			var factor float64 = float64(p.Alpha) / 255.0
			(*pixels)[index] = canvas.BlendPixel((*pixels)[index], uint32(p.Color), factor)
		}

	}
}

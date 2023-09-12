package canvassubjects

import canvas "github.com/rocco-gossmann/go_wasmcanvas"

type Fill struct {
	Color canvas.Color

	Alpha uint8
}

// Implement CanvasSubject =====================================================
// ==============================================================================
func (f Fill) Draw(w, h uint16, pixels *[]uint32) {

	if f.Alpha == 0xff || f.Alpha == 0 {
		for index := 0; index < len(*pixels); index++ {
			(*pixels)[index] = uint32(f.Color)
		}
	} else {
		var factor float64 = float64(f.Alpha) / 255.0

		for index := 0; index < len(*pixels); index++ {
			(*pixels)[index] = canvas.BlendPixel((*pixels)[index], uint32(f.Color), factor)
		}
	}

}

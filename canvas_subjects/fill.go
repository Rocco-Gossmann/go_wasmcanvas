package canvassubjects

import canvas "github.com/rocco-gossmann/go_wasmcanvas"

type Fill struct {
	Color canvas.Color

	Alpha uint8
}

// Implement CanvasSubject =====================================================
// ==============================================================================
func (f Fill) Draw(c *canvas.Canvas) {

	px := c.GetPixelIndex(0)
	nxt := uint32(1)

	if f.Alpha == 0xff || f.Alpha == 0 {
		for px != nil {
			*px = uint32(f.Color)
			px = c.GetPixelIndex(nxt)
			nxt++
		}
	} else {
		var factor float64 = float64(f.Alpha) / 255.0

		for px != nil {
			*px = canvas.BlendPixel(*px, uint32(f.Color), factor)
			px = c.GetPixelIndex(nxt)
			nxt++
		}
	}

}

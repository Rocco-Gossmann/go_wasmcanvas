package canvasfragments

import (
	Ca "github.com/rocco-gossmann/go_wasmcanvas"
)

type Fill struct {
	Color Ca.Color
	Alpha byte
}

func (f *Fill) Draw(_ uint32, _, _ uint16, pixels *[]uint32) {
	if f.Alpha == 0 || f.Alpha == 0xff {
		for i, _ := range *pixels {
			(*pixels)[i] = uint32(f.Color)
		}

	} else {
		var factor float64 = (float64(f.Alpha) / 255.0)
		for i, _ := range *pixels {
			Ca.BlendPixel(&(*pixels)[i], uint32(f.Color), factor)
		}

	}

}

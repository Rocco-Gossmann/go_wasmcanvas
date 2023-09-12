package go_wasmcanvas

import (
	"math"
)

func roundBlend(v float64) float64 {
	if v > 0 {
		return math.Ceil(v)
	} else {
		return math.Floor(v)
	}
}

func IndexFromCoords(x, y, w, h uint16) (uint32, bool) {
	if x < 0 || x > w-1 || y < 0 || y > h-1 {
		return 0, false
	}
	return uint32(y*w + x), true
}

func BlendPixel(existingPixel uint32, newPixel uint32, factor float64) (ret uint32) {

	er := float64(existingPixel & (255 << 16) >> 16)
	eg := float64(existingPixel & (255 << 8) >> 8)
	eb := float64(existingPixel & (255))

	ret = (uint32(er-roundBlend(er-float64(newPixel&(255<<16)>>16))*factor) << 16) +
		(uint32(eg-roundBlend(eg-float64(newPixel&(255<<8)>>8))*factor) << 8) +
		uint32(eb-roundBlend(eb-float64(newPixel&(255)))*factor)

	return
}

package go_wasmcanvas

import (
	"math"
)

func ExtractRGB(px Color) (r float64, g float64, b float64) {
	r = float64(px & (255 << 16) >> 16)
	g = float64(px & (255 << 8) >> 8)
	b = float64(px & (255))

	return
}

func CombineRGB(r, g, b byte) Color {
	return Color((uint32(r) << 16) + (uint32(g) << 8) + uint32(b))
}

func BlendPixel(existingPixel *uint32, newPixel uint32, factor float64) {

	nr, ng, nb := ExtractRGB(Color(newPixel))
	er, eg, eb := ExtractRGB(Color(*existingPixel))

	*existingPixel = ((*existingPixel) & (0xff000000)) +
		(uint32(er-roundBlend((er-nr)*factor)) << 16) +
		(uint32(eg-roundBlend((eg-ng)*factor)) << 8) +
		uint32(eb-roundBlend((eb-nb)*factor))

}

// Private Helpers
//==============================================================================

func roundBlend(v float64) float64 {
	if v > 0 {
		return math.Ceil(v)
	} else {
		return math.Floor(v)
	}
}

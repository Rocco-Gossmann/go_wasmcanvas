package canvassubjects

import (
	"fmt"
	"math"

	canvas "github.com/rocco-gossmann/go_wasmcanvas"
)

type Line struct {
	Startx, Starty uint16
	Endx, Endy     uint16

	Color canvas.Color

	Alpha uint8
}

// =============================================================================
// Implement CanvasSubject
// =============================================================================
func (l Line) Draw(w, h uint16, pixels *[]uint32) {

	x1, x2, y1, y2 := l.Startx, l.Endx, l.Starty, l.Endy
	i1, okstart := canvas.IndexFromCoords(x1, y1, w, h)
	i2, okend := canvas.IndexFromCoords(x2, y2, w, h)

	if !okstart {
		panic("invalid start coordinates")
	}
	if !okend {
		panic("invalid end coordinates")
	}

	fmt.Println(x1, y1, "(", i1, ") to ", x2, y2, "(", i2, ")")

	var drawFragment func(x, y uint16)
	var factor = float64(l.Alpha) / 255.0

	if l.Alpha == 0x0 || l.Alpha == 0xff {
		drawFragment = func(x, y uint16) {
			(*pixels)[y*w+x] = uint32(l.Color)
		}
	} else {
		drawFragment = func(x, y uint16) {
			var i = y*w + x
			(*pixels)[i] = canvas.BlendPixel((*pixels)[i], uint32(l.Color), factor)
		}
	}

	if x1 == x2 {
		// Vertical Lines
		//--------------------------------------------------------------------------
		if y1 > y2 {
			y2, y1 = y1, y2
		}

		for y := y1; y <= y2; y++ {
			drawFragment(x1, y)
		}

	} else if y1 == y2 {
		// Horizontal Line
		//--------------------------------------------------------------------------
		if x1 > x2 {
			x2, x1 = x1, x2
		}

		for x := x1; x <= x2; x++ {
			drawFragment(x, y1)
		}

	} else {
		// Diagonal line
		//--------------------------------------------------------------------------
		aspect := math.Abs(float64(int(x1)-int(x2))) / math.Abs(float64(int(y1)-int(y2)))
		if aspect <= 1.0 {
			// Y-Dominat
			//--------------------------------------------------------------------------
			ystart, yend, x, xstep := prepLineVars(y1, y2, x1, x2, aspect)

			for y := ystart; y <= yend; y++ {
				drawFragment(uint16(math.Floor(x)), y)
				x += xstep
			}

		} else if aspect > 1.0 {
			// X-Dominat
			//--------------------------------------------------------------------------
			xstart, xend, y, ystep := prepLineVars(x1, x2, y1, y2, 1/aspect)

			for x := xstart; x <= xend; x++ {
				drawFragment(x, uint16(math.Floor(y)))
				y += ystep
			}

		}
	}

}

// =============================================================================
// Private Helpers
// =============================================================================
func prepLineVars(d1 uint16, d2 uint16, s1 uint16, s2 uint16, as float64) (uint16, uint16, float64, float64) {

	ds := min(d1, d2)
	de := max(d1, d2)
	ss := int(s1)
	se := int(s2)

	dd := de - ds

	if ds == d2 {
		ss, se = se, ss
	}

	sd := se - ss
	st := float64(sd) / float64(dd)

	return ds, de, float64(ss), st
}
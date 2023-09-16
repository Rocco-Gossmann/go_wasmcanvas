package canvasfragments

import (
	"math"

	"github.com/rocco-gossmann/go_throwable"
	Canvas "github.com/rocco-gossmann/go_wasmcanvas"
)

type Line struct {
	Startx, Starty uint16
	Endx, Endy     uint16

	Points []Canvas.Point

	Color Canvas.Color

	Alpha uint8
}

// =============================================================================
// Implement CanvasSubject
// =============================================================================
func (l Line) Draw(_ uint32, w, h uint16, pixels *[]uint32) {

	noPoints := false
	if len(l.Points) == 0 {
		l.Points = append(l.Points, Canvas.Point{l.Startx, l.Starty})
		l.Points = append(l.Points, Canvas.Point{l.Endx, l.Endy})
		noPoints = true
	}

	if len(l.Points) < 2 {
		go_throwable.Throw(Canvas.CanvasPanic{
			Msg:     "lines needs more points",
			Subject: "Line",
			Value:   len(l.Points),
			Allowed: " min. 2 Points",
		})
	}

	var drawFragment func(x, y uint16)
	var factor = float64(l.Alpha) / 255.0

	if l.Alpha == 0x0 || l.Alpha == 0xff {
		drawFragment = func(x, y uint16) {
			if xyOK(x, y, w, h) {
				(*pixels)[y*w+x] = uint32(l.Color)
			}
		}
	} else {
		drawFragment = func(x, y uint16) {
			if xyOK(x, y, w, h) {
				Canvas.BlendPixel(&((*pixels)[y*w+x]), uint32(l.Color), factor)
			}
		}
	}

	index := 2
	lastPoint := &(l.Points[0])
	nextPoint := &(l.Points[1])

	for nextPoint != nil {

		x1, y1, x2, y2 := lastPoint.X, lastPoint.Y, nextPoint.X, nextPoint.Y

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
					drawFragment(uint16(math.Round(x)), y)
					x += xstep
				}

			} else if aspect > 1.0 {
				// X-Dominat
				//--------------------------------------------------------------------------
				xstart, xend, y, ystep := prepLineVars(x1, x2, y1, y2, 1/aspect)

				for x := xstart; x <= xend; x++ {
					drawFragment(x, uint16(math.Round(y)))
					y += ystep
				}
			}
		}

		lastPoint = nextPoint
		if index < len(l.Points) {
			nextPoint = &(l.Points[index])
			index++
		} else {
			nextPoint = nil
		}

	}

	if noPoints {
		l.Points = l.Points[0:0]
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

func xyOK(x, y, w, h uint16) bool {
	return (x >= 0 && x < w && y >= 0 && y < h)
}

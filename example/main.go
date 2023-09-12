package main

import (
	"fmt"

	gt "github.com/rocco-gossmann/go_throwable"
	canvas "github.com/rocco-gossmann/go_wasmcanvas"
	cs "github.com/rocco-gossmann/go_wasmcanvas/canvas_subjects"
)

func tick(c *canvas.Canvas, deltaTime float64) canvas.CanvasTickFunction {

	ret, _ := interface{}(gt.Try(func() any {

		// Should slowly blend in a Green pixel at coords 10 x 10
		//	c.Draw(cs.Pixel{X: 10, Y: 10, Color: canvas.COLOR_GREEN, Alpha: 4})

		c.Draw(cs.Pixel{X: 0, Y: 0, Color: canvas.COLOR_WHITE})
		c.Draw(cs.Pixel{X: 319, Y: 0, Color: canvas.COLOR_WHITE})
		c.Draw(cs.Pixel{X: 0, Y: 199, Color: canvas.COLOR_WHITE})
		c.Draw(cs.Pixel{X: 319, Y: 199, Color: canvas.COLOR_WHITE})

		c.Draw(cs.Line{
			Startx: 0, Starty: 1,
			Endx: 319, Endy: 200,
			Color: canvas.COLOR_ORANGE,
		})

		return tick

	}, gt.TryOpts{
		Catch: func(p any) any {
			fmt.Println("cought panic => shutdown canvas", p)
			return nil
		},
	})).(canvas.CanvasTickFunction)

	return ret
}

func main() {
	canv := canvas.Create(320, 200)

	// Initial background fill
	canv.Draw(cs.Fill{Color: canvas.COLOR_DARKGRAY})

	canv.Run(tick)
}

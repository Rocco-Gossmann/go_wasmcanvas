package main

import Ca "github.com/rocco-gossmann/go_wasmcanvas"

var pixelX float64 = 0     //<- hold the pixels position
const duration float64 = 5 //<- move the pixel in 5 seconds

func tick(c *Ca.Canvas, deltaTime float64) Ca.CanvasTickFunction {

	var pxPerSec = float64(c.Width()) / duration
	pixelX += pxPerSec * deltaTime
	c.SetPixel(uint16(pixelX), 100, Ca.COLOR_GREEN)

	return tick
}

func main() {
	canv := Ca.Create(320, 200)

	canv.Run(tick)
}

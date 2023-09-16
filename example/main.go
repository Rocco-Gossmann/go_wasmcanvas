package main

import (
	Ca "github.com/rocco-gossmann/go_wasmcanvas"
	Cf "github.com/rocco-gossmann/go_wasmcanvas/canvas_fragments"
)

var pixelX float64 = 0     //<- hold the pixels position
const duration float64 = 5 //<- move the pixel in 5 seconds

var fill = Cf.Fill{Color: Ca.COLOR_DARKGRAY}

var player = Cf.Line{
	Startx: 0, Starty: 100,
	Endx: 0, Endy: 100,
	Color: Ca.COLOR_LIGHTGREEN,
}

func tick(c *Ca.Canvas, deltaTime float64) Ca.CanvasTickFunction {
	var pxPerSec = float64(c.Width()) / duration

	c.Draw(&fill)

	player.Startx = player.Endx

	player.Endx = uint16(float64(player.Endx) + pxPerSec*deltaTime)

	if player.Startx > c.Width() {
		player.Endx -= player.Startx
		player.Startx = 0
	}

	c.Draw(&player)

	return tick
}

func main() {
	canv := Ca.Create(320, 200)

	canv.Draw(&fill)
	fill.Alpha = 16

	canv.Run(tick)
}

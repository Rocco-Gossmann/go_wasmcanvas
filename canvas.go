package go_wasmcanvas

import (
	ex "github.com/rocco-gossmann/go_throwable"
)

const max_dimension = 10000

type canvas struct {
	width, height uint16

	pixelCount uint32
	byteSize   uint32
	pixels     []uint32
}

func (c *canvas) Width() uint16  { return c.width }
func (c *canvas) Height() uint16 { return c.height }

func Create(width, height uint16) canvas {

	if width > max_dimension {
		ex.Throw(&canvasPanic{"width can't be more than " + string(max_dimension) + "pixels"})
	}
	if height > max_dimension {
		ex.Throw(&canvasPanic{"height can't be more than " + string(max_dimension) + "pixels"})
	}

	var canv canvas
	canv.width = width
	canv.height = height

	canv.pixelCount = uint32(width) * uint32(height)
	canv.byteSize = canv.pixelCount * 4

	canv.pixels = make([]uint32, canv.pixelCount)

	return canv
}

package go_wasmcanvas

type CanvasSubject interface {
	Draw(canvaswidth, canvasheight uint16, pixels *[]uint32)
}

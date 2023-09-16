package go_wasmcanvas

type CanvasFragmentFunction func(pixelCnt uint32, pixelPerRow uint16, rowCount uint16, pixels *[]uint32)

type CanvasFragment interface {
	Draw(pixelCnt uint32, pixelPerRow uint16, rowCount uint16, pixels *[]uint32)
}

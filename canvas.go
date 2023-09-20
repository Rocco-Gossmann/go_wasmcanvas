package go_wasmcanvas

import (
	"syscall/js"
	"time"

	ex "github.com/rocco-gossmann/go_throwable"
)

// Types =======================================================================
// =============================================================================
type CanvasTickFunction func(c *Canvas, deltaTime float64) CanvasTickFunction
type Canvas struct {
	id            int
	width, height uint16
	maxx, maxy    uint16

	pixelCount uint32
	byteSize   uint32
	pixels     []uint32
}

// Getters =====================================================================
// =============================================================================
func (c *Canvas) Width() uint16    { return c.width }
func (c *Canvas) Height() uint16   { return c.height }
func (c *Canvas) PixelCnt() uint32 { return c.pixelCount }

// Browser <=> Go Communication ================================================
// =============================================================================
var vblankchannel chan byte = make(chan byte)

func onMessage(this js.Value, args []js.Value) interface{} {
	ev := args[0].Get("data")
	if ev.Type() != js.Undefined().Type() {
		switch ev.Get("0").String() {
		case "vblankdone":
			vblankchannel <- 1
		}
	}

	return nil
}

// Constructor =================================================================
// =============================================================================
const max_dimension = 10000

var canvasid_autoincrement = 0

func Create(width, height uint16) Canvas {

	canvasid_autoincrement++

	if width > max_dimension {
		ex.Throw(&CanvasPanic{"invalid", "width", width, max_dimension})
	}
	if height > max_dimension {
		ex.Throw(&CanvasPanic{"invalid", "height", height, max_dimension})
	}

	var c Canvas
	c.id = canvasid_autoincrement
	c.width = width
	c.height = height
	c.maxx = width - 1
	c.maxy = height - 1

	c.pixelCount = uint32(width) * uint32(height)
	c.byteSize = c.pixelCount * 4
	c.pixels = make([]uint32, c.pixelCount)

	js.Global().Call("postMessage", []interface{}{"createCanvas", c.id, width, height})
	js.Global().Call("addEventListener", "message", js.FuncOf(onMessage), false)

	return c
}

func (c Canvas) GetPixel(x, y uint16) *uint32 {
	if x < 0 || x > c.maxx || y < 0 || y > c.maxy {
		return nil
	}
	return &(c.pixels[uint32(y*c.width+x)])
}

func (c Canvas) GetPixelIndex(index uint32) (pixel *uint32) {
	if index >= 0 && index < c.pixelCount {
		return &(c.pixels[index])
	} else {
		return nil
	}
}

func (c Canvas) SetPixel(x, y uint16, color Color) bool {
	if x < 0 || x > c.maxx || y < 0 || y > c.maxy {
		return false
	}

	c.pixels[uint32(y*c.width+x)] = uint32(color)
	return true
}

func (c Canvas) Draw(frag CanvasFragment) {
	frag.Draw(c.pixelCount, c.width, c.height, &(c.pixels))
}

func (c Canvas) Apply(fnc CanvasFragmentFunction) {
	fnc(c.pixelCount, c.width, c.height, &c.pixels)
}

func (c Canvas) Run(tick CanvasTickFunction) {
	fnc := tick
	var last int64 = 9223372036854775807 // <- time starts out with the higest
	//    									   value possible to skip the
	//    									   initial, inacurate update
	//    									   in a way, that takes advantage
	//    									   of the time rollover mechanic

	for fnc != nil {
		//-----------------------------------------------------------------------
		// Transfere pixels to JS Space
		//-----------------------------------------------------------------------
		// convert  from [pixelCnt]uint32  to  [byteSize]byte first
		buf := make([]byte, c.byteSize)
		var ptr = 0
		for _, val := range c.pixels {
			buf[ptr] = byte((val & (255 << 16)) >> 16) // r
			buf[ptr+1] = byte((val & (255 << 8)) >> 8) // g
			buf[ptr+2] = byte(val & 255)               // b
			buf[ptr+3] = 0xff                          // a
			ptr += 4
		}

		arr := js.Global().Get("Uint8ClampedArray").New(js.ValueOf(c.byteSize))
		js.CopyBytesToJS(arr, buf)

		//-----------------------------------------------------------------------
		// Request the Browser to redraw and wait until it is done
		//-----------------------------------------------------------------------
		js.Global().Call(
			"postMessage",                      // command
			[]interface{}{"vblank", c.id, arr}, // data
			[]interface{}{arr.Get("buffer")},   // worker transfere
		)

		//-----------------------------------------------------------------------
		// Update the Canvas
		//-----------------------------------------------------------------------
		var now = time.Now().UnixMilli()
		delta := max(0, now-last)
		if delta > 0 {
			var deltaTime float64 = float64(now-last) / 1000.0
			fnc = fnc(&c, deltaTime)
			last = now
		}
		last = now
		<-vblankchannel
	}

	js.Global().Call("removeEventListener", "message", js.FuncOf(onMessage), false)
	js.Global().Call("postMessage", []interface{}{"destroyCanvas", c.id})
}

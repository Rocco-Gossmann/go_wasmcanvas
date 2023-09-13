package go_wasmcanvas

import (
	"fmt"
	"syscall/js"
	"time"

	ex "github.com/rocco-gossmann/go_throwable"
)

// Types =======================================================================
// =============================================================================
type CanvasTickFunction func(c *Canvas, deltaTime float64) CanvasTickFunction
type Canvas struct {
	width, height uint16

	pixelCount uint32
	byteSize   uint32
	pixels     []uint32

	Draw func(CanvasSubject)
	Run  func(CanvasTickFunction)
}

// Getters =====================================================================
// =============================================================================
func (c *Canvas) Width() uint16  { return c.width }
func (c *Canvas) Height() uint16 { return c.height }

// Browser <=> Go Communication ================================================
// =============================================================================
var vblankchannel chan byte = make(chan byte)

func onMessage(this js.Value, args []js.Value) interface{} {

	ev := args[0].Get("data")
	if ev.Type() != js.Undefined().Type() {
		switch ev.Get("0").String() {
		case "vblankdone":
			vblankchannel <- 1
		default:
			fmt.Println("cant handle message", ev.Get("0"))
		}
	}

	return nil
}

// Constructor =================================================================
// =============================================================================
const max_dimension = 10000

func Create(width, height uint16) Canvas {

	if width > max_dimension {
		ex.Throw(&CanvasPanic{"invalid", "width", width, max_dimension})
	}
	if height > max_dimension {
		ex.Throw(&CanvasPanic{"invalid", "height", height, max_dimension})
	}

	var c Canvas
	c.width = width
	c.height = height

	c.pixelCount = uint32(width) * uint32(height)
	c.byteSize = c.pixelCount * 4

	c.pixels = make([]uint32, c.pixelCount)

	c.Draw = func(s CanvasSubject) { s.Draw(c.width, c.height, &c.pixels) }

	c.Run = func(tick CanvasTickFunction) {
		fnc := tick
		var last int64 = 9223372036854775807 // <- time starts out with the higest
		//    value possible to skip the
		//    initial, inacurate update
		//    in a way, that takes advantage
		//    of the time rollover mechanic

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
				"postMessage",                    // command
				[]interface{}{"vblank", arr},     // data
				[]interface{}{arr.Get("buffer")}, // worker transfere
			)

			<-vblankchannel

			//-----------------------------------------------------------------------
			// Update the Canvas
			//-----------------------------------------------------------------------
			var now = time.Now().UnixMilli()

			if now > last {
				var deltaTime float64 = float64(last-now) / 1000.0
				fnc = fnc(&c, deltaTime)
				continue
			}

			last = now
		}
	}

	js.Global().Call("postMessage", []interface{}{"createCanvas", width, height})
	js.Global().Call("addEventListener", "message", js.FuncOf(onMessage), false)

	return c
}

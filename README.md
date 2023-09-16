# Go-WASMCanvas

Provides everything nessary to control a HTML-Canvas Pixel by Pixel via Go.

check out a demo here: https://rocco-gossmann.github.io/go_wasmcanvas/

## HTML-Preparation:
This is only usable in WebBrowsers, therefore some preparations need to be made.

### 1. First you need the __wasm_exec.js__ from `$GOROOT/misc/wasm`  
In Bash, type 
```bash 
cd `go env GOROOT`/misc/wasm
```
To Find it.

---
### 2. The WebAssembly must be run from a WebWorker, as to not impact the UI

This is what it should look like 
```javascript
// Loading the script copied in step 1
importScripts("./wasm_exec.js");

// Loading and starting the WebAssembly
const go = new Go();
WebAssembly.instantiateStreaming(fetch("./main.wasm"), go.importObject)
    .then( gowasm => {
        if(!gowasm) {
            console.error("failed to instantiate gowasm", go, gowasm);
            alert("technical error");
            return;
        }
        self.postMessage("wasm has started")
        return go.run(gowasm.instance)
    })      
    .then( () => self.postMessage("wasm has ended"))
```

---

### 3. Communication between UI (Main-Thread) and WebWorker

The WASM - Machine inside the WebWorker will send various Post-Messages to the Main-Thread/UI.

It is the UI's Job to react to them and respond with its own post messages.

If you just want something to get going.  
Here you go:
```javascript
if (WebAssembly) {

    const worker = new Worker("./worker.js");

    const canvases = new Map();
    const contexts = new Map();

    worker.addEventListener("message", (ev) => {
        if (ev.data instanceof Array) {
            const canvasid = ev.data[1];
            const canv = canvases.get(canvasid)
            const ctx = contexts.get(canvasid);

            switch (ev.data[0]) {

                case "createCanvas": {
                    const canv = document.createElement("canvas");
                    canv.width = ev.data[2];
                    canv.height = ev.data[3];
                    canv.className="go-wasm-canvas";

                    canvases.set(canvasid, canv);
                    contexts.set(canvasid, canv.getContext("2d"));

                    document.body.appendChild(canv);
                    reactToScreenSize(canv)
                } break;


                case "vblank": {
                    const imgDat = new ImageData(ev.data[2], canv.width, canv.height);
                    window.requestAnimationFrame( () => { 
                        ctx.putImageData(imgDat, 0, 0);
                        worker.postMessage(["vblankdone"]) 
                    }); 
                } break;


                case "destroyCanvas": {
                    canv.parentNode.removeChild(canv);
                    canvases.delete(canvasid);
                    contexts.delete(canvasid);
                } break;

                default:{
                    console.log( "[WORKER MESSAGE]", ev.data[0], " in ", ev);
                }
            }
        }
        else console.log("unknwon worker message", ev.data);
    }, {capture: true})

}
else alert("Your Browser does not support WebAssembly");


```
Now for a more detailed explaination:

Each event send by or to Go-WASMCanvas will be structured like follows:
```typescript
[ task: string, canvas.id: number,  ...args: any]
```
and Array whos first element always determins, how the other elements
should be interpreted. 
The Second element is always the canvas, that issued the postMessage 
(You can create multiple canvases via Go)

These are the possible Events/Messages that are send from Go right now.

| Event |  Input   |  Output | Description |
| ------|----------|---------|-------------|
| createCanvas |  `width: number`, `height: number` |  none  |  The Dimensions of the canvas are determined by Go, therefore it will instruct the Browser on how to construct its canvas. <br><br>|
| vblank |  none  | `[ 'vblankdone' ]`| At the beginning of each Tick, Go will send the current PixelBuffer to the UI and then wait until it receives a 'vblankdone' message.<br ><br>When receiving a `vblank` message, the MainThread should create a new `ImageData`-Instance using the given pixels.<br><br>The UI should use `requestAnimationFrame` to wait until it is save to change the canvas pixels.<br><br> Once a new frame is availble, it should put the `ImageData` on the canvas that was created during `createCanvas`.<br><br> Ending the process by posting a `['vblankdone']` message back to the WebWorker running GoWasm. 




## Usage in Go:

### 1.) Import the Package
```go
import Canvas "github.com/rocco-gossmann/go_wasmcanvas"
```

### 2.) Create a Canvas in your `main` function.
```go
Canvas.Create(width uint16, height uint16) Canvas.Canvas
```

The parameters are 
| | | |
|-|-|-|
| `width`  | `uint16` | Pixel width of the canvas. Max: 10000px|
| `height` | `uint16` | Pixel height of th canvas. Max: 10000px|


```go
func main() {
    ca := Canvas.Create(320, 200) //<- creates 320 by 200 px canvas
    //...
}
```
As descripted in the __HTML-Preparations__ The module will then talk to the Main-Threads/UIs Javascript via `postMessage` .

So it has no control over what is shown on screen. Rather it tells the Browser, what it wishes to happen. In a sence the MainThread acts more as a virtual Graphics and IO processor, while Go acts as a data processor.
(Similar to how a 6502 Processor would interact with the PPU-Chip on a NES console or the VIC-Chip on a C64)

### 3.) Provide a `tick` function
```go
// Doc
type CanvasTickFunction func(c *Canvas, deltaTime float64) CanvasTickFunction
```
A `TickFunction` is a function that gets executed every vblank cycle.

#### Return: 
| | | 
|-|-|
|`TickFunction`| It nedds to Return a pointer to the function that will run next tick.<br>If `nil` is returned, the Canvas will shutdown and end its Execution |

#### Params:
The TickFunction takes 2 parameters.

|   |  | |
|-|-|-|
| `canvas` | `*Canvas.Canvas` | a Pointer to the `Canvas` - insetance that called it |
| `deltaTime` |`float64` | provieds a percentage value of how many fractions of a second have passed since the last tick execution.

#### More Details on deltaTime
```
example:
1.2   == 1200ms == 1 Second and 2ms
0.5   ==  500ms == half a second
0.016 ==   16ms == 1/60th of a second.
```
#### Lets give an Example:

Here we want a Pixel to cross a canvas in 5 Seconds (Regardless of Canvas size or Browser Performance).
The browser performance (and thus execution times) can differ wildly, but thanks to the `deltaTime`, we can still reach our goal of always having the pixel cross the canvas in 5 seconds.

In return the fluidity of the animation depends on how powerfull the browser is
```go
var pixelX float64 = 0      //<- hold the pixels position
const duration float64 = 5 //<- move the pixel across the canvas in 5 seconds

func tick(c *Canvas.Canvas, deltaTime float64) Canvas.CanvasTickFunction {

	var pxPerSec = float64(c.Width()) / duration  // <- define how many pixels we must move 
    //                                                  per seconds to cross the canvas in the time we need

	pixelX += pxPerSec * deltaTime  //<- use Delta Time to define how much we must move this tick. 

	c.SetPixel(uint16(pixelX), 100, Canvas.COLOR_GREEN) //<- Set the pixel

	return tick
}
```



### 4. Telling the Canvas to run the Tick-Function

```go
Canvas.Run( tick TickFunction)
```

This function must be called on the canvas instance.
```go
func main() {
    ca := Canvas.Create( 320, 200)
    ca.Run(tick)
}
```


So lets bring the last 4 points all together 

```go
package main

import Canvas "github.com/rocco-gossmann/go_wasmcanvas"

var pixelX float64 = 0      //<- hold the pixels position
const duration float64 = 5 //<- move the pixel in 5 seconds

func tick(c *Canvas.Canvas, deltaTime float64) Canvas.CanvasTickFunction {

	var pxPerSec = float64(c.Width()) / duration
	pixelX += pxPerSec * deltaTime
	c.SetPixel(uint16(pixelX), 100, Canvas.COLOR_GREEN)

	return tick
}

func main() {
	ca := Canvas.Create(320, 200)
	ca.Run(tick)
}
```

### 5. Important things to note
The example above should result in a green line being drawn from left to right within the span of 5 seconds.

Depending on your Browser, you will see little gaps within the line. These stam from the Browsers internal functions taking up more or less time.
On some browsers, these gaps are more consistant than other.

Unfortunatly due to the nature of Javascript, we can't do much about these.  I'll try my gest to mitigate them within the Go-Packages functions,
But they will never go await to 100%.

But with some clever Application design, we can at least hide them.  (See the "Advanced" Drawing section further down);

## Colors

Go-WASMCanvas uses `24Bit` colors.
Each pixel of the Canvas is represented by a `uint32`
```
# 00 rr gg bb
   |  |  |  |
   |  |  |  --- 8Bit Blue channel
   |  |  -------8bit Green channel  
   |  ----------8bit Red chanel
   -------------unused memory 
```
The first 8bit of each pixel are not transfered to the canvas, and are only visible to Go. 
They can for example be used by a method from the  __Advanced Drawing__ Section 
to give some state to each pixel.

Go-WASMCanvas comes with a set of predefined colors. These are made to 
emulate the Color-Palette of a C64

You are however free do define your own pixel colors, as the Canvas.Color - Type is base on the `uint32` type 
```go
type Canvas.Color uint32
```

## Advanced Drawing

Drawing everything Pixel by pixel via the Canvas.SetPixel Method would be highly inefitient though.

Thats why Go-Canvas supports two more methods to draw things


### `Canvas.Apply()`
```go
ca.Apply(func(pixelCount uint32, canvasWidth, canvasHeight uint16, pixels *[]uint32))
```
The `Apply` method takes in a function the get access to the essentialy information needed to draw pixels on the canvas. including `pixels`-Memory.  

It is best used for things that apply to to the entire Canvas. (Like Filters, etc.)


### `Canvas.Draw()`
```go
ca.Draw(fragment Canvas.CanvasFragment) 
```
the `Draw` method is a bit more sophisticated than the `Apply` method.
Allowing you to define what to draw via the definition of a `struct`.

The Only requirement is, that this `struct` must implement the `Canvas.CanvasFragment` - `interface`

The Interface is defined as follows.
```go
type CanvasFragment interface {
	Draw(pixelCnt uint32, pixelPerRow uint16, rowCount uint16, pixels *[]uint32)
}
```

### Defing A canvas Fragment
As you can see the All the Interface requires is a `Draw` method, that is provided the vital fields that are required to manipulate each pixel of the canvas.

`*pixels` grants direct access to the entire pixelbuffer of the canvas.

In a sence the Fragments act more like shaders, as the do conventional drawing functions.

Given the direct pixel access, you can come up with all kinds of things to draw this way.

As an Example, This package includes 2 predefined CanvasFragments

### Included Canvas-Fragments.
This package comes with two available `CanvasFragments` out of the box.

```go
// Fills the canvas with the given color
canvasfragment.Fill{
    Color Canvas.Color
    Alpha byte
}

```

```go
// Draws a line between 2 or more points
canvasfragment.Line{
    Startx, Starty, Endx, Endy uint16, //<- can be ignored if Points are defined

	Points []Canvas.Point,  //<- can be ignored if Startx/y and Endx/y are defined

	Color Canvas.Color
    Alpha byte
}


```
Using a Fragment is as simple, as creating an instance using its structure and passing its Adress to the `Canvas.Draw()` function.
```go
var backgroundFill = CanvasFragment.Fill{ 
    Color: Canvas.COLOR_DARKGRAY 
};

func tick(c *Ca.Canvas, dt float64) Ca.CanvasTickFunction {

    c.Draw(&backgroundFill);

    return tick;

}
```

## Additional Canvas Function

In additon to `Run`, `Apply` and `Draw`, there are also the following functions
available

```go
func (c Canvas) GetPixel(x, y uint16) *uint32 
// Returns a pointer to the pixel at coordinates x,y 
//   or `nil` if x or y are outside the Canvas bounds


func (c Canvas) GetPixelIndex(index uint32) (pixel *uint32) 
// If you know the index of the pixel you can get its Adress 
//   using this method
//   or `nil` if the index is out of the pixel buffers bounds 


func (c Canvas) SetPixel(x, y uint16, color Color) bool {
// Directly sets a Pixel on the canvas to the given color
//   returns true if x,y are valid coordinates 
//           false if x or y are out of the canvas bounds

```


## Additional Package Function

The Package provieds a few helper functions:

```go
func Canvas.ExtractRGB(c Canvas.Color) (r, g, b float64) 
//To split a given color into its conterparts 
//  For technical reasons concerning the next function, the results
//  a are cast as `float64`


func Canvas.BlendPixel(existingPixel *uint32, newPixel uint32, factor float64) 
// Mixes the color of `newPixel` into `*existingPixel` by the given factor
//  factor 1.0 => the new pixel completely overrides the old one
//  factor 0.5 => Existing Pixel becomes a half and half blend of 
//                Itself and the new Pixel
//  factor 0.0 => the old Pixel stays unchanged


func Canvas.CombineRGB(r, g, b byte) Canvas.Color
// A Helper to create a Color from 8bit r, g and b values
```

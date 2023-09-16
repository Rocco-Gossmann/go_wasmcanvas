# Go-WASMCanvas

Proves everything nessary to control a HTML-Canvas Pixel via Go.

## HTML-Preparation:
This is only usable in WebBrowsers, therefore some preparations need to be made.

### 1. First you need the __wasm_exec.js__ from `$GOROOT/misc/wasm`  
In Bash, type 
```bash 
cd `go env GOROOT`/misc/wasm
```
To Find it.

---
### 2. The WebAssembly must be run from a WebWorker 
Reason being, as to not impact the UI

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

Each event send by or to GoWasm will be structured like follows:
```typescript
[ task: string,  ...args: any]
```
and Array whos first element always determins, how the other elements
should be interpreted.

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



// TODO: Describe available functions.
For now you can have a look at the given example. Sorry.


// TODO: Advanced Drawing Section
For now you can have a look at the given example. Sorry.
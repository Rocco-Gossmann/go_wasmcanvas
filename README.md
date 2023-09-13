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
    let canv
    let ctx 

    worker.addEventListener("message", (ev) => {
        if (ev.data instanceof Array) {
            switch (ev.data[0]) {
                case "createCanvas": 
                    // [ 'createCanvas', canvasWidth, canvasHeight ]
                    canv = document.createElement("canvas");
                    canv.width = ev.data[1];
                    canv.height = ev.data[2];
                    canv.className="go-wasm-canvas";
                    ctx = canv.getContext("2d");

                    document.body.appendChild(canv);
                    break;


                case "vblank": 
                    // [ 'vblank` ] => [ 'vblankdone` ]
                    const imgDat = new ImageData(ev.data[1], 320, 200);
                    window.requestAnimationFrame( () => { 
                        ctx.putImageData(imgDat, 0, 0);
                        worker.postMessage(["vblankdone"]) 
                    }); 

                    break;

                default:
                    console.log(
                        "[WORKER MESSAGE]", ev.data[0], " in ", ev
                    );
            }
        }
        else console.log("unknwon worker message", ev.data);
    }, {capture: true})
}
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

// TODO: Describe available functions.

For now you can have a look at the given example on how to use it. Sorry.
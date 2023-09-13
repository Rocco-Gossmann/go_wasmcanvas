# Go-WASMCanvas

Proves everything nessary to control a HTML-Canvas Pixel via Go.

## Usage:
- The project must be run from a WebWorker
```javascript
if (WebAssembly) {

    const worker = new Worker("./worker.js");
    worker.addEventListener("message", (ev) => {
        if (ev.data instanceof Array) {
            switch (ev.data[0]) {
                case "createCanvas": 
                    initCanvas(worker, ev.data[1], ev.data[2]); 
                    break;
            }
        }
        else console.log("unknwon worker message", ev.data);
    }, {once: true, capture: true})

}
else alert("Your Browser does not support WebAssembly");

```

---


The worker will then send a `postMessage` with the following array as data.
```javascript
var msg = [ "createCanvas", pixelwidth, pixelheight ];
```
---

The Main JS-Thread should use this data, to create the Canvas inside the DOM.

```javascript
const canv = document.createElement("canvas");
canv.width = width;
canv.height = height;
canv.className="go-wasm-canvas";

ctx = canv.getContext("2d");

document.body.appendChild(canv);
```


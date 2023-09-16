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


/**
 * Reacts to screesize changes and makes sure the canvas always fills
 * the Screen, as much as possible without distortion
 * @param {HTMLCanvasElement} canv 
 */
function reactToScreenSize(canv) {
    const handleScreenSize = () => {
        const winWidth = window.innerWidth;
        const winHeight = window.innerHeight;

        const aspectwin = winWidth / winHeight;
        const aspectcanv = canv.width / canv.height;

        if (aspectwin < aspectcanv) {
            canv.style.width = "100vw";
            canv.style.height = "initial";
        }
        else {
            canv.style.width = "initial";
            canv.style.height = "100vh";
        }
    }
    window.addEventListener("resize", handleScreenSize)
    handleScreenSize();
}

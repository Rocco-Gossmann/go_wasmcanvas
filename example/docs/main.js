if (WebAssembly) {

    const worker = new Worker("./worker.js");
    let ctx 

    worker.addEventListener("message", (ev) => {
        if (ev.data instanceof Array) {
            switch (ev.data[0]) {
                case "createCanvas": 
                    // [ 'createCanvas', canvasWidth, canvasHeight ]
                    const canv = document.createElement("canvas");
                    canv.width = ev.data[1];
                    canv.height = ev.data[2];
                    canv.className="go-wasm-canvas";
                    ctx = canv.getContext("2d");

                    document.body.appendChild(canv);
                    reactToScreenSize(canv)
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

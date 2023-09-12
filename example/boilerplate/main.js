function initCanvas(worker, width, height) {

    const canv = document.createElement("canvas");
    canv.width = width;
    canv.height = height;
    canv.className="go-wasm-canvas";

    ctx = canv.getContext("2d");

    document.body.appendChild(canv);
    // Window UI
    //------------------------------------------------------------
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

    worker.addEventListener( "message", (ev) => {
        if (ev.data instanceof Array) {

            switch (ev.data[0]) {

                case "vblank": window.requestAnimationFrame(
                    () => { 
                        const imgDat = new ImageData(ev.data[1], 320, 200);
                        ctx.putImageData(imgDat, 0, 0);
                        worker.postMessage(["vblankdone"]) 
                    }
                ); break;

                default:
                    console.log(
                        "[WORKER MESSAGE]", ev.data[0], " in ", ev
                    );
            }
        }
        else console.log("[WORKER MESSAGE]", ev.data);
    })
}

if (WebAssembly) {

    const worker = new Worker("./worker.js");

    worker.addEventListener("message", (ev) => {
        if (ev.data instanceof Array) {
            switch (ev.data[0]) {
                case "createCanvas": initCanvas(worker, ev.data[1], ev.data[2]); break;
            }
        }
        else console.log("unknwon worker message", ev.data);
    }, {once: true, capture: true})

}
else alert("Your Browser does not support WebAssembly");

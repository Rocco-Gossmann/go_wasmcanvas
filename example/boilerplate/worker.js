importScripts("./wasm_exec.js");

//==============================================================================
// Worker <=> Main-Thread communication
//==============================================================================
addEventListener("message", (ev) => {
    switch(ev?.data[0]) {
        case "vblankdone": reportVBlank(); break;
        default: console.warn("don't know how to handle message =>", ev?.data);
    }
}, false)

//==============================================================================
// Let's Gooooo !!!
//==============================================================================
const go = new Go();
WebAssembly.instantiateStreaming(fetch("./main.wasm"), go.importObject)
    .then( gowasm => {
        if(!gowasm) {
            console.error("failed to instantiate gowasm", go, gowasm);
            alert("technical error");
            return;
        }

        console.log("lets goooo !!!!!")
        return go.run(gowasm.instance)
    })      
    .then( () => self.postMessage("wasm has ended"))

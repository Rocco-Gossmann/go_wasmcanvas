BUILDDIR:= ./docs
BOILERPLATE:=./boilerplate

BOILERFILES:=$(BUILDDIR)/  $(BUILDDIR)/index.html
JSFILES:=$(BUILDDIR)/wasm_exec.js $(BUILDDIR)/main.js $(BUILDDIR)/worker.js $(BUILDDIR)/polyfills.js

$(BUILDDIR)/main.wasm: main.go $(BOILERFILES) $(JSFILES)
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o $@ .

$(BUILDDIR)/.htaccess: $(BOILERPLATE)/.htaccess
	cp $^ $@ 

$(BUILDDIR)/%.html: $(BOILERPLATE)/%.html
	cp $^ $@ 

$(BUILDDIR)/%.js: $(BOILERPLATE)/%.js
	cp $^ $@ 


.phony: run stop clean remake $(BUILDDIR)/

run:
	./startserver.zsh

stop:
	./stopserver.zsh

remake: 
	rm -f $(BUILDDIR)/main.wasm
	make $(BUILDDIR)/main.wasm

$(BUILDDIR)/:
	mkdir $(BUILDDIR)

clean:
	rm -rf ./build

package main

import (
	"fmt"
	"log"

	"os"

	"github.com/valyala/fasthttp"
)

func main() {

	dir, _ := os.Getwd()
	fs := &fasthttp.FS{
		Root:               dir,
		AcceptByteRange:    true,
		GenerateIndexPages: true,
	}
	fsHandler := fs.NewRequestHandler()

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		fsHandler(ctx)
	}

	fmt.Println("Starting HTTP server on PORT 80")
	go func() {
		if err := fasthttp.ListenAndServe(":80", requestHandler); err != nil {
			log.Fatalln("error in ListenAndServe: " + err.Error())
		}
	}()

	fmt.Printf("Serving files from directory: '%s'\n", dir)

	// Wait forever.
	select {}
}

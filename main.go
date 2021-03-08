package main

import (
	"flag"
	"fmt"

	"os"

	"github.com/valyala/fasthttp"
)

func main() {
	port := 80
	flag.IntVar(&port, "p", port, "define port - shorthand")
	flag.IntVar(&port, "port", port, "define port")

	dir, err := os.Getwd()
	if err != nil {
		dir = "."
	}
	flag.StringVar(&dir, "d", dir, "define file server directory - shorthand")
	flag.StringVar(&dir, "dir", dir, "define file server directory")
	flag.Parse()

	fs := &fasthttp.FS{
		Root:               dir,
		AcceptByteRange:    true,
		GenerateIndexPages: true,
	}
	fsHandler := fs.NewRequestHandler()

	fmt.Printf("Starting File server on http://localhost:%v \n", port)
	fmt.Printf("Serving files from directory: '%s'\n", dir)

	if err := fasthttp.ListenAndServe(
		fmt.Sprintf(":%v", port),
		func(ctx *fasthttp.RequestCtx) { fsHandler(ctx) },
	); err != nil {
		fmt.Printf("error in ListenAndServe: %v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"flag"
	"fmt"
	"strings"

	"os"

	"github.com/valyala/fasthttp"
)

const (
	homeRoute     = "/"
	uploaderRoute = "/~send"
	explorerRoute = "/~explorer"
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

	s := &fasthttp.Server{
		MaxRequestBodySize: 1024 * 1024 * 5, // 100mb
		StreamRequestBody:  true,
		Handler: func(ctx *fasthttp.RequestCtx) {
			pt := string(ctx.Path())
			switch {
			case pt == homeRoute:
				ctx.SetContentType("text/html; charset=utf-8")
				ctx.SetBody([]byte(homePage))

			// serve files from explorer route
			case strings.HasPrefix(pt, explorerRoute):
				uri := ctx.URI()
				uri.SetPath(pt[len(explorerRoute):])
				ctx.Request.SetRequestURI(uri.String())
				fsHandler(ctx)

			// handle uploading
			case pt == uploaderRoute:
				if ctx.IsGet() {
					ctx.SetContentType("text/html; charset=utf-8")
					ctx.SetBody([]byte(uploaderPage))
					return
				}
				// else check for file
				mf, err := ctx.MultipartForm()
				if err != nil {
					fmt.Println("failed in upload parsing form", err)
					ctx.SetBody([]byte("failed-p " + err.Error()))
					return
				}
				if mf.File == nil {
					ctx.SetBody([]byte("no files provided"))
					return
				}
				for _, fls := range mf.File {
					fl := fls[0]
					if err := fasthttp.SaveMultipartFile(fl, fl.Filename); err != nil {
						fmt.Println("failed in saving file", err)
						ctx.SetBody([]byte("failed-s " + err.Error()))
						return
					}
					ctx.SetBody([]byte(`done`))
				}

			default:
				fsHandler(ctx)
			}
		},
	}
	if err := s.ListenAndServe(fmt.Sprintf(":%v", port)); err != nil {
		fmt.Printf("error in ListenAndServe: %v\n", err)
		os.Exit(1)
	}
}

const (
	homePage = `
		<head>
			<meta charSet="utf-8" />
			<meta name="viewport" content="width=device-width, initial-scale=1" />
		</head>
		<style>
			html, body {
				background: black;
				color: white;
				font-family: sans-serif, system-ui, Arial;
			}
			p,a {
				display: block;
				margin-block: 5px;
				color: white;
			}
		</style>
		<body>
			<p>Hello 👋</p>
			<a href="` + uploaderRoute + `">Use this to send</a>
			<a href="` + explorerRoute + `">Explore here</a>
		</body>
	`
	uploaderPage = `
		<head>
			<meta charSet="utf-8" />
			<meta name="viewport" content="width=device-width, initial-scale=1" />
		</head>
		<style>
			html, body {
				background: black;
				color: white;
				font-family: sans-serif, system-ui, Arial;
			}
			p,a {
				display: block;
				margin-block: 5px;
				color: white;
			}
		</style>
		<body>
			<h1>File Upload</h1>
			<form id="form" enctype="multipart/form-data">
				<label for="files">Select files</label>
				<input id="file" type="file" multiple />
				<button type="submit">Upload</button>
			</form>
			<p id="waiting"></p>
			<script>
				const form = document.getElementById("form");
				const input = document.getElementById("file");
				const waiting = document.getElementById("waiting");
				
				const handleSubmit = (event) => {
					event.preventDefault();
					waiting.innerText="uploading..."
				
					const formData = new FormData();
					[...input.files].forEach(
						(file, i) => {
							key = "file"
							if(input.files.length != 1)
								key+= i
							formData.append(key, file);
						}
					)

				
					fetch("` + uploaderRoute + `", {
						method: "post",
						body: formData,
					}).
						then(r => r.text()).
						then(r => {
							alert(r);
							waiting.innerText=""
						}).
						catch(err => alert("Something went wrong!" + err));
				};
				
				form.addEventListener("submit", handleSubmit);
			</script>
		</body>
	`
)

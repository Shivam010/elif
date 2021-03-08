package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func write(w io.Writer, list ...interface{}) {
	fmt.Println(list...)
	_, _ = fmt.Fprintln(w, list...)
}

func greet(w http.ResponseWriter, r *http.Request) {
	defer write(w)
	if r.URL.Scheme == "" {
		r.URL.Scheme = "http"
	}
	write(w, r.Proto, r.Method, "at", r.URL.Scheme+"://"+r.Host+r.RequestURI)
	write(w, "Remote:", r.RemoteAddr)
	write(w, "Headers:")
	for key, val := range r.Header {
		write(w, "\t", key, "-", val)
	}
	write(w, "Content Length:", r.ContentLength)
	if r.Body != nil {
		write(w, "Body")
		body, _ := ioutil.ReadAll(r.Body)
		write(w, string(body))
	}
}

func main() {
	port := 80
	flag.IntVar(&port, "p", port, "define port - shorthand")
	flag.IntVar(&port, "port", port, "define port")
	flag.Parse()

	http.HandleFunc("/", greet)
	portStr := fmt.Sprintf(":%v", port)
	write(ioutil.Discard, "Starting server at http://localhost"+portStr)
	if err := http.ListenAndServe(portStr, nil); err != nil {
		write(ioutil.Discard, "Error starting:", err)
		os.Exit(1)
	}
}

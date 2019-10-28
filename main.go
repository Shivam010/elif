package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println(http.ListenAndServe(":80", http.FileServer(http.Dir("."))))
}

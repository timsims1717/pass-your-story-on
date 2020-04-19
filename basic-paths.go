package main

import (
	"io"
	"net/http"
)

func DefaultEndPoint(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!\n")
}
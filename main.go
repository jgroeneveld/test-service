package main

import (
	"io"
	"net/http"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	io.Copy(w, r.Body)
}

func main() {
	http.HandleFunc("/echo", echoHandler)
	http.ListenAndServe(":8080", nil)
}

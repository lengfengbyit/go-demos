package main

import (
	"net/http"
	"time"
)

func streamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "text/even-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for i := 0; i < 10; i++ {
		_, _ = w.Write([]byte("data: Hello world\n\n"))
		w.(http.Flusher).Flush()
		time.Sleep(time.Second)
	}
}

func main() {
	http.HandleFunc("/", streamHandler)
	http.ListenAndServe(":9090", nil)
}

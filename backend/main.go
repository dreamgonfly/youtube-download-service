package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"youtube-download-backend/server"
	"youtube-download-backend/storage"
	"youtube-download-backend/videodownload"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	st := &storage.GCS{}
	dw := &videodownload.YoutubeDl{}
	srv := server.NewServer(st, dw)

	h2s := &http2.Server{
		MaxReadFrameSize:             16 << 20,
		MaxUploadBufferPerConnection: 1 << 30,
		MaxUploadBufferPerStream:     1 << 30,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	addr := "0.0.0.0:" + port
	server := &http.Server{
		Addr:              addr,
		Handler:           h2c.NewHandler(srv, h2s),
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
	}

	log.Printf("Server up & running at %s", addr)

	return server.ListenAndServe()
}

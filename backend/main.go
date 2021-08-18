package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"youtube-download-backend/internal/config"
	"youtube-download-backend/internal/gcs"
	"youtube-download-backend/internal/logging"
	"youtube-download-backend/internal/server"
	"youtube-download-backend/internal/youtubefile"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const DefaultPort = "8000"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	logging.InitMultiLogger()

	// dependencies
	ctx := context.Background()
	c := youtubefile.Command
	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "storage.NewClient")
	}
	defer gcsClient.Close()
	g := &gcs.Client{Client: gcsClient}
	h := http.DefaultClient
	sf := storage.SignedURL

	app := server.NewServer(ctx, c, g, h, sf)
	h2s := &http2.Server{
		MaxReadFrameSize:             1 << 20, // Default value. 1MB
		MaxUploadBufferPerConnection: 1 << 20, // Default value. 1MB
		MaxUploadBufferPerStream:     1 << 20, // Default value. 1MB
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	addr := "0.0.0.0:" + port
	server := &http.Server{
		Addr:              addr,
		Handler:           h2c.NewHandler(app, h2s),
		ReadTimeout:       0, // Default value.
		ReadHeaderTimeout: 0, // Default value.
		WriteTimeout:      0, // Default value.
	}

	logging.Logger.Info(fmt.Sprintf("Server up & running at %s on %s environment", addr, config.Env))

	return server.ListenAndServe()
}

package server

import (
	"context"
	"fmt"
	"net/http"

	"youtube-download-backend/internal/gcs"
	"youtube-download-backend/internal/httpfile"
	"youtube-download-backend/internal/youtubefile"

	"github.com/gorilla/mux"
)

type Server struct {
	router     *mux.Router
	context    context.Context
	youtubedl  youtubefile.YoutubeDl
	gcsClient  gcs.Clienter
	httpClient httpfile.Clienter
	signFunc   gcs.SignFunc
}

func NewServer(ctx context.Context, c youtubefile.Commander, g gcs.Clienter, h httpfile.Clienter, sf gcs.SignFunc) *Server {
	r := mux.NewRouter()
	s := &Server{
		router:     r,
		context:    ctx,
		youtubedl:  youtubefile.YoutubeDl{ExecCommand: c},
		gcsClient:  g,
		httpClient: h,
		signFunc:   sf,
	}
	s.routes()
	return s
}

// Make Server an http.Handle
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(&w)
	if r.Method == "OPTIONS" { // Handling pre-flight OPTIONS requests
		return
	}
	s.router.ServeHTTP(w, r)
}

func (s *Server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!\n")
	}
}

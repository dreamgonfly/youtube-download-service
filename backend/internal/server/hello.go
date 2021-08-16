package server

import (
	"fmt"
	"net/http"
)

func (s *Server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!\n")
	}
}

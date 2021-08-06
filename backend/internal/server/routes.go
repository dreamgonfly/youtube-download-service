package server

import "net/http"

func (s *Server) routes() {
	s.router.HandleFunc("/hello", s.handleHello()).Methods("GET")
	s.router.HandleFunc("/preview/{id}", s.handlePreview()).Methods("GET")
	s.router.HandleFunc("/update-thumbnail/{id}", s.handleUpdateThumbnail()).Methods("POST")
	s.router.HandleFunc("/download/{id}", s.handleDownload()).Methods("GET")
	s.router.HandleFunc("/save/{id}", s.handleSave()).Methods("GET")
}

func (s *Server) enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

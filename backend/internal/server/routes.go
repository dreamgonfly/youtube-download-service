package server

import . "youtube-download-backend/internal/middleware"

func (s *Server) routes() {
	s.router.HandleFunc("/hello", HandleLogging(s.handleHello())).Methods("GET")
	s.router.HandleFunc("/preview/{id}", HandleLogging(s.handlePreview())).Methods("GET")
	s.router.HandleFunc("/update-thumbnail/{id}", HandleLogging(s.handleUpdateThumbnail())).Methods("POST")
	s.router.HandleFunc("/download/{id}", HandleLogging(s.handleDownload())).Methods("GET")
	s.router.HandleFunc("/play/{id}", HandleLogging(s.handlePlay())).Methods("GET")
}

package server

func (s *Server) routes() {
	s.router.HandleFunc("/hello", s.handleHello()).Methods("GET")
	s.router.HandleFunc("/preview/{id}", s.handlePreview()).Methods("GET")
	s.router.HandleFunc("/update-thumbnail/{id}", s.handleUpdateThumbnail()).Methods("POST")
	s.router.HandleFunc("/download/{id}", s.handleDownload()).Methods("GET")
	s.router.HandleFunc("/play/{id}", s.handlePlay()).Methods("GET")
}

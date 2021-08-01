package server

func (s *Server) routes() {
	s.router.HandleFunc("/hello", s.handleHello()).Methods("GET")
	s.router.HandleFunc("/preview/{id}", s.handlePreview()).Methods("GET")
	s.router.HandleFunc("/download/{id}", s.handleDownload()).Methods("GET")
}

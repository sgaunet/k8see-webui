package main

import (
	"net/http"
)

func (s *appServer) routes() {
	// static := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	// s.router.PathPrefix("/static").Handler(static)

	// http.FS can be used to create a http Filesystem
	var staticFS = http.FS(staticCSS)
	fsStatic := http.FileServer(staticFS)
	// static := http.StripPrefix("/static/", fsStatic)
	s.router.PathPrefix("/static").Handler(fsStatic)
	s.router.HandleFunc("/", s.IndexHandler)
}

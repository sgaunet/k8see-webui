package main

import "net/http"

func (s *appServer) routes() {
	//fs := http.FileServer(http.Dir("./static"))
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	s.router.PathPrefix("/static").Handler(static)
	s.router.HandleFunc("/", s.IndexHandler)
}

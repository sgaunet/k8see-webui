package main

import (
	"fmt"
	"io/fs"
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

	tmplFiles, err := fs.ReadDir(staticCSS, "static")
	if err != nil {
		panic(err)
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}
		fmt.Println(tmpl.Name())
	}

	s.router.HandleFunc("/", s.IndexHandler)
}

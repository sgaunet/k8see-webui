package main

import (
	"fmt"
	"net/http"
	"text/template"
)

type dataErr struct {
	ErrorMsg string
}

func (s *appServer) HandlerError(response http.ResponseWriter, data dataErr) {
	// tmplt := template.New("error.html")
	// tmplt, _ = tmplt.ParseFiles("./templates/error.html")
	tmplt, err := template.ParseFS(htmlFiles, "templates/error.html")
	if err != nil {
		panic(err)
	}

	err = tmplt.Execute(response, data)
	if err != nil {
		fmt.Printf("Error when generating template error: %s\n", err.Error())
	}
}

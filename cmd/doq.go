package main

import (
	"log"
	"net/http"

	"github.com/itsektionen/doq/app/pages"
)

func main() {
	p := pages.New()
	mux := http.NewServeMux()

	mux.Handle("GET /static/", p.HandleStatic())
	mux.HandleFunc("GET /{$}", p.HandleIndex)
	mux.HandleFunc("GET /", p.HandleNotFound)
	mux.HandleFunc("/", p.HandleNotFound)

	log.Println("Server starting on http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}

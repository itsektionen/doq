package main

import (
	"log"
	"net/http"

	"github.com/itsektionen/doq/app/handler"
	"github.com/itsektionen/doq/app/notify"
	"github.com/itsektionen/doq/app/pages"
	"github.com/itsektionen/doq/app/service"
	"github.com/itsektionen/doq/app/storage"
)

func main() {
	// storageAdapter := storage.NewLocalStorage(filepath.Join("tmp", "img"))
	storageAdapter := storage.NewNopStorage()
	photoService := service.NewPhotoService(storageAdapter)
	photoHandler := handler.NewPhotoHandler(photoService)
	p := pages.New()
	notify.SetTemplate(p.Template())

	mux := http.NewServeMux()

	mux.Handle("GET /static/", p.HandleStatic())
	mux.HandleFunc("GET /{$}", p.HandleIndex)
	mux.HandleFunc("GET /camera", p.HandleCamera)
	mux.HandleFunc("POST /upload-photo", photoHandler.HandleUploadPhoto)
	mux.HandleFunc("GET /", p.HandleNotFound)
	mux.HandleFunc("/", p.HandleNotFound)

	log.Println("Server starting on http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}

package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/itsektionen/doq/app/config"
	"github.com/itsektionen/doq/app/handler"
	"github.com/itsektionen/doq/app/notify"
	"github.com/itsektionen/doq/app/pages"
	"github.com/itsektionen/doq/app/service"
	"github.com/itsektionen/doq/app/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	ctx := context.Background()
	var storageAdapter storage.Storage

	switch strings.ToLower(cfg.StorageType) {
	case "file", "local":
		storageAdapter = storage.NewLocalStorage(cfg.FilePath)
		log.Printf("Using Local File storage at: %s", cfg.FilePath)
	case "gdrive", "googledrive":
		gdriveStorage, err := storage.NewGoogleDriveStorage(
			ctx, cfg.GDriveCredentialsFile, cfg.GDriveFolderID)
		if err != nil {
			log.Fatalf("failed to initialize google drive storage: %v", err)
		}
		storageAdapter = gdriveStorage
		log.Printf("Using Google Drive storage with folder ID: %s", cfg.GDriveFolderID)
	default:
		storageAdapter = storage.NewNopStorage()
		log.Println("Using Nop (No-Op) storage")
	}

	photoService := service.NewPhotoService(storageAdapter)
	photoHandler := handler.NewPhotoHandler(photoService)
	p := pages.New()
	notify.SetTemplate(p.Template())

	mux := http.NewServeMux()

	mux.Handle("GET /static/", p.HandleStatic())
	mux.HandleFunc("GET /{$}", p.HandleIndex)
	mux.HandleFunc("GET /camera", p.HandleCamera)
	mux.HandleFunc("POST /upload-photo", photoHandler.HandleUploadPhoto)
	mux.HandleFunc("GET /privacy", p.HandlePrivacy)
	mux.HandleFunc("GET /", p.HandleNotFound)
	mux.HandleFunc("/", p.HandleNotFound)

	log.Printf("Server starting on %s ...", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}

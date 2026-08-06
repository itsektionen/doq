package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/itsektionen/doq/app/notify"
	"github.com/itsektionen/doq/app/service"
)

type PhotoHandler struct {
	photoService *service.PhotoService
}

func NewPhotoHandler(ps *service.PhotoService) *PhotoHandler {
	return &PhotoHandler{photoService: ps}
}

func (h *PhotoHandler) parsePhoto(r *http.Request) ([]byte, string, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		_ = r.ParseForm()
	}

	file, header, err := r.FormFile("photo")
	if err == nil {
		defer file.Close()

		imgBytes, err := io.ReadAll(file)
		if err != nil {
			return nil, "", fmt.Errorf("failed to read photo file: %w", err)
		}

		ext := ".jpg"

		if strings.HasSuffix(strings.ToLower(header.Filename), ".png") {
			ext = ".png"
		}

		return imgBytes, ext, nil
	}

	photo := r.FormValue("photo")
	if photo == "" {
		return nil, "", errors.New("missing photo field")
	}

	ext := ".jpg"
	raw := photo
	if idx := strings.Index(photo, ","); idx != -1 {
		headerPart := photo[:idx]
		if strings.Contains(headerPart, "image/png") {
			ext = ".png"
		} else if strings.Contains(headerPart, "image/webp") {
			ext = ".webp"
		}
		raw = photo[idx+1:]
	}

	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, "", errors.New("invalid base64 image data")
	}

	if len(decoded) == 0 {
		return nil, "", errors.New("empty photo data")
	}

	return decoded, ext, nil
}

func (h *PhotoHandler) HandleUploadPhoto(w http.ResponseWriter, r *http.Request) {
	imgBytes, ext, err := h.parsePhoto(r)
	if err != nil {
		_ = notify.NotifyError(w, r, "Failed to save photo")
		return
	}

	_, err = h.photoService.SavePhoto(r.Context(), imgBytes, ext)
	if err != nil {
		_ = notify.NotifyError(w, r, "Failed to save photo")
		return
	}

	_ = notify.NotifySuccess(w, r, "Photo saved successfully!")
}

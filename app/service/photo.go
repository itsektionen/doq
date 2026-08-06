package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/itsektionen/doq/app/storage"
)

type PhotoService struct {
	storage storage.Storage
}

func NewPhotoService(s storage.Storage) *PhotoService {
	return &PhotoService{storage: s}
}

func (s *PhotoService) SavePhoto(ctx context.Context, imgBytes []byte, ext string) (string, error) {
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	err := s.storage.Save(ctx, filename, bytes.NewReader(imgBytes))
	if err != nil {
		return "", err
	}
	return filename, nil
}

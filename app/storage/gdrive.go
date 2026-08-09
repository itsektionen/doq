package storage

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDriveStorage struct {
	service  *drive.Service
	folderID string
}

func NewGoogleDriveStorage(ctx context.Context, keyFilePath string, folderID string) (*GoogleDriveStorage, error) {
	srv, err := drive.NewService(
		ctx,
		option.WithCredentialsFile(keyFilePath),
		option.WithScopes(drive.DriveScope),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create google drive service: %w", err)
	}

	return &GoogleDriveStorage{
		service:  srv,
		folderID: folderID,
	}, nil
}

func (s *GoogleDriveStorage) Save(ctx context.Context, filename string, r io.Reader) error {
	fileMeta := &drive.File{
		Name: filename,
	}

	if s.folderID != "" {
		fileMeta.Parents = []string{s.folderID}
	}

	_, err := s.service.Files.Create(fileMeta).Media(r).SupportsAllDrives(true).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to upload file to google drive: %w", err)
	}

	return nil
}

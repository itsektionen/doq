package storage_test

import (
	"context"
	"testing"

	"github.com/itsektionen/doq/app/storage"
)

func TestGoogleDriveStorage_Interface(t *testing.T) {
	var _ storage.Storage = (*storage.GoogleDriveStorage)(nil)
}

func TestNewGoogleDriveStorage_InvalidFilePath(t *testing.T) {
	_, err := storage.NewGoogleDriveStorage(context.Background(), "non_existent_file.json", "")
	if err == nil {
		t.Errorf("expected error initializing with non-existent file, got nil")
	}
}

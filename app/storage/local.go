package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	dir string
}

func NewLocalStorage(dir string) *LocalStorage {
	return &LocalStorage{dir: dir}
}

func (s *LocalStorage) Save(ctx context.Context, filename string, r io.Reader) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(s.dir, filename)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, r); err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	return nil
}

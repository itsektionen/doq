package storage

import (
	"context"
	"io"
)

type NopStorage struct{}

func NewNopStorage() *NopStorage {
	return &NopStorage{}
}

func (s *NopStorage) Save(ctx context.Context, filename string, r io.Reader) error {
	return nil
}

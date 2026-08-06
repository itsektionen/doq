package storage

import (
	"context"
	"io"
)

type Storage interface {
	Save(ctx context.Context, filename string, r io.Reader) error
}

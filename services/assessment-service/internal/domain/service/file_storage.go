package service

import (
	"context"
	"io"
)

type FileStorage interface {
	Upload(ctx context.Context, path string, content io.Reader) (string, error)
	GetURL(ctx context.Context, path string) (string, error)
}

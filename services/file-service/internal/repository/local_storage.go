package repository

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorageProvider struct {
	BaseDir string
	BaseURL string
}

func NewLocalStorageProvider(baseDir, baseURL string) (*LocalStorageProvider, error) {
	if err := os.MkdirAll(baseDir, 0750); err != nil {
		return nil, err
	}
	return &LocalStorageProvider{
		BaseDir: baseDir,
		BaseURL: baseURL,
	}, nil
}

func (s *LocalStorageProvider) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, path string) (string, error) {
	fullPath := filepath.Join(s.BaseDir, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", err
	}

	dst, err := os.Create(fullPath) // #nosec G304
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return path, nil
}

func (s *LocalStorageProvider) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.BaseDir, path)
	return os.Remove(fullPath)
}

func (s *LocalStorageProvider) GetURL(ctx context.Context, path string) (string, error) {
	// Normalize path to ensure forward slashes
	path = strings.ReplaceAll(path, "\\", "/")
	if strings.HasSuffix(s.BaseURL, "/") {
		return fmt.Sprintf("%s%s", s.BaseURL, path), nil
	}
	return fmt.Sprintf("%s/%s", s.BaseURL, path), nil
}

func (s *LocalStorageProvider) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.BaseDir, path)
	return os.Open(fullPath) // #nosec G304
}

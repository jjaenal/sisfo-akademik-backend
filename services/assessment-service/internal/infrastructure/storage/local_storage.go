package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/service"
)

type localStorage struct {
	basePath string
	baseURL  string
}

func NewLocalStorage(basePath, baseURL string) service.FileStorage {
	// Ensure base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		panic(fmt.Sprintf("failed to create storage directory: %v", err))
	}
	return &localStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

func (s *localStorage) Upload(ctx context.Context, path string, content io.Reader) (string, error) {
	fullPath := filepath.Join(s.basePath, path)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, content); err != nil {
		return "", fmt.Errorf("failed to write file content: %w", err)
	}

	return fmt.Sprintf("%s/%s", s.baseURL, path), nil
}

func (s *localStorage) GetURL(ctx context.Context, path string) (string, error) {
	return fmt.Sprintf("%s/%s", s.baseURL, path), nil
}

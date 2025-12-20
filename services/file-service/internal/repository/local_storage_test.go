package repository

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalStorage_Upload_Get_Delete(t *testing.T) {
	dir := t.TempDir()
	s, err := NewLocalStorageProvider(dir, "http://localhost:9098")
	if err != nil {
		t.Fatalf("new storage: %v", err)
	}
	content := []byte("hello")
	file := &fileLike{Reader: bytes.NewReader(content)}
	header := &multipart.FileHeader{Filename: "hello.txt"}
	path := "uploads/sub/hello.txt"
	p, err := s.Upload(context.Background(), file, header, path)
	if err != nil {
		t.Fatalf("upload: %v", err)
	}
	if p != path {
		t.Fatalf("path mismatch")
	}
	u, err := s.GetURL(context.Background(), path)
	if err != nil {
		t.Fatalf("url: %v", err)
	}
	if u != "http://localhost:9098/"+path {
		t.Fatalf("unexpected url: %s", u)
	}
	rc, err := s.Get(context.Background(), path)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	data, _ := io.ReadAll(rc)
	_ = rc.Close()
	if string(data) != "hello" {
		t.Fatalf("content mismatch")
	}
	if err := s.Delete(context.Background(), path); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, path)); !os.IsNotExist(err) {
		t.Fatalf("file should be deleted")
	}
}

type fileLike struct {
	*bytes.Reader
}
func (f *fileLike) Close() error { return nil }

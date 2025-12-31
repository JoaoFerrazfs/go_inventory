package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct{
    BaseDir string
    BaseURL string // e.g. http://localhost:3000/
}

func NewLocalStorage(baseDir string, baseURL string) *LocalStorage{
    return &LocalStorage{BaseDir: baseDir, BaseURL: baseURL}
}

func (s *LocalStorage) Upload(path string, content io.Reader) (string, error){
    fullPath := filepath.Join(s.BaseDir, path)
    if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
        return "", err
    }
    f, err := os.Create(fullPath)
    if err != nil {
        return "", err
    }
    defer f.Close()
    if _, err := io.Copy(f, content); err != nil {
        return "", err
    }
    // return the relative path inside storage
    return path, nil
}

func (s *LocalStorage) GetURL(path string) (string, error){
    // path is relative inside base dir
    return fmt.Sprintf("%s%s", s.BaseURL, path), nil
}

func (s *LocalStorage) Delete(path string) error{
    fullPath := filepath.Join(s.BaseDir, path)
    return os.Remove(fullPath)
}

package storage

import "io"

type Storage interface {
    Upload(path string, content io.Reader) (string, error)
    GetURL(path string) (string, error)
    Delete(path string) error
}

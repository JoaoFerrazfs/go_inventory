package mocks

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type MockStorage struct{ mock.Mock }

func (m *MockStorage) Upload(path string, content io.Reader) (string, error) {
	args := m.Called(path, content)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) GetURL(path string) (string, error) {
	args := m.Called(path)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) Delete(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

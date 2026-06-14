package mock

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

// Storage is an in-memory implementation of contract.Storage for testing.
type Storage struct {
	mu    sync.RWMutex
	files map[string][]byte
}

func NewStorage() *Storage {
	return &Storage{files: make(map[string][]byte)}
}

func (s *Storage) Upload(_ context.Context, bucket, name string, reader io.Reader, _ int64, _ string) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	key := s.key(bucket, name)
	s.files[key] = data
	return key, nil
}

func (s *Storage) Download(_ context.Context, bucket, name string) (io.ReadCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.files[s.key(bucket, name)]
	if !ok {
		return nil, fmt.Errorf("mock: file not found: %s/%s", bucket, name)
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (s *Storage) Delete(_ context.Context, bucket, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.files, s.key(bucket, name))
	return nil
}

func (s *Storage) URL(_ context.Context, bucket, name string, _ time.Duration) (string, error) {
	return fmt.Sprintf("http://mock-storage/%s/%s", bucket, name), nil
}

// Has returns true if the file exists in the mock store — useful for assertions.
func (s *Storage) Has(bucket, name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.files[s.key(bucket, name)]
	return ok
}

func (s *Storage) key(bucket, name string) string {
	return fmt.Sprintf("%s/%s", bucket, name)
}

package storage

import (
	"context"
	"sync"
)

// InMemoryStorage is an in-memory implementation of the Storage interface.
type InMemoryStorage struct {
	mu   sync.RWMutex
	urls map[string]bool
}

// NewInMemoryStorage creates a new in-memory storage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		urls: make(map[string]bool),
	}
}

// StoreURLs stores the given URLs in memory.
func (s *InMemoryStorage) StoreURLs(ctx context.Context, urls map[string]bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for u := range urls {
		s.urls[u] = true
	}
	return nil
}

// GetURLs retrieves all URLs from memory.
func (s *InMemoryStorage) GetURLs(ctx context.Context) (map[string]bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.urls, nil
}

// GetSites retrieves all sites from memory.
func (s *InMemoryStorage) GetSites(ctx context.Context) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sites := []string{"https://djinni.co/jobs/keyword-golang/",
		"https://www.globallogic.com/ua/career-search-page/?keywords=golang&experience=none&location=ukraine/"}

	// This is a mock implementation. In a real scenario, this would fetch links from a persistent storage.
	return sites, nil
}

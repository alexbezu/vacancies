package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/alexbezu/vacancies/internal/model"
	"google.golang.org/api/iterator"
)

type FireStore struct {
	mu   sync.RWMutex
	urls map[string]bool
	fs   *firestore.Client
}

func NewFireStore(ctx context.Context, projectID string) (*FireStore, error) {
	fs, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &FireStore{
		urls: make(map[string]bool),
		fs:   fs,
	}, nil
}

func (s *FireStore) StoreURLs(ctx context.Context, urls map[string]bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for url := range urls {
		s.urls[url] = true

		_, _, err := s.fs.Collection("positions").Add(ctx, model.Position{
			URL:       url,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			return fmt.Errorf("failed to add an url %s with err: %s", url, err)
		}
	}
	return nil
}

func (s *FireStore) GetURLs(ctx context.Context) (map[string]bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var p model.Position

	iter := s.fs.Collection("positions").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate next url: %s", err)
		}
		if err = doc.DataTo(&p); err != nil {
			return nil, fmt.Errorf("failed to parse an url from json: %s", err)
		}
		s.urls[p.URL] = true
	}

	return s.urls, nil
}

func (s *FireStore) GetSites(ctx context.Context) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sites := []string{"https://djinni.co/jobs/keyword-golang/",
		"https://www.globallogic.com/ua/career-search-page/?keywords=golang&experience=none&location=ukraine/"}

	return sites, nil
}

package service

import (
	"context"

	"github.com/alexbezu/vacancies/internal/model"
	"github.com/sirupsen/logrus"
)

//go:generate go tool go.uber.org/mock/mockgen -source=$GOFILE -destination=mocks/service.go -package=mocks Storage

// Storage is the interface for storing and retrieving URLs.
type Storage interface {
	StoreURLs(ctx context.Context, urls map[string]bool) error
	GetURLs(ctx context.Context) (map[string]bool, error)
	GetSites(ctx context.Context) ([]model.JobSite, error)
}

// Webhook is the interface for sending notifications.
type Webhook interface {
	Send(ctx context.Context, message string) error
}

type Scraper interface {
	UrlsFromSite(ctx context.Context, link, filter string) ([]string, error)
}

// Service is the main service for vacancies.
type Service struct {
	storage Storage
	scraper Scraper
	webhook Webhook
	log     *logrus.Logger
}

// New creates a new vacancies service.
func New(storage Storage, scraper Scraper, webhook Webhook, log *logrus.Logger) *Service {
	return &Service{
		storage: storage,
		scraper: scraper,
		webhook: webhook,
		log:     log,
	}
}

// ProcessURLs processes the URLs from storage.
func (s *Service) ProcessURLs(ctx context.Context) error {
	// Latest urls from vacancy sites
	var urls []string
	sites, err := s.storage.GetSites(ctx)
	if err != nil {
		return err
	}

	for _, site := range sites {
		vacancies, err := s.scraper.UrlsFromSite(ctx, site.URL, site.Filter)
		if err != nil {
			continue
		}
		urls = append(urls, vacancies...)
	}
	// Get existing URLs from storage.
	existingURLs, err := s.storage.GetURLs(ctx)
	if err != nil {
		return err
	}

	// Find new URLs (have to be added to db)
	newURLs := make(map[string]bool)
	for _, u := range urls {
		if _, ok := existingURLs[u]; !ok {
			newURLs[u] = true
		}
	}

	// Store the new URLs.
	if len(newURLs) > 0 {
		if err := s.storage.StoreURLs(ctx, newURLs); err != nil {
			return err
		}

		// Send a notification for each new URL.
		for u := range newURLs {
			if err := s.webhook.Send(ctx, u); err != nil {
				s.log.WithError(err).Error("failed to send webhook")
			}
		}
	}

	return nil
}

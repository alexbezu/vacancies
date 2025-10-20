package service

import (
	"context"
	"regexp"
	"time"

	"fmt"
	"net/http"
	"net/url"

	"github.com/alexbezu/vacancies/internal/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

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

// Service is the main service for vacancies.
type Service struct {
	storage Storage
	webhook Webhook
	log     *logrus.Logger
}

// New creates a new vacancies service.
func New(storage Storage, webhook Webhook, log *logrus.Logger) *Service {
	return &Service{
		storage: storage,
		webhook: webhook,
		log:     log,
	}
}

func (s *Service) urlsFromSite(ctx context.Context, link, filter string) ([]string, error) {
	var ret []string

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() // Ensure the cancel function is called to release resources

	// Create a new HTTP GET request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil, err
	}

	// Create an HTTP client
	client := &http.Client{}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error performing request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					// Check if we have relative link. If so, add a host
					u, err := url.Parse(attr.Val)
					if err != nil {
						s.log.Error(err)
					}
					if u.Host == "" {
						l, err := url.Parse(link)
						if err != nil {
							s.log.Error(err)
						}
						u.Host = l.Host
						u.Scheme = l.Scheme
					}

					if filter != "" {
						r, err := regexp.Compile(filter)
						if err != nil {
							s.log.Error(err)
						} else {
							str := u.String()
							if r.MatchString(str) {
								ret = append(ret, str)
							}
						}
					} else {
						ret = append(ret, u.String())
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return ret, nil
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
		vacancies, err := s.urlsFromSite(ctx, site.URL, site.Filter)
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

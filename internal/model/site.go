package model

import "time"

type JobSite struct {
	ID            string         `firestore:"-" json:"id"` // "-" ignores this field for Firestore, we'll use it as the Doc ID
	Name          string         `firestore:"name" json:"name"`
	URL           string         `firestore:"url" json:"url"`
	LastScraped   *time.Time     `firestore:"lastScraped,omitempty" json:"lastScraped,omitempty"` // omitempty for optional fields
	Status        string         `firestore:"status" json:"status"`
	ScraperConfig map[string]any `firestore:"scraperConfig,omitempty" json:"scraperConfig,omitempty"`
	CreatedAt     time.Time      `firestore:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time      `firestore:"updatedAt" json:"updatedAt"`
}

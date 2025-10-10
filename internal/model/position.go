package model

import "time"

type Position struct {
	ID                  string         `firestore:"-" json:"id"`
	Title               string         `firestore:"title" json:"title"`
	Company             string         `firestore:"company" json:"company"`
	Location            string         `firestore:"location" json:"location"`
	URL                 string         `firestore:"url" json:"url"`
	Description         string         `firestore:"description" json:"description"`
	JobSiteRef          string         `firestore:"jobSiteRef" json:"jobSiteRef"` // Storing ID as string for simplicity
	PostedDate          *time.Time     `firestore:"postedDate,omitempty" json:"postedDate,omitempty"`
	ScrapedDate         time.Time      `firestore:"scrapedDate" json:"scrapedDate"`
	Status              string         `firestore:"status" json:"status"`
	Keywords            []string       `firestore:"keywords,omitempty" json:"keywords,omitempty"`
	SalaryRange         map[string]any `firestore:"salaryRange,omitempty" json:"salaryRange,omitempty"`
	ApplicationDeadline *time.Time     `firestore:"applicationDeadline,omitempty" json:"applicationDeadline,omitempty"`
	CreatedAt           time.Time      `firestore:"createdAt" json:"createdAt"`
	UpdatedAt           time.Time      `firestore:"updatedAt" json:"updatedAt"`
}

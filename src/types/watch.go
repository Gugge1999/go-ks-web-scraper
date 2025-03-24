package types

import (
	"time"
)

type Watch struct {
	Id             string         `json:"id"`
	WatchToScrape  string         `json:"watchToScrape"`
	Label          string         `json:"label"`
	Watches        string         `json:"-"`
	Active         bool           `json:"active"`
	LastEmailSent  *time.Time     `json:"lastEmailSent"`
	Added          time.Time      `json:"added"`
	ScrapedWatches []ScrapedWatch `json:"watches"`
	Notifications  []time.Time    `json:"notifications"`
}

type ScrapedWatch struct {
	Name       string `json:"name"`
	PostedDate string `json:"postedDate"`
	Link       string `json:"link"`
}

type SaveWatchDto struct {
	WatchToScrape string `json:"watchToScrape"`
	Label         string `json:"label"`
}

package types

import (
	"time"
)

type Watch struct {
	Id            string       `json:"id"`
	WatchToScrape string       `json:"watchToScrape"`
	Label         string       `json:"label"`
	Watches       string       `json:"-"`
	Active        bool         `json:"active"`
	LastEmailSent *time.Time   `json:"lastEmailSent"`
	Added         time.Time    `json:"added"`
	Notifications []time.Time  `json:"notifications"`
	LatestWatch   ScrapedWatch `json:"latestWatch"`
}

type ScrapedWatch struct {
	Name       string `json:"name"`
	PostedDate string `json:"postedDate"`
	Link       string `json:"link"`
}

type SaveWatchDto struct {
	WatchToScrape string `json:"watchToScrape" binding:"required"`
	Label         string `json:"label" binding:"required"`
}

type ToggleActiveStatusesDto struct {
	Ids []string `json:"ids" binding:"required"`
	// OBS: Den kan inte vara binding:"required" eftersom c.ShouldBindJSON kommer returnera fel då
	NewActiveStatus bool `json:"newActiveStatus"`
}

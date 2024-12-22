package types

import (
	"database/sql"
	"time"
)

type Watch struct {
	Id            string       `json:"id"` // TODO: Hur ska guid hanteras?
	WatchToScrape string       `json:"watchToScrape"`
	Label         string       `json:"label"`
	Watches       string       `json:"watches"` // TODO: Hantera det så: https://stackoverflow.com/a/75944972
	Active        bool         `json:"active"`
	LastEmailSent sql.NullTime `json:"lastEmailSent"`
	Added         time.Time    `json:"added"`
}

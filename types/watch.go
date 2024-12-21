package types

import (
	"database/sql"
	"time"
)

type Watch struct {
	Id            string // TODO: Hur ska guid hanteras?
	WatchToScrape string
	Label         string
	Watches       string // TODO: Hantera det så: https://stackoverflow.com/a/75944972
	Active        bool
	LastEmailSent sql.NullTime
	Added         time.Time
}

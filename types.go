package main

import (
	"database/sql"
	"time"
)

type watch struct {
	id            string // TODO: Hur ska guid hanteras?
	watchToScrape string
	label         string
	watches       string // TODO: Hantera det så: https://stackoverflow.com/a/75944972
	active        bool
	lastEmailSent sql.NullTime
	added         time.Time
}

package types

import (
	"time"
)

type Notification struct {
	Id      string    `json:"id"`
	WatchId string    `json:"watch_id"`
	Sent    time.Time `json:"sent"`
}

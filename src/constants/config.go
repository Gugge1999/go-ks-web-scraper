package constants

import (
	"github.com/gorilla/websocket"
	"net/http"
	"os"
)

const IntervalInMin uint = 10

func IntervalInMs() uint {
	env := os.Getenv("ENV")

	if env == "dev" {
		return IntervalInMin * 1_500
	}

	return IntervalInMin * 60_000
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

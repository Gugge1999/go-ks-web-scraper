package constants

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gorilla/websocket"
)

const IntervalInMin uint = 10

func IntervalInMs() uint {
	if os.Getenv("ENV") == "dev" {
		return IntervalInMin * 1_500
	}

	return IntervalInMin * 60_000
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":3000"
	}

	return ":" + port
}

var CorsConfig = cors.New(cors.Config{
	AllowOrigins:  []string{"*"},
	AllowMethods:  []string{"*"},
	AllowHeaders:  []string{"*"},
	AllowWildcard: true,
})

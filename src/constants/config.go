package constants

import (
	"os"

	"github.com/gin-contrib/cors"
)

const IntervalInMin uint = 10

func IntervalInMs() uint {
	if os.Getenv("ENV") == "dev" {
		return IntervalInMin * 1_500
	}

	return IntervalInMin * 60_000
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

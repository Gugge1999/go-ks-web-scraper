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

var CorsConfig = cors.New(cors.Config{
	AllowOrigins:  []string{"*"},
	AllowMethods:  []string{"*"},
	AllowHeaders:  []string{"*"},
	AllowWildcard: true,
})

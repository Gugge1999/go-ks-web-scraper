package constants

import "os"

const IntervalInMin uint = 10

func IntervalInMs() uint {
	env := os.Getenv("ENV")

	if env == "dev" {
		return IntervalInMin * 1_500
	}

	return IntervalInMin * 60_000
}

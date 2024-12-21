package constants

import "os"

const intervalInMin uint32 = 10

func IntervalInMs() uint32 {
	env := os.Getenv("ENV")

	if env == "dev" {
		return intervalInMin * 1_500
	}

	return intervalInMin * 60_000
}

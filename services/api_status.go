package services

import (
	"ks-web-scraper/constants"
	"ks-web-scraper/types"
	"runtime"
	"runtime/metrics"
	"time"
)

func GetApiStatus(startTime time.Time) types.ApiStatus {
	return types.ApiStatus{
		Status:                    "active",
		ScrapingIntervalInMinutes: constants.IntervalInMin,
		NumberOfCpus:              runtime.NumCPU(),
		MemoryUsage:               getMemoryUsageInMb(),
		Uptime:                    getUptime(startTime),
	}
}

// TODO: Den här verkar endast öka med belastning men minskar aldrig
func getMemoryUsageInMb() uint64 {
	const myMetric = "/memory/classes/total:bytes"

	sample := make([]metrics.Sample, 1)
	sample[0].Name = myMetric

	// Sample the metric.
	metrics.Read(sample)

	bytesInMb := sample[0].Value.Uint64() / 1024 / 1024

	return bytesInMb
}

func getUptime(startTime time.Time) types.Uptime {
	uptime := time.Since(startTime)
	seconds := uint8(uptime.Seconds()) % 60
	minutes := uint8(uptime.Minutes()) % 60
	hours := uint8(uptime.Hours()) % 24
	days := uint16(float64(hours/24)) % 30
	months := uint8(float64(days/30)) % 12
	years := days / 365

	return types.Uptime{
		Seconds: seconds,
		Minutes: minutes,
		Hours:   hours,
		Days:    days,
		Months:  months,
		Years:   years,
	}
}

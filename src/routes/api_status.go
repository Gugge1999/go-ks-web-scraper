package routes

import (
	"ks-web-scraper/src/constants"
	"ks-web-scraper/src/logger"
	"ks-web-scraper/src/types"
	"net/http"
	"runtime"
	"runtime/metrics"
	"sync/atomic"
	"time"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var startTime = time.Now()
var connectionCount uint32

func ApiRoutesApiStatus(router *gin.Engine) {
	router.GET("/api/status", func(c *gin.Context) {
		logger := logger.GetLogger()

		conn, wsError := upgrader.Upgrade(c.Writer, c.Request, nil)

		if wsError != nil {
			logger.Error().Msg("Kunde inte skapa websocket. Error: " + wsError.Error())
			return
		}

		defer conn.Close()

		atomic.AddUint32(&connectionCount, 1)

		for {
			status := getApiStatus(startTime)

			err := conn.WriteJSON(status)

			if err != nil {
				atomic.AddUint32(&connectionCount, ^uint32(0))
				logger.Info().Msg("User disconnected from WebSocket. Total connected users: " + fmt.Sprint(connectionCount))
				break
			}

			time.Sleep(time.Second)
		}
	})
}

func getApiStatus(startTime time.Time) types.ApiStatus {
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
	seconds := uint(uptime.Seconds()) % 60
	minutes := uint(uptime.Minutes()) % 60
	hours := uint(uptime.Hours()) % 24
	days := uint(float64(hours/24)) % 30
	months := uint(float64(days/30)) % 12
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

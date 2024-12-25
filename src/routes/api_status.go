package routes

import (
	"ks-web-scraper/src/constants"
	"ks-web-scraper/src/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RegisterRoutesApiStatus(router *gin.Engine, startTime time.Time) {
	router.GET("/api/api-status", func(c *gin.Context) {
		conn, wsError := constants.Upgrader.Upgrade(c.Writer, c.Request, nil)

		if wsError != nil {
			log.Error().Msg("Kunde inte skapa websocket: " + wsError.Error())
			return
		}

		defer conn.Close()

		for {
			status := services.GetApiStatus(startTime)

			err := conn.WriteJSON(status)

			if err != nil {
				return
			}

			time.Sleep(time.Second)
		}
	})

}

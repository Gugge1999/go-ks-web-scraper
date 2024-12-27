package routes

import (
	"ks-web-scraper/src/constants"
	"ks-web-scraper/src/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var startTime = time.Now()

func RegisterRoutesApiStatus(router *gin.Engine) {
	router.GET("/api/api-status", func(c *gin.Context) {
		conn, wsError := constants.Upgrader.Upgrade(c.Writer, c.Request, nil)

		if wsError != nil {
			log.Error().Msg("Kunde inte skapa websocket. Error: " + wsError.Error())
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

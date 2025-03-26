package routes

import (
	"ks-web-scraper/src/logger"
	"ks-web-scraper/src/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var startTime = time.Now()

func ApiRoutesApiStatus(router *gin.Engine) {
	router.GET("/api/status", func(c *gin.Context) {
		logger := logger.GetLogger()

		conn, wsError := upgrader.Upgrade(c.Writer, c.Request, nil)

		if wsError != nil {
			logger.Error().Msg("Kunde inte skapa websocket. Error: " + wsError.Error())
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

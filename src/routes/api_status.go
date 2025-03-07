package routes

import (
	"ks-web-scraper/src/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var startTime = time.Now()

func RegisterRoutesApiStatus(router *gin.Engine) {
	router.GET("/api/api-status", func(c *gin.Context) {
		conn, wsError := upgrader.Upgrade(c.Writer, c.Request, nil)

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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

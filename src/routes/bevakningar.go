package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"ks-web-scraper/src/database"
	"os"
)

func RegisterRoutesBevakningar(router *gin.Engine, conn *pgx.Conn) {
	router.GET("/api/bevakningar/all-watches", func(c *gin.Context) {
		allNotifications := database.GetAllNotifications(conn)
		allWatches := database.GetAllWatches(conn)

		fmt.Fprintf(os.Stderr, "type: %v\n", allNotifications)

		c.JSON(200, allWatches)
	})
}

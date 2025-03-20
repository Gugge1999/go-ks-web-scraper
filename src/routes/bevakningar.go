package routes

import (
	"fmt"
	"ks-web-scraper/src/database"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func RegisterRoutesBevakningar(router *gin.Engine, conn *pgx.Conn) {
	router.GET("/api/bevakningar/all-watches", func(c *gin.Context) {
		// TODO: Ska den meddela användaren om notiser inte kan hämtas?
		allNotifications, _ := database.GetAllNotifications(conn)
		allWatches := database.GetAllWatches(conn)

		fmt.Fprintf(os.Stderr, "type: %v\n", allNotifications)

		c.JSON(200, allWatches)
	})
}

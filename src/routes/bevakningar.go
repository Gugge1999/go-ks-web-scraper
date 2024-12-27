package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"ks-web-scraper/src/database"
)

func RegisterRoutesBevakningar(router *gin.Engine, conn *pgx.Conn) {
	router.GET("/api/bevakningar/all-watches", func(c *gin.Context) {
		allWatches := database.GetAllWatches(conn)

		c.JSON(200, allWatches)
	})
}

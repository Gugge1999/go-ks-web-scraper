package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutesBevakningar(router *gin.Engine) {

	router.GET("/api/bevakningar/all-watches", func(c *gin.Context) {
		// TODO: anropa db
		c.JSON(200, gin.H{"message": "hejsan"})
	})
}

package routes

import (
	"ks-web-scraper/src/database"
	"ks-web-scraper/src/services"
	"ks-web-scraper/src/types"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

const apiBaseUrl = "/api/bevakningar/"

func ApiRoutesBevakningar(router *gin.Engine, conn *pgx.Conn) {
	router.GET(apiBaseUrl+"all-watches", func(c *gin.Context) {
		allNotifications, err1 := database.GetAllNotifications(conn)
		allWatches, err2 := database.GetAllWatches(conn)

		if err1 != nil || err2 != nil {
			c.JSON(500, gin.H{"message": "Kunde inte hämta bevakningar ", "stack": "Error notiser " + err1.Error() + ". Error bevakningar" + err2.Error()})
			return
		}

		res := createWatchDto(allWatches, allNotifications)

		c.JSON(200, res)
	})

	router.POST(apiBaseUrl+"save-watch", func(c *gin.Context) {
		var saveWatchDto types.SaveWatchDto
		// TODO: Här ska man kolla om dto innehåller rätt properties
		if err := c.ShouldBindJSON(&saveWatchDto); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		scrapedWatches := services.ScrapeWatchInfo(saveWatchDto.WatchToScrape)
		dbRes := database.SaveWatch(conn, saveWatchDto, scrapedWatches)

		dbRes[0].Notifications = []time.Time{}

		c.JSON(200, dbRes[0])
	})

	router.DELETE(apiBaseUrl+"delete-watch/:id", func(c *gin.Context) {
		id := c.Param("id")
		dbRes, err := database.DeleteWatch(conn, id)

		if err != nil {
			c.JSON(500, gin.H{"message": "Kunde inte radera bevakning med id: " + id})
			return
		}

		c.JSON(200, gin.H{"deleteWatchId": dbRes})
	})
}

func createWatchDto(watches []types.Watch, notifications []types.Notification) []types.Watch {
	var watchDtos []types.Watch

	for _, w := range watches {

		var notiserForBevakning []time.Time
		for _, notis := range notifications {
			if w.Id == notis.WatchId {
				notiserForBevakning = append(notiserForBevakning, notis.Sent)
			}
		}

		watchDtos = append(watchDtos, types.Watch{
			// TODO: Finns det något sätt att göra detta med spread operator?
			Id:            w.Id,
			WatchToScrape: w.WatchToScrape,
			Label:         w.Label,
			Active:        w.Active,
			LastEmailSent: w.LastEmailSent,
			Added:         w.Added,
			LastestWatch:  w.LastestWatch,
			Notifications: notiserForBevakning,
		})
	}

	return watchDtos

}

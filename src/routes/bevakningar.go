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
		// TODO: Ska den meddela användaren om notiser inte kan hämtas?
		allNotifications, _ := database.GetAllNotifications(conn)
		allWatches := database.GetAllWatches(conn)

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

		// TODO: Den här ska bara skicka tillbaka den senaste klockan, inte alla 30
		c.JSON(200, dbRes[0])
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
			Id:             w.Id,
			WatchToScrape:  w.WatchToScrape,
			Label:          w.Label,
			Active:         w.Active,
			LastEmailSent:  w.LastEmailSent,
			Added:          w.Added,
			ScrapedWatches: w.ScrapedWatches,
			Notifications:  notiserForBevakning,
		})
	}

	return watchDtos

}

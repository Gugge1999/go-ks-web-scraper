package routes

import (
	"ks-web-scraper/src/database"
	"ks-web-scraper/src/types"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func ApiRoutesBevakningar(router *gin.Engine, conn *pgx.Conn) {
	router.GET("/api/bevakningar/all-watches", func(c *gin.Context) {
		// TODO: Ska den meddela användaren om notiser inte kan hämtas?
		allNotifications, _ := database.GetAllNotifications(conn)
		allWatches := database.GetAllWatches(conn)

		res := createWatchDto(allWatches, allNotifications)

		c.JSON(200, res)
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

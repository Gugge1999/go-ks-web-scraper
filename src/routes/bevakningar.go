package routes

import (
	"ks-web-scraper/src/database"
	"ks-web-scraper/src/services"
	"ks-web-scraper/src/types"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const apiBaseUrl = "/api/bevakningar/"

func ApiRoutesBevakningar(router *gin.Engine, dbPoolConn *pgxpool.Pool) {
	router.GET(apiBaseUrl+"all-watches", func(c *gin.Context) {
		// TODO: Fixa goroutine
		allNotifications, err1 := database.GetAllNotifications(dbPoolConn)
		allWatches, err2 := database.GetAllWatches(dbPoolConn)

		if err1 != nil || err2 != nil {
			c.JSON(500, gin.H{"message": "Kunde inte hämta bevakningar ", "stack": "Error notiser " + err1.Error() + ". Error bevakningar" + err2.Error()})
			return
		}

		dbRes := createWatchDto(allWatches, allNotifications)

		c.JSON(200, dbRes)
	})

	router.POST(apiBaseUrl+"save-watch", func(c *gin.Context) {
		var saveWatchDto types.SaveWatchDto

		validatedSaveWatchDto, shouldReturn := validateSaveWatchBody(c, saveWatchDto)
		if shouldReturn {
			return
		}

		watchToScrapeUrl := "https://klocksnack.se/search/1/?q=" + validatedSaveWatchDto.WatchToScrape + "&t=post&c[child_nodes]=1&c[nodes][0]=11&c[title_only]=1&o=date"

		scrapedWatches := services.ScrapeWatchInfo(watchToScrapeUrl)
		if len(scrapedWatches) == 0 {
			c.JSON(422, gin.H{"message": "Kunde inte hitta några klockor med sökordet: " + validatedSaveWatchDto.WatchToScrape})
			return
		}

		dbRes, err := database.SaveWatch(dbPoolConn, validatedSaveWatchDto.Label, watchToScrapeUrl, scrapedWatches)

		if err != nil {
			c.JSON(500, gin.H{"message": "Kunde inte spara bevakning"})
			return
		}

		dbRes[0].Notifications = []time.Time{}

		c.JSON(200, dbRes[0])
	})

	router.DELETE(apiBaseUrl+"delete-watch/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := uuid.Validate(id); err != nil {
			c.JSON(422, gin.H{"message": "id måste vara av typen uuid v4"})
			return
		}

		dbRes, err := database.DeleteWatch(dbPoolConn, id)

		if err != nil {
			c.JSON(500, gin.H{"message": "Kunde inte radera bevakning med id: " + id})
			return
		}

		c.JSON(200, gin.H{"deleteWatchId": dbRes})
	})

	router.PUT(apiBaseUrl+"toggle-active-statuses", func(c *gin.Context) {
		var toggleActiveStatusesDto types.ToggleActiveStatusesDto

		if err := c.ShouldBindJSON(&toggleActiveStatusesDto); err != nil {
			c.JSON(422, gin.H{"message": "Body måste innehålla ett object med två properties. ids: en array med ids. newActiveStatus: en boolean för ny aktiv status"})
			return
		}

		for _, id := range toggleActiveStatusesDto.Ids {
			if err := uuid.Validate(id); err != nil {
				c.JSON(422, gin.H{"message": "Samtliga ids måste vara av typen uuid v4"})
				return
			}
		}

		dbRes, err := database.ToggleActiveStatuses(dbPoolConn, toggleActiveStatusesDto.Ids, toggleActiveStatusesDto.NewActiveStatus)

		if err != nil {
			c.JSON(500, gin.H{"message": "Kunde inte uppdatera bevakningarna"})
			return
		}

		c.JSON(200, dbRes)
	})
}

func validateSaveWatchBody(c *gin.Context, saveWatchDto types.SaveWatchDto) (types.SaveWatchDto, bool) {
	if err := c.ShouldBindJSON(&saveWatchDto); err != nil {
		c.JSON(422, gin.H{"message": "Body måste finnas och måste innehålla WatchToScrape och Label"})
		return saveWatchDto, true
	}

	if saveWatchDto.WatchToScrape == "" || saveWatchDto.Label == "" {
		c.JSON(422, gin.H{"message": "saveWatchDto måste innehålla WatchToScrape och Label"})
		return saveWatchDto, true
	}

	if len(saveWatchDto.WatchToScrape) <= 3 || len(saveWatchDto.Label) <= 2 {
		c.JSON(422, gin.H{"message": "watchToScrape och label måste vara minst 3 respektive 2 tecken"})
		return saveWatchDto, true
	}

	if len(saveWatchDto.WatchToScrape) >= 35 || len(saveWatchDto.Label) >= 30 {
		c.JSON(422, gin.H{"message": "watchToScrape och label måste vara minst 35 respektive 30 tecken"})
		return saveWatchDto, true
	}

	return saveWatchDto, false
}

func createWatchDto(watches []types.Watch, notifications []types.Notification) []types.Watch {
	var watchDtos []types.Watch

	for _, w := range watches {

		notiserForBevakning := []time.Time{}
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
			LatestWatch:   w.LatestWatch,
			Notifications: notiserForBevakning,
		})
	}

	return watchDtos
}

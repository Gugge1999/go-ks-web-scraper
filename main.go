package main

import (
	"fmt"
	"ks-web-scraper/src/constants"
	"ks-web-scraper/src/database"
	"ks-web-scraper/src/logger"
	"ks-web-scraper/src/routes"
	"ks-web-scraper/src/services"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	initApiMsg := "Init api @ \x1b[32m" + time.Now().Format("15:04:05") + "\x1b[0m\n\n" // 32 = grön. OBS: Format måste vara exakt 15:04:05
	fmt.Fprint(os.Stderr, initApiMsg)

	// OBS: Måste ske tidigt
	services.LoadDotEnvFile()

	dbPoolConn := database.InitDB()

	router := gin.Default()

	router.Use(constants.CorsConfig)

	routes.ApiRoutesApiStatus(router)
	routes.ApiRoutesBevakningar(router, dbPoolConn)

	routerRunErr := router.Run(services.GetPort())

	if routerRunErr != nil {
		logger := logger.GetLogger()

		logger.Error().Msg("Kunde inte starta server:" + routerRunErr.Error())
	}

	// TODO: Kolla om det här är rätt
	defer dbPoolConn.Close()
}

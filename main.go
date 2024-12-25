package main

import (
	"context"
	"fmt"
	"ks-web-scraper/src/constants"
	"ks-web-scraper/src/database"
	"ks-web-scraper/src/routes"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {

	// TODO: Kolla varför måste det vara 15:04:05?
	fmt.Fprintf(os.Stderr, "Init api @ \x1b[%dm%s\x1b[0m\n\n", 32, time.Now().Format("15:04:05"))

	// TODO: Det kanske går att bryta ut skapandet av logger till funktion?
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	runLogFile, logFileError := os.OpenFile("logs/logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if logFileError != nil {
		fmt.Fprintf(os.Stderr, "Kunde inte hitta / skapa filen logs.log \n")
	}
	multi := zerolog.MultiLevelWriter(consoleWriter, runLogFile)
	log := zerolog.New(multi).With().Timestamp().Logger()

	envErr := godotenv.Load()
	if envErr != nil {
		log.Panic().Err(envErr).Msg("Error vid hämtning av .env")
	}

	dbUrl := database.GetDbUrl()

	dbConfig, confParseErr := pgx.ParseConfig(dbUrl)

	if confParseErr != nil {
		log.Panic().Err(confParseErr).Msg("Ogiltig url för databas")
	}

	// TODO: Byt till connection pool. Undersök vidare vad det är
	conn, dbConConfErr := pgx.ConnectConfig(context.Background(), dbConfig)
	if dbConConfErr != nil {
		log.Panic().Err(dbConConfErr).Msg("Kunde inte ansluta till databasen")
	}

	defer conn.Close(context.Background())

	database.TestingQuery(conn)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))

	routes.RegisterRoutesApiStatus(router)
	routes.RegisterRoutesBevakningar(router)

	router.Run(":" + constants.GetPort())
}

// TODO: Fixa
func setUpLogger() {

}

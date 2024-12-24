package main

import (
	"context"
	"fmt"
	"ks-web-scraper/src/constants"
	"ks-web-scraper/src/services"
	"ks-web-scraper/src/types"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	startTime := time.Now()

	// TODO: Kolla varför måste det vara 15:04:05?
	fmt.Fprintf(os.Stderr, "Init api @ \x1b[%dm%s\x1b[0m\n\n", 32, time.Now().Format("15:04:05"))

	// TODO: Det kanske går att bryta ut skapandet av logger till funktion?
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	runLogFile, logFileError := os.OpenFile("logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if logFileError != nil {
		fmt.Fprintf(os.Stderr, "Kunde inte hitta filen error.log \n")
	}
	multi := zerolog.MultiLevelWriter(consoleWriter, runLogFile)
	log := zerolog.New(multi).With().Timestamp().Logger()

	envErr := godotenv.Load()
	if envErr != nil {
		log.Panic().Err(envErr).Msg("Error vid hämtning av .env")
	}

	dbUrl, dbUrlError := getDbUrl()

	if dbUrlError != nil {
		log.Panic().Err(dbUrlError).Msg("Kunde inte skapa database url")
	}

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

	selectQuery := "select id, watch_to_scrape, label, watches, active, last_email_sent, added from watch"

	rows, queryErr := conn.Query(context.Background(), selectQuery)
	if queryErr != nil {
		log.Error().Msg("SQL query för att hämta bevakningar misslyckades: " + queryErr.Error())
	}

	defer rows.Close()

	var watches []types.Watch
	for rows.Next() {
		var w types.Watch
		scanErr := rows.Scan(&w.Id, &w.WatchToScrape, &w.Label, &w.Watches, &w.Active, &w.LastEmailSent, &w.Added)

		if scanErr != nil {
			log.Error().Msg("Kunde inte köra scan av raden: " + scanErr.Error())
			return
		}

		fmt.Fprintf(os.Stderr, "%v\n", w.Label)

		watches = append(watches, w)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))

	// TODO: Gör om api till constant
	router.GET("/api/api-status", func(c *gin.Context) {
		conn, wsError := constants.Upgrader.Upgrade(c.Writer, c.Request, nil)

		if wsError != nil {
			log.Error().Msg("Kunde inte skapa websocket: " + wsError.Error())
			return
		}

		defer conn.Close()

		for {
			status := services.GetApiStatus(startTime)

			err := conn.WriteJSON(status)

			if err != nil {
				return
			}

			time.Sleep(time.Second)
		}
	})

	router.GET("/api/bevakningar/all-watches", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hejsan"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router.Run(":" + port)
}

// TODO: Det är kanske lättare om postgres använder url i dev också. Då slipper man sätta så många variabler
func getDbUrl() (string, error) {
	envHost := os.Getenv("PGHOST")
	envPort := os.Getenv("PGPORT")
	envUsername := os.Getenv("PGUSERNAME")
	envPassword := os.Getenv("PGPASSWORD")
	envDatabase := os.Getenv("PGDATABASE")

	var dbUrl strings.Builder

	if os.Getenv("ENV") != "dev" {
		envDatabaseUrl := os.Getenv("DATABASE_URL")
		dbUrl.WriteString(envDatabaseUrl)
		return dbUrl.String(), nil
	}

	dbUrl.WriteString("user=" + envUsername)
	dbUrl.WriteString(" password=" + envPassword)
	dbUrl.WriteString(" host=" + envHost)
	dbUrl.WriteString(" port=" + envPort)
	dbUrl.WriteString(" dbname=" + envDatabase)

	return dbUrl.String(), nil
}

// TODO: Fixa
func setUpLogger() {

}

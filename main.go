package main

import (
	"context"
	"fmt"
	"ks-web-scraper/src/constants"
	"ks-web-scraper/src/database"
	"ks-web-scraper/src/routes"
	"ks-web-scraper/src/services"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	initApiMsg := "Init api @ \x1b[32m" + time.Now().Format("15:04:05") + "\x1b[0m\n\n" // 32 = grön
	fmt.Fprint(os.Stderr, initApiMsg)

	log := setUpLogger()

	loadDotEnvFile(log)

	conn := setUpDb(log)

	defer conn.Close(context.Background())

	services.ScrapeWatchInfo("")

	// database.GetAllWatches(conn)
	// database.GetAllNotifications(conn)

	router := gin.Default()

	router.Use(constants.CorsConfig)

	routes.RegisterRoutesApiStatus(router)
	routes.RegisterRoutesBevakningar(router, conn)

	routerRunErr := router.Run(constants.GetPort())

	if routerRunErr != nil {
		log.Error().Msg("Kunde inte starta server:" + routerRunErr.Error())
	}
}

func setUpLogger() zerolog.Logger {
	if _, err := os.Stat("logs/logs.log"); os.IsNotExist(err) {
		os.MkdirAll("logs/", 0700)
	}

	runLogFile, logFileError := os.OpenFile("logs/logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if logFileError != nil {
		fmt.Fprintf(os.Stderr, "Kunde inte hitta / skapa filen logs.log \n%v", logFileError)
		defer runLogFile.Close()
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	multi := zerolog.MultiLevelWriter(consoleWriter, runLogFile)
	log := zerolog.New(multi).With().Timestamp().Logger()

	return log
}

func setUpDb(log zerolog.Logger) *pgx.Conn {
	dbConfig, confParseErr := pgx.ParseConfig(database.GetDbUrl())

	if confParseErr != nil {
		log.Panic().Err(confParseErr).Msg("Ogiltig url för databas")
	}

	// TODO: Byt till connection pool. Undersök vidare vad det är
	conn, dbConConfErr := pgx.ConnectConfig(context.Background(), dbConfig)
	if dbConConfErr != nil {
		log.Panic().Err(dbConConfErr).Msg("Kunde inte ansluta till databasen")
	}

	return conn
}

func loadDotEnvFile(log zerolog.Logger) {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Panic().Err(envErr).Msg("Error vid hämtning av .env")
	}
}

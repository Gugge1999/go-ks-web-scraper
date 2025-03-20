package services

import (
	"context"
	"fmt"
	"ks-web-scraper/src/database"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func SetUpLogger() zerolog.Logger {
	_, err := os.Stat("logs/logs.log")

	if os.IsNotExist(err) {
		// Förklaring: rwx | 7 | Read, write and execute för user. Mer info finns här: https://stackoverflow.com/a/31151508
		os.MkdirAll("logs/", 0700)
	}

	runLogFile, logFileError := os.OpenFile("logs/logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if logFileError != nil {
		fmt.Fprintf(os.Stderr, "Kunde inte hitta / skapa filen logs.log\n%v", logFileError)
		defer runLogFile.Close()
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	multi := zerolog.MultiLevelWriter(consoleWriter, runLogFile)
	log := zerolog.New(multi).With().Timestamp().Logger()

	return log
}

func SetUpDb(log zerolog.Logger) *pgx.Conn {
	dbConfig, confParseErr := pgx.ParseConfig(database.GetDbUrl())

	if confParseErr != nil {
		log.Panic().Err(confParseErr).Msg("Ogiltig config för databas")
	}

	// TODO: Byt till connection pool. Undersök vidare vad det är
	conn, dbConConfErr := pgx.ConnectConfig(context.Background(), dbConfig)
	if dbConConfErr != nil {
		log.Panic().Err(dbConConfErr).Msg("Kunde inte ansluta till databasen")
	}

	return conn
}

func LoadDotEnvFile(log zerolog.Logger) {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Panic().Err(envErr).Msg("Error vid hämtning av .env")
	}
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":3000"
	}

	return ":" + port
}

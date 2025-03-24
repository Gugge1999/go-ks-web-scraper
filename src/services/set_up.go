package services

import (
	"context"
	"ks-web-scraper/src/database"
	"ks-web-scraper/src/logger"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func SetUpDb() *pgx.Conn {
	logger := logger.GetLogger()

	dbConfig, confParseErr := pgx.ParseConfig(database.GetDbUrl())

	if confParseErr != nil {
		logger.Panic().Err(confParseErr).Msg("Ogiltig config för databas")
	}

	// TODO: Byt till connection pool. Undersök vidare vad det är
	conn, dbConConfErr := pgx.ConnectConfig(context.Background(), dbConfig)
	if dbConConfErr != nil {
		logger.Panic().Err(dbConConfErr).Msg("Kunde inte ansluta till databasen")
	}

	return conn
}

func LoadDotEnvFile() {
	logger := logger.GetLogger()

	envErr := godotenv.Load()
	if envErr != nil {
		logger.Panic().Err(envErr).Msg("Error vid hämtning av .env")
	}
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":3000"
	}

	return ":" + port
}

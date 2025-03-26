package database

import (
	"context"
	"ks-web-scraper/src/logger"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() *pgxpool.Pool {
	logger := logger.GetLogger()

	dbConfig, confParseErr := pgxpool.ParseConfig(getDbUrl())

	if confParseErr != nil {
		logger.Panic().Err(confParseErr).Msg("Ogiltig config för databas")
	}

	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		logger.Panic().Msg("Kunde inte skapa pgx connection pool: " + err.Error())
	}

	return dbpool
}

func getDbUrl() string {
	if os.Getenv("ENV") != "dev" {
		envDatabaseUrl := os.Getenv("DATABASE_URL")
		return envDatabaseUrl
	}

	return "postgres://" + os.Getenv("PGUSERNAME") + ":" + os.Getenv("PGPASSWORD") + "@localhost:5432/" + os.Getenv("PGDATABASE")
}

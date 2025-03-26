package database

import (
	"context"
	"ks-web-scraper/src/logger"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

// TODO: Ska den använda once do?
func InitDB() {
	logger := logger.GetLogger()

	dbConfig, confParseErr := pgxpool.ParseConfig(getDbUrl())

	if confParseErr != nil {
		logger.Panic().Err(confParseErr).Msg("Ogiltig config för databas")
	}

	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		logger.Panic().Msg("Kunde inte skapa pgx connection pool: " + err.Error())
	}

	db = dbpool
}

func GetDB() *pgxpool.Pool {
	return db
}

func getDbUrl() string {
	if os.Getenv("ENV") != "dev" {
		envDatabaseUrl := os.Getenv("DATABASE_URL")
		return envDatabaseUrl
	}

	return "postgres://" + os.Getenv("PGUSERNAME") + ":" + os.Getenv("PGPASSWORD") + "@localhost:5432/" + os.Getenv("PGDATABASE")
}

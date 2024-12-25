package database

import (
	"os"
)

func GetDbUrl() string {
	if os.Getenv("ENV") != "dev" {
		envDatabaseUrl := os.Getenv("DATABASE_URL")
		return envDatabaseUrl
	}

	envUsername := os.Getenv("PGUSERNAME")
	envPassword := os.Getenv("PGPASSWORD")
	envDatabase := os.Getenv("PGDATABASE")

	return "postgres://" + envUsername + ":" + envPassword + "@localhost:5432/" + envDatabase
}

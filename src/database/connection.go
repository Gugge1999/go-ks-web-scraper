package database

import (
	"os"
)

func GetDbUrl() string {
	if os.Getenv("ENV") != "dev" {
		envDatabaseUrl := os.Getenv("DATABASE_URL")
		return envDatabaseUrl
	}

	return "postgres://" + os.Getenv("PGUSERNAME") + ":" + os.Getenv("PGPASSWORD") + "@localhost:5432/" + os.Getenv("PGDATABASE")
}

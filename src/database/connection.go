package database

import (
	"os"
	"strings"
)

// GetDbUrl TODO: Det är kanske lättare om postgres använder url i dev också. Då slipper man sätta så många variabler
func GetDbUrl() (string, error) {
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

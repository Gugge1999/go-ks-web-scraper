package services

import (
	"ks-web-scraper/src/logger"
	"os"

	"github.com/joho/godotenv"
)

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

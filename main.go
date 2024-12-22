package main

import (
	"context"
	"encoding/json"
	"fmt"
	"ks-web-scraper/constants"
	"ks-web-scraper/middleware"
	"ks-web-scraper/types"
	"net/http"
	"os"
	"runtime"
	"runtime/metrics"
	"strings"
	"time"

	"github.com/gorilla/mux"
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

	conn, dbConErr := pgx.ConnectConfig(context.Background(), dbConfig)
	if dbConErr != nil {
		log.Panic().Err(dbConErr).Msg("Kunde inte ansluta till databasen")
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

	r := mux.NewRouter()

	r.Use(middleware.ContentTypeApplicationJsonMiddleware)

	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	r.HandleFunc("/api-status", func(w http.ResponseWriter, r *http.Request) {
		status := types.ApiStatus{
			Active:                    true,
			ScrapingIntervalInMinutes: constants.IntervalInMin,
			NumberOfCpus:              runtime.NumCPU(),
			MemoryUsage:               getMemoryUsageInMb(),
			Uptime:                    getUptime(startTime),
		}

		err := json.NewEncoder(w).Encode(status)
		if err != nil {
			return
		}
	})

	http.ListenAndServe(":3000", r)
}

// TODO: Den här verkar endast öka med belastning men minskar aldrig
func getMemoryUsageInMb() uint64 {
	const myMetric = "/memory/classes/total:bytes"

	// Create a sample for the metric.
	sample := make([]metrics.Sample, 1)
	sample[0].Name = myMetric

	// Sample the metric.
	metrics.Read(sample)

	bytesInMb := sample[0].Value.Uint64() / 1024 / 1024

	return bytesInMb
}

func getUptime(startTime time.Time) types.Uptime {
	uptime := time.Since(startTime)
	seconds := uint8(uptime.Seconds()) % 60
	minutes := uint8(uptime.Minutes()) % 60
	hours := uint8(uptime.Hours()) % 24
	days := uint16(float64(hours/24)) % 30
	months := uint8(float64(days/30)) % 12
	years := days / 365

	return types.Uptime{
		Seconds: seconds,
		Minutes: minutes,
		Hours:   hours,
		Days:    days,
		Months:  months,
		Years:   years,
	}
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

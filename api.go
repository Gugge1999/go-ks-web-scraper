package main

import (
	"context"
	"encoding/json"
	"fmt"
	"ks-web-scraper/constants"
	"ks-web-scraper/types"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/metrics"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	startTime := time.Now()
	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", envErr)
		os.Exit(1)
	}

	dbUrl, envErr := getDbUrl()

	if envErr != nil {
		os.Exit(1)
	}

	dbConfig, confParseErr := pgx.ParseConfig(dbUrl)

	if confParseErr != nil {
		fmt.Fprintf(os.Stderr, "Invalid url to database: %v\n", confParseErr)
		os.Exit(1)
	}

	conn, dbConErr := pgx.ConnectConfig(context.Background(), dbConfig)
	if dbConErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dbConErr)
		log.Fatalf("unexpected error while tried to connect to database: %v\n", dbConErr)
	}

	defer conn.Close(context.Background())

	selectQuery := "select id, watch_to_scrape, label, watches, active, last_email_sent, added from watch"

	rows, queryErr := conn.Query(context.Background(), selectQuery)
	if queryErr != nil {
		fmt.Fprintf(os.Stderr, "Could not get watches: %v\n", queryErr)
		os.Exit(1)
	}

	defer rows.Close()

	var watches []types.Watch
	for rows.Next() {
		var w types.Watch
		scanErr := rows.Scan(&w.Id, &w.WatchToScrape, &w.Label, &w.Watches, &w.Active, &w.LastEmailSent, &w.Added)

		if scanErr != nil {
			fmt.Fprintf(os.Stderr, "Could not scan row: %v\n", scanErr)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "%v\n", w.Label)

		watches = append(watches, w)
	}

	r := mux.NewRouter()

	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	// TODO: För att använda lowercase i json dto: https://stackoverflow.com/a/11694255/14671400
	r.HandleFunc("/api-status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		//var m runtime.MemStats
		//runtime.ReadMemStats(&m)
		//// For info on each, see: https://golang.org/pkg/runtime/#MemStats
		////fmt.Printf("\nMemory usage = %v MB", byesToMb(m.Sys))
		//fmt.Printf("\nthe host has %d cpus\n", runtime.NumCPU())
		//
		//go heartBeat()
		//time.Sleep(time.Second * 5)
		//
		//sum := 0
		//for i := 0; i < 1_000_000_000; i++ {
		//	sum += i
		//}
		//fmt.Println(sum)

		uptime := time.Since(startTime)
		seconds := uint8(uptime.Seconds()) % 60
		minutes := uint8(uptime.Minutes()) % 60
		hours := uint8(uptime.Hours()) % 24
		days := uint16(float64(hours/24)) % 30
		months := uint8(float64(days/30)) % 12
		years := days / 365

		status := types.ApiStatus{
			Active:                    true,
			ScrapingIntervalInMinutes: constants.IntervalInMin,
			NumberOfCpus:              runtime.NumCPU(),
			MemoryUsage:               getMemoryUsageInMb(),
			Uptime: types.Uptime{
				Seconds: seconds,
				Minutes: minutes,
				Hours:   hours,
				Days:    days,
				Months:  months,
				Years:   years,
			},
		}

		err := json.NewEncoder(w).Encode(status)
		if err != nil {
			return
		}
	})

	// TODO: Kolla varför måste det vara 15:04:05?
	fmt.Fprintf(os.Stderr, "Init API@ %v\n", time.Now().Format("15:04:05"))

	http.ListenAndServe(":3000", r)

}

func getMemoryUsageInMb() uint64 {
	const myMetric = "/gc/heap/allocs:bytes"

	// Create a sample for the metric.
	sample := make([]metrics.Sample, 1)
	sample[0].Name = myMetric

	// Sample the metric.
	metrics.Read(sample)

	return sample[0].Value.Uint64()
}

func byesToMb(bytes uint64) uint64 {
	return bytes / 1024 / 1024
}

func heartBeat() {
	for range time.Tick(time.Second * 1) {
		//fmt.Println("Foo")
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

	dbUrl.WriteString(" user=" + envUsername)
	dbUrl.WriteString(" password=" + envPassword)
	dbUrl.WriteString(" host=" + envHost)
	dbUrl.WriteString(" port=" + envPort)
	dbUrl.WriteString(" dbname=" + envDatabase)

	return dbUrl.String(), nil
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var startTime time.Time

func uptime() time.Duration {
	return time.Since(startTime)
}

func main() {
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

	var watches []watch
	for rows.Next() {
		var w watch
		scanErr := rows.Scan(&w.id, &w.watchToScrape, &w.label, &w.watches, &w.active, &w.lastEmailSent, &w.added)

		if scanErr != nil {
			fmt.Fprintf(os.Stderr, "Could not scan row: %v\n", scanErr)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "%v\n", w.label)

		watches = append(watches, w)
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("\nSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\nthe host has %d cpus\n", runtime.NumCPU())

	r := mux.NewRouter()

	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	r.HandleFunc("/api-status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "uptime %s\n", uptime())
	})

	http.ListenAndServe(":3000", r)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func getDbUrl() (string, error) {
	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", envErr)
		return "", envErr
	}

	envHost := os.Getenv("PGHOST")
	envPort := os.Getenv("PGPORT")
	envUsername := os.Getenv("PGUSERNAME")
	envPassword := os.Getenv("PGPASSWORD")
	envDatabase := os.Getenv("PGDATABASE")
	envDatabaseUrl := os.Getenv("DATABASE_URL")

	var dbUrl strings.Builder

	if envDatabaseUrl != "" {
		dbUrl.WriteString(envDatabaseUrl) // Url för prod
		return dbUrl.String(), nil
	}

	dbUrl.WriteString(" user=" + envUsername)
	dbUrl.WriteString(" password=" + envPassword)
	dbUrl.WriteString(" host=" + envHost)
	dbUrl.WriteString(" port=" + envPort)
	dbUrl.WriteString(" dbname=" + envDatabase)

	return dbUrl.String(), nil
}

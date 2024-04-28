package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
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

func init() {
	startTime = time.Now()
}

type watch struct {
	id            string // TODO: Hur ska guid hanteras?
	watchToScrape string
	label         string
	watches       string
	active        bool
	lastEmailSent sql.NullTime
	added         time.Time
}

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", envErr)
		os.Exit(1)
	}

	envHost := os.Getenv("PGHOST")
	envPort := os.Getenv("PGPORT")
	envUsername := os.Getenv("PGUSERNAME")
	envPassword := os.Getenv("PGPASSWORD")
	envDatabase := os.Getenv("PGDATABASE")

	var dbUrl strings.Builder
	dbUrl.WriteString(" user=" + envUsername)
	dbUrl.WriteString(" password=" + envPassword)
	dbUrl.WriteString(" host=" + envHost)
	dbUrl.WriteString(" port=" + envPort)
	dbUrl.WriteString(" dbname=" + envDatabase)

	dbConfig, confParseErr := pgx.ParseConfig(dbUrl.String())

	if confParseErr != nil {
		fmt.Fprintf(os.Stderr, "Invalid dbUrl: %v\n", confParseErr)
		os.Exit(1)
	}

	conn, dbConErr := pgx.ConnectConfig(context.Background(), dbConfig)
	if dbConErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dbConErr)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	selectQuery := "select id, \"watchToScrape\", label, watches, active, \"lastEmailSent\", added from watch"

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

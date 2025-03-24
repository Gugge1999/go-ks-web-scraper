package database

import (
	"context"
	"encoding/json"
	"ks-web-scraper/src/types"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func GetAllWatches(conn *pgx.Conn) []types.Watch {
	dbQuery := `select id, watch_to_scrape, label, watches, active, last_email_sent, added from watch`

	rows, queryErr := conn.Query(context.Background(), dbQuery)
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
			return nil
		}

		var scrapedWatches []types.ScrapedWatch
		unmarshalErr := json.Unmarshal([]byte(w.Watches), &scrapedWatches)
		if unmarshalErr != nil {
			log.Error().Msg("Kunde inte göra unmarshal av watches. Error:" + unmarshalErr.Error())
		}

		w.ScrapedWatches = scrapedWatches

		watches = append(watches, w)
	}

	return watches
}

func SaveWatch(conn *pgx.Conn, saveWatchDto types.SaveWatchDto, scrapedWatches []types.ScrapedWatch) []types.Watch {
	dbQuery := `INSERT INTO watch (label, watch_to_scrape, active, watches) VALUES (@label, @watchToScrape, @active, @scrapedWatches) RETURNING *`
	args := pgx.NamedArgs{
		"label":          saveWatchDto.Label,
		"watchToScrape":  saveWatchDto.WatchToScrape,
		"active":         true,
		"scrapedWatches": scrapedWatches,
	}

	rows, err := conn.Query(context.Background(), dbQuery, args)
	// TODO: Logger verkar inte fungera
	if err != nil {
		log.Error().Msg("Kunde inte spara ny bevakning. Error:" + err.Error())
	}

	var watches []types.Watch
	defer rows.Close()
	for rows.Next() {
		var w types.Watch
		scanErr := rows.Scan(&w.Id, &w.Label, &w.Watches, &w.Active, &w.WatchToScrape, &w.LastEmailSent, &w.Added)

		if scanErr != nil {
			log.Error().Msg("Kunde inte köra scan av raden: " + scanErr.Error())
			return nil
		}

		var scrapedWatches []types.ScrapedWatch
		unmarshalErr := json.Unmarshal([]byte(w.Watches), &scrapedWatches)
		if unmarshalErr != nil {
			log.Error().Msg("Kunde inte göra unmarshal av watches. Error:" + unmarshalErr.Error())
		}

		w.ScrapedWatches = scrapedWatches

		watches = append(watches, w)
	}

	return watches
}

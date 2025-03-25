package database

import (
	"context"
	"encoding/json"
	"ks-web-scraper/src/logger"
	"ks-web-scraper/src/types"

	"github.com/jackc/pgx/v5"
)

func GetAllWatches(conn *pgx.Conn) []types.Watch {
	logger := logger.GetLogger()

	const dbQuery = `SELECT id, watch_to_scrape, label, watches, active, last_email_sent, added 
						FROM watch
							ORDER BY added`

	rows, queryErr := conn.Query(context.Background(), dbQuery)
	if queryErr != nil {
		logger.Error().Msg("SQL query för att hämta bevakningar misslyckades: " + queryErr.Error())
	}

	defer rows.Close()

	var watches []types.Watch
	for rows.Next() {
		var w types.Watch
		scanErr := rows.Scan(&w.Id, &w.WatchToScrape, &w.Label, &w.Watches, &w.Active, &w.LastEmailSent, &w.Added)

		if scanErr != nil {
			logger.Error().Msg("Kunde inte köra scan av raden: " + scanErr.Error())
			return nil
		}

		var scrapedWatches []types.ScrapedWatch
		unmarshalErr := json.Unmarshal([]byte(w.Watches), &scrapedWatches)
		if unmarshalErr != nil {
			logger.Error().Msg("Kunde inte göra unmarshal av watches. Error:" + unmarshalErr.Error())
		}

		w.ScrapedWatches = scrapedWatches

		watches = append(watches, w)
	}

	return watches
}

func SaveWatch(conn *pgx.Conn, saveWatchDto types.SaveWatchDto, scrapedWatches []types.ScrapedWatch) []types.Watch {
	logger := logger.GetLogger()

	const dbQuery = `INSERT INTO watch (label, watch_to_scrape, active, watches)
						VALUES
	 						(@label, @watchToScrape, @active, @scrapedWatches)
								RETURNING *`
	args := pgx.NamedArgs{
		"label":          saveWatchDto.Label,
		"watchToScrape":  saveWatchDto.WatchToScrape,
		"active":         true,
		"scrapedWatches": scrapedWatches,
	}

	rows, err := conn.Query(context.Background(), dbQuery, args)
	if err != nil {
		logger.Error().Msg("Kunde inte spara ny bevakning. Error:" + err.Error())
	}

	var watches []types.Watch
	// TODO: Kolla om det går att göra en egen funktion för att skapa en watch
	defer rows.Close()
	for rows.Next() {
		var w types.Watch
		scanErr := rows.Scan(&w.Id, &w.Label, &w.Watches, &w.Active, &w.WatchToScrape, &w.LastEmailSent, &w.Added)

		if scanErr != nil {
			logger.Error().Msg("Kunde inte köra scan av raden: " + scanErr.Error())
			return nil
		}

		var scrapedWatches []types.ScrapedWatch
		unmarshalErr := json.Unmarshal([]byte(w.Watches), &scrapedWatches)
		if unmarshalErr != nil {
			logger.Error().Msg("Kunde inte göra unmarshal av watches. Error:" + unmarshalErr.Error())
		}

		w.ScrapedWatches = scrapedWatches

		watches = append(watches, w)
	}

	return watches
}

func DeleteWatch(conn *pgx.Conn, watchId string) (string, error) {
	logger := logger.GetLogger()

	const dbQuery = `DELETE FROM watch
						WHERE id = @watchId
							RETURNING *`

	args := pgx.NamedArgs{
		"watchId": watchId,
	}

	rows, err := conn.Query(context.Background(), dbQuery, args)
	if err != nil {
		logger.Error().Msg("Kunde inte radera bevakning. Error:" + err.Error())
		return "", err
	}

	defer rows.Close()

	return watchId, nil
}

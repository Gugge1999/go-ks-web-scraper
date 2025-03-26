package database

import (
	"context"
	"encoding/json"
	"ks-web-scraper/src/logger"
	"ks-web-scraper/src/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func GetAllWatches(dbPoolConn *pgxpool.Pool) ([]types.Watch, error) {
	logger := logger.GetLogger()

	const dbQuery = `SELECT id, label, watches, active, watch_to_scrape, last_email_sent, added 
						FROM watch
							ORDER BY added`

	rows, queryErr := dbPoolConn.Query(context.Background(), dbQuery)
	if queryErr != nil {
		logger.Error().Msg("SQL query för att hämta bevakningar misslyckades: " + queryErr.Error())
		return nil, queryErr
	}

	watches := mapDbRowToWatchDto(rows, logger)

	return watches, nil
}

func SaveWatch(dbPoolConn *pgxpool.Pool, label string, watchToScrapeUrl string, scrapedWatches []types.ScrapedWatch) ([]types.Watch, error) {
	logger := logger.GetLogger()

	const dbQuery = `INSERT INTO watch (label, watches, active, watch_to_scrape)
						VALUES
	 						(@label, @scrapedWatches, @active, @watchToScrape)
								RETURNING *`
	args := pgx.NamedArgs{
		"label":          label,
		"scrapedWatches": scrapedWatches,
		"active":         true,
		"watchToScrape":  watchToScrapeUrl,
	}

	rows, err := dbPoolConn.Query(context.Background(), dbQuery, args)
	if err != nil {
		logger.Error().Msg("Kunde inte spara ny bevakning. Error:" + err.Error())
		return nil, err
	}

	watches := mapDbRowToWatchDto(rows, logger)

	return watches, nil
}

func DeleteWatch(dbPoolConn *pgxpool.Pool, watchId string) (string, error) {
	logger := logger.GetLogger()

	const dbQuery = `DELETE FROM watch
						WHERE id = @watchId
							RETURNING *`

	args := pgx.NamedArgs{
		"watchId": watchId,
	}

	rows, err := dbPoolConn.Query(context.Background(), dbQuery, args)
	if err != nil {
		logger.Error().Msg("Kunde inte radera bevakning. Error:" + err.Error())
		return "", err
	}

	defer rows.Close()

	return watchId, nil
}

func ToggleActiveStatuses(dbPoolConn *pgxpool.Pool, ids []string, newActiveStatus bool) ([]types.Watch, error) {
	logger := logger.GetLogger()

	const dbQuery = `UPDATE watch
						SET active = @newActiveStatus
							WHERE id = ANY (@ids)
								RETURNING *`

	args := pgx.NamedArgs{
		"newActiveStatus": newActiveStatus,
		"ids":             ids,
	}

	rows, err := dbPoolConn.Query(context.Background(), dbQuery, args)
	if err != nil {
		logger.Error().Msg("Kunde inte ändra aktiv status. Error:" + err.Error())

		return nil, err
	}

	watches := mapDbRowToWatchDto(rows, logger)

	return watches, nil
}

func mapDbRowToWatchDto(rows pgx.Rows, logger zerolog.Logger) []types.Watch {
	defer rows.Close()

	var watches []types.Watch
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

		w.LatestWatch = scrapedWatches[0]

		watches = append(watches, w)
	}

	return watches
}

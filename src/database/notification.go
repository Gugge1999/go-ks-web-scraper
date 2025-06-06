package database

import (
	"context"
	"ks-web-scraper/src/logger"
	"ks-web-scraper/src/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllNotifications(dbPoolConn *pgxpool.Pool) ([]types.Notification, error) {
	logger := logger.GetLogger()
	const dbQuery = `SELECT * 
						FROM NOTIFICATION
							ORDER BY sent ASC`

	rows, queryErr := dbPoolConn.Query(context.Background(), dbQuery)

	if queryErr != nil {
		logger.Error().Msg("SQL query för att hämta notiser misslyckades: " + queryErr.Error())
		return nil, queryErr
	}

	return getNotificationRows(rows)
}

func InsertNewNotification(dbPoolConn *pgxpool.Pool, watchId string) ([]types.Notification, error) {
	logger := logger.GetLogger()

	const dbQuery = `INSERT INTO notification(watch_id)
							VALUES ($1)
								RETURNING *`

	args := pgx.NamedArgs{
		"watchId": watchId,
	}

	rows, queryErr := dbPoolConn.Query(context.Background(), dbQuery, args)

	if queryErr != nil {
		logger.Error().Msg("SQL query för att skapa ny notification misslyckades: " + queryErr.Error())
		return nil, queryErr
	}

	return getNotificationRows(rows)
}

func getNotificationRows(rows pgx.Rows) ([]types.Notification, error) {
	logger := logger.GetLogger()

	defer rows.Close()

	var notifications []types.Notification

	for rows.Next() {
		var n types.Notification
		scanErr := rows.Scan(&n.Id, &n.WatchId, &n.Sent)

		if scanErr != nil {
			logger.Error().Msg("Kunde inte hämta köra scan av notiser. Error:" + scanErr.Error())
			return nil, scanErr
		}

		notifications = append(notifications, n)
	}

	return notifications, nil
}

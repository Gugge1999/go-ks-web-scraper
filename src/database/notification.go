package database

import (
	"context"
	"ks-web-scraper/src/types"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func GetAllNotifications(conn *pgx.Conn) ([]types.Notification, error) {
	selectQuery := `SELECT * FROM NOTIFICATION`

	rows, queryErr := conn.Query(context.Background(), selectQuery)

	if queryErr != nil {
		log.Error().Msg("SQL query för att hämta notiser misslyckades: " + queryErr.Error())
		return nil, queryErr
	}

	return getNotificationRows(rows)
}

// TODO: Kolla på https://hexacluster.ai/postgresql/connecting-to-postgresql-with-go-using-pgx/
func InsertNewNotification(conn *pgx.Conn, watchId string) ([]types.Notification, error) {
	insertQuery := `INSERT INTO notification(watch_id) VALUES ($1) RETURNING *`

	args := pgx.NamedArgs{
		"watchId": watchId,
	}

	rows, queryErr := conn.Query(context.Background(), insertQuery, args)

	if queryErr != nil {
		log.Error().Msg("SQL query för att skapa ny notification misslyckades: " + queryErr.Error())
		return nil, queryErr
	}

	defer rows.Close()

	return getNotificationRows(rows)
}

func getNotificationRows(rows pgx.Rows) ([]types.Notification, error) {
	defer rows.Close()

	var notifications []types.Notification

	for rows.Next() {
		var n types.Notification
		scanErr := rows.Scan(&n.Id, &n.WatchId, &n.Sent)

		if scanErr != nil {
			log.Error().Msg("Kunde inte hämta köra scan av notiser. Error:" + scanErr.Error())
			return nil, scanErr
		}

		notifications = append(notifications, n)
	}

	return notifications, nil
}

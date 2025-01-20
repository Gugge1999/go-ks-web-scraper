package database

import (
	"context"
	"ks-web-scraper/src/types"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func GetAllNotifications(conn *pgx.Conn) []types.Notification {
	selectQuery := "select id, watch_id, sent from notification"
	rows, queryErr := conn.Query(context.Background(), selectQuery)
	if queryErr != nil {
		log.Error().Msg("SQL query för att hämta notiser misslyckades: " + queryErr.Error())
	}

	defer rows.Close()

	var notifications []types.Notification
	for rows.Next() {
		var n types.Notification
		scanErr := rows.Scan(&n.Id, &n.WatchId, &n.Sent)

		if scanErr != nil {
			log.Error().Msg("Kunde inte hämta köra scan av notiser. Error:" + scanErr.Error())
		}

		notifications = append(notifications, n)
	}

	return notifications
}

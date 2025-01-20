package database

import (
	"context"
	"encoding/json"
	"fmt"
	"ks-web-scraper/src/types"
	"os"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type TheParam struct {
	Id       string `json:"id"`
	Provider string `json:"provider"`
}

func GetAllWatches(conn *pgx.Conn) []types.Watch {
	selectQuery := "select id, watch_to_scrape, label, watches::jsonb, active, last_email_sent, added from watch"

	rows, queryErr := conn.Query(context.Background(), selectQuery)
	if queryErr != nil {
		log.Error().Msg("SQL query för att hämta bevakningar misslyckades: " + queryErr.Error())
	}

	defer rows.Close()

	var watches []types.Watch
	for rows.Next() {
		var w types.Watch
		// TODO: Läs den https://tillitsdone.com/blogs/pgx-with-postgresql-json-in-go/
		scanErr := rows.Scan(&w.Id, &w.WatchToScrape, &w.Label, &w.Watches, &w.Active, &w.LastEmailSent, &w.Added)

		fmt.Fprintf(os.Stderr, "type: %v\n", reflect.TypeOf(&w.Watches))

		testing := "{\"id\":\"some-id\",\"provider\":\"any-provider\"}"
		param := TheParam{}
		err := json.Unmarshal([]byte(testing), &param)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var vette []types.ScrapedWatch
		const hejsan = "[{\"name\":\"hejsan\",\"postedDate\":\"2021-09-01T00:00:00Z\",\"link\":\"https://hejsan.se\"}]"
		const hejsanDb = "[{\"name\":\"Gummiband och länk till Explorer 1 39mm\",\"postedDate\":\"2024-12-25T17:16:39+0100\",\"link\":\"https://klocksnack.se/threads/gummiband-och-l%C3%A4nk-till-explorer-1-39mm.195113/\"},{\"name\":\"Rolex Milgauss 116400 40mm\",\"postedDate\":\"2024-12-24T12:42:26+0100\",\"link\":\"https://klocksnack.se/threads/rolex-milgauss-116400-40mm.195092/\"},{\"name\":\"Rolex submariner 126613LN  116613LN\",\"postedDate\":\"2024-12-24T10:54:15+0100\",\"link\":\"https://klocksnack.se/threads/rolex-submariner-126613ln-116613ln.195089/\"},{\"name\":\"Omega Aqua Terra  Rolex ExplorerOP\",\"postedDate\":\"2024-12-24T07:46:16+0100\",\"link\":\"https://klocksnack.se/threads/omega-aqua-terra-rolex-explorer-op.195087/\"},{\"name\":\"Rolex Submariner 126610LV\",\"postedDate\":\"2024-12-23T21:34:57+0100\",\"link\":\"https://klocksnack.se/threads/rolex-submariner-126610lv.195082/\"},{\"name\":\"Rolex Submariner 14060M\",\"postedDate\":\"2024-12-23T21:25:25+0100\",\"link\":\"https://klocksnack.se/threads/rolex-submariner-14060m.195080/\"},{\"name\":\"Rolex Oyster Perpetual 36 126000\",\"postedDate\":\"2024-12-23T16:51:10+0100\",\"link\":\"https://klocksnack.se/threads/rolex-oyster-perpetual-36-126000.195073/\"},{\"name\":\"Rolex Lady Datejust 6917\",\"postedDate\":\"2024-12-23T15:36:59+0100\",\"link\":\"https://klocksnack.se/threads/rolex-lady-datejust-6917.195072/\"},{\"name\":\"Rolex Datejust 36mm\",\"postedDate\":\"2024-12-23T09:43:12+0100\",\"link\":\"https://klocksnack.se/threads/rolex-datejust-36mm.195062/\"},{\"name\":\"Rolex Datejust 36mm\",\"postedDate\":\"2024-12-23T08:50:11+0100\",\"link\":\"https://klocksnack.se/threads/rolex-datejust-36mm.195060/\"},{\"name\":\"Rolex Datejust + GMT 126710\",\"postedDate\":\"2024-12-22T19:21:05+0100\",\"link\":\"https://klocksnack.se/threads/rolex-datejust-gmt-126710.195055/\"},{\"name\":\"Rolex Sea-Dweller 116600\",\"postedDate\":\"2024-12-22T15:57:17+0100\",\"link\":\"https://klocksnack.se/threads/rolex-sea-dweller-116600.195045/\"},{\"name\":\"Bytes Garmin Fenix mot RolexOmega accesoarer\",\"postedDate\":\"2024-12-22T13:11:38+0100\",\"link\":\"https://klocksnack.se/threads/garmin-fenix-mot-rolex-omega-accesoarer.195040/\"},{\"name\":\"Rolex Datejust 69173 Tiffany&co\",\"postedDate\":\"2024-12-21T19:27:31+0100\",\"link\":\"https://klocksnack.se/threads/rolex-datejust-69173-tiffany-co.195030/\"},{\"name\":\"Rolex Submariner 114060\",\"postedDate\":\"2024-12-21T17:20:46+0100\",\"link\":\"https://klocksnack.se/threads/rolex-submariner-114060.195028/\"},{\"name\":\"Rolex Submariner 16610LV\",\"postedDate\":\"2024-12-21T17:13:46+0100\",\"link\":\"https://klocksnack.se/threads/rolex-submariner-16610lv.195027/\"},{\"name\":\"Rolex Datejust 279173 28mm\",\"postedDate\":\"2024-12-21T15:33:27+0100\",\"link\":\"https://klocksnack.se/threads/rolex-datejust-279173-28mm.195024/\"},{\"name\":\"Rolex Datejust 41 Rhodium - 126334\",\"postedDate\":\"2024-12-21T11:08:55+0100\",\"link\":\"https://klocksnack.se/threads/rolex-datejust-41-rhodium-126334.195019/\"},{\"name\":\"Rolex OP 114200 0live Green\",\"postedDate\":\"2024-12-21T09:55:30+0100\",\"link\":\"https://klocksnack.se/threads/rolex-op-114200-0live-green.195016/\"},{\"name\":\"Rolex 16700 PB\",\"postedDate\":\"2024-12-19T13:51:04+0100\",\"link\":\"https://klocksnack.se/threads/rolex-16700-pb.194980/\"},{\"name\":\"18K dressareRolex 16805513\",\"postedDate\":\"2024-12-19T11:08:40+0100\",\"link\":\"https://klocksnack.se/threads/18k-dressare-rolex-1680-5513.194975/\"},{\"name\":\"Rolex 116520, VTNR eller BLRO\",\"postedDate\":\"2024-12-19T08:13:45+0100\",\"link\":\"https://klocksnack.se/threads/rolex-116520-vtnr-eller-blro.194968/\"},{\"name\":\"Rolex GMT-Master II 16710\",\"postedDate\":\"2024-12-17T19:55:19+0100\",\"link\":\"https://klocksnack.se/threads/rolex-gmt-master-ii-16710.194946/\"},{\"name\":\"Rolex Datejust 41 - 126334\",\"postedDate\":\"2024-12-17T11:16:32+0100\",\"link\":\"https://klocksnack.se/threads/rolex-datejust-41-126334.194937/\"},{\"name\":\"Rolex 126000, 17013, Turn O Graph  Breitling Navitimer B01  Zenith El Primero  Heuer 2447  Cartier Tank  Must de\",\"postedDate\":\"2024-12-17T09:02:05+0100\",\"link\":\"https://klocksnack.se/threads/rolex-126000-17013-turn-o-graph-breitling-navitimer-b01-zenith-el-primero-heuer-2447-cartier-tank-must-de.194934/\"},{\"name\":\"Rolex Day-Date 118138\",\"postedDate\":\"2024-12-16T20:10:26+0100\",\"link\":\"https://klocksnack.se/threads/rolex-day-date-118138.194930/\"},{\"name\":\"Rolex Sea-Dweller Deepsea 116660\",\"postedDate\":\"2024-12-16T10:01:43+0100\",\"link\":\"https://klocksnack.se/threads/rolex-sea-dweller-deepsea-116660.194915/\"},{\"name\":\"Rolex travel case!\",\"postedDate\":\"2024-12-16T01:49:59+0100\",\"link\":\"https://klocksnack.se/threads/rolex-travel-case.194911/\"},{\"name\":\"Rolex Submariner 14060M\",\"postedDate\":\"2024-12-14T13:32:42+0100\",\"link\":\"https://klocksnack.se/threads/rolex-submariner-14060m.194858/\"},{\"name\":\"Rolex Yachtmaster - 126622\",\"postedDate\":\"2024-12-13T23:14:00+0100\",\"link\":\"https://klocksnack.se/threads/rolex-yachtmaster-126622.194844/\"}]"

		unmarshalErr := json.Unmarshal([]byte(hejsanDb), &vette)
		if unmarshalErr != nil {
			log.Error().Msg("Kunde inte unmarshal av watches. Error:" + unmarshalErr.Error())
		}

		var vetteNy []types.ScrapedWatch
		const hejsanNy = "[{\"name\":\"hejsan\",\"postedDate\":\"2021-09-01T00:00:00Z\",\"link\":\"https://hejsan.se\"}]"
		//r := []rune(w.Watches)

		//var withoutLastChar = string(r[:len(r)-1])
		//var vetteDb = withoutLastChar[1:]

		// TODO: Läs https://stackoverflow.com/a/65434116/14671400
		unmarshalErrNy := json.Unmarshal([]byte(w.Watches), &vetteNy)
		if unmarshalErrNy != nil {
			log.Error().Msg("Kunde inte unmarshal av watches. Error:" + unmarshalErrNy.Error())
		}

		fmt.Fprintf(os.Stderr, "vette: %v\n", w.Watches)

		if scanErr != nil {
			log.Error().Msg("Kunde inte köra scan av raden: " + scanErr.Error())
			return nil
		}

		fmt.Fprintf(os.Stderr, "%v\n", w.Label)

		watches = append(watches, w)
	}

	return watches
}

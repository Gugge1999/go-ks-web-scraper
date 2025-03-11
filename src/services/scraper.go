package services

import (
	"ks-web-scraper/src/types"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

func ScrapeWatchInfo(watchToScrape string) []types.ScrapedWatch {
	var names []string
	var links []string
	var postedDates []string
	var scrapedWatches []types.ScrapedWatch

	c := colly.NewCollector()

	c.OnHTML(".contentRow-title > a", func(e *colly.HTMLElement) {
		annonsLink := "https://klocksnack.se" + e.Attr("href")
		links = append(links, annonsLink)

		annonsTextArr := strings.SplitAfter(e.Text, "\u00a0")
		lastElement := annonsTextArr[len(annonsTextArr)-1]
		names = append(names, lastElement)
	})

	c.OnHTML(".u-dt", func(e *colly.HTMLElement) {
		dataTime := e.Attr("data-time")
		unixTimestamp, errParseInt := strconv.ParseInt(dataTime, 10, 64)

		if errParseInt != nil {
			log.Error().Msg("Kunde inte skapa UNIX timestamp från data-time. Error: " + errParseInt.Error())
		}

		// OBS! Måste vara 2006-01-02T15:04:05-0700 för ISO 8601
		annonsDate := time.Unix(unixTimestamp, 0).UTC().Format("2006-01-02T15:04:05-0700")

		postedDates = append(postedDates, annonsDate)
	})

	c.Visit(watchToScrape)

	// Bygger på att alla 3 array:er är av samma storlek
	for index := range names {
		scrapedWatch := types.ScrapedWatch{
			Name:       names[index],
			PostedDate: postedDates[index],
			Link:       links[index],
		}

		scrapedWatches = append(scrapedWatches, scrapedWatch)
	}

	return scrapedWatches
}

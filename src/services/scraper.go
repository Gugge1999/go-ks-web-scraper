package services

import (
	"fmt"
	"ks-web-scraper/src/types"

	"github.com/gocolly/colly"
)

func ScrapeWatchInfo(watchToScrape string) []types.ScrapedWatch {
	// var titles []string
	// var dates []string
	// var links []string
	var scrapedWatches []types.ScrapedWatch

	c := colly.NewCollector()

	fmt.Printf("Här")

	// On every a element which has href attribute call callback
	c.OnHTML(".contentRow-title > a", func(e *colly.HTMLElement) {
		link := e.Text
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://www.klocksnack.se/search/1/?q=Omega&t=post&c[child_nodes]=1&c[nodes][0]=11&c[title_only]=1&o=date")

	return scrapedWatches
}

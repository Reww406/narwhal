package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/reww406/narwhal/config"
)

// https://www.espn.com/nba/player/_/id/3908845/john-collins


var playerStatsRegex = regexp.MustCompile(`https://www\.espn\.com/nba/player/_/id/(\d{7})/([a-z]\-]+)`)
var gameLogRegex = regexp.MustCompile(`https://www\.espn\.com/nba/player/gamelog/_/id/(\d{7})/([a-z]\-]+)`)

func isGameOrStatsLink(link string) bool {
  // There will be more playerStats so compare it first.
  return playerStatsRegex.Match([]byte(link)) || gameLogRegex.Match([]byte(link))
}


func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.espn.com"),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       5 * time.Second,
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

    if isGameOrStatsLink(link) {  
      e.Request.Visit(link)
    }

		e.Request.Visit(link)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

  // TODO Check if this is a game log page and then scrape the tabs from it if it is.
  c.OnHtml("div[id='fittPageContainer']", func(e *colly.HTMLElement) {})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error", r.Request.URL, err)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	// Team -> Players -> Game log -> Scrape

  c.Visit(config.TeamUrls[0])
	c.Wait()
}

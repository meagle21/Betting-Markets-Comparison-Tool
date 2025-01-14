package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
)

var sportsbook_url string = "https://sportsbook.draftkings.com/leagues/basketball/nba"

func main() {
	fmt.Println("Hello, World!")
	c := colly.NewCollector()
	c.Visit(sportsbook_url)
}

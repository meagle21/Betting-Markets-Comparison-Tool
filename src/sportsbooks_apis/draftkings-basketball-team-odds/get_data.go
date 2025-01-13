package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type structTeamOddsLines struct {
	Team       string
	Gameday    string
	Gametime   string
	Spread     string
	SpreadOdds string
	totalLine  string
	totalOdds  string
	Moneyline  string
}

type Game struct {
	Home *structTeamOddsLines
	Away *structTeamOddsLines
}

func main() {

	var games []Game

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited:", r.Request.URL)
		fmt.Println("Status Code:", r.StatusCode)
	})
	var homeAwayCheck int = 0
	game := &Game{}
	c.OnHTML("tbody.sportsbook-table__body > tr", func(e *colly.HTMLElement) { //returns HTML that matches the selected HTML tag info
		var gameDayFormat string = "01-02-2006"
		var gameTimeFormat string = "3:04PM"
		teamOddsLines := &structTeamOddsLines{} //this creates a pointer for the struct that allows it to be dynamically updated
		teamOddsLines.Gameday = time.Now().Format(gameDayFormat)
		homeAwayCheck = homeAwayCheck + 1

		e.ForEach("td, th", func(i int, el *colly.HTMLElement) { //for each element in the HTML that matches "td" or "th" tag
			cellNumber := i + 1                                           //need to count cell numbers as each cell number contains different info to be parsed differently
			cellText := el.Text                                           //extract cell text from the HTML
			cleanedText := strings.TrimSpace(cellText)                    // String cleaning that needs to occur generally so the if statements work properly
			normalizedText := strings.ReplaceAll(cleanedText, "–", "-")   // Replace en dash
			normalizedText = strings.ReplaceAll(normalizedText, "—", "-") // Replace em dash
			normalizedText = strings.ReplaceAll(normalizedText, "−", "-") // Replace minus sign
			normalizedText = strings.ReplaceAll(normalizedText, "＋", "+") // Fullwidth plus
			normalizedText = strings.ReplaceAll(normalizedText, "±", "+") // Plus-minus (if treated as a simple "+")
			if cellNumber == 1 {                                          // The first cell contains gametime and team name information
				gametimeTeamNameCellSplit := strings.Split(cellText, ":") //split on the colon included in the time
				if len(gametimeTeamNameCellSplit[0]) == 2 {               //this accounts for when a time is 1X:XX
					parsedTime, err := time.Parse(gameTimeFormat, cellText[:7]) //parse the string with the referenced datetime format
					if err != nil {
						fmt.Println("Error parsing time:", err)
					}
					teamOddsLines.Gametime = parsedTime.Add(19 * time.Hour).Format(gameTimeFormat) // add 19 hours to the time, it seems the time returned from the website is 19 hours ahead of game time
					teamOddsLines.Team = strings.TrimSpace(cellText[7:])                           //index to get the team name
				}
				if len(gametimeTeamNameCellSplit[0]) == 1 { //this accouts for gametimes that are X:XX
					parsedTime, err := time.Parse(gameTimeFormat, cellText[:6]) //parse the string with the referenced datetime format
					if err != nil {
						fmt.Println("Error parsing time:", err)
					}
					teamOddsLines.Gametime = parsedTime.Add(19 * time.Hour).Format(gameTimeFormat) // add 19 hours to the time, it seems the time returned from the website is 19 hours behind game time
					teamOddsLines.Team = strings.TrimSpace(cellText[6:])                           //index to get the team name
				}
			}
			if cellNumber == 2 { //get the second cell which is the spread
				textForOddsDiscovery := normalizedText[1:]       //extract the sign that is setting the spread
				if strings.Contains(textForOddsDiscovery, "-") { //if without the spread sign there's a minus sign
					teamOddsLines.SpreadOdds = textForOddsDiscovery[strings.Index(textForOddsDiscovery, "-"):] //find index of minus and get whatever is after it, this will be the odds
				}
				if strings.Contains(textForOddsDiscovery, "+") { //if without the spread sign there's a plus sign
					teamOddsLines.SpreadOdds = textForOddsDiscovery[strings.Index(textForOddsDiscovery, "+"):] //find index of plus and get whatever is after it, this will be the odds
				}
				if strings.Contains(normalizedText, ".") {
					teamOddsLines.Spread = normalizedText[:4]
				}
				if strings.Contains(normalizedText, ".") == false {
					teamOddsLines.Spread = normalizedText[:3]
				}

			}
			if cellNumber == 3 { //get the third cell which is the over under odds
				if strings.Contains(normalizedText, "-") { // Check if a minus exists
					totalCellSplit := strings.Split(normalizedText, "-")                 // If a minus exists, split on it to get the odds and the line as separate variables
					teamOddsLines.totalLine = strings.TrimSpace(totalCellSplit[0])       //add
					teamOddsLines.totalOdds = "-" + strings.TrimSpace(totalCellSplit[1]) //need to prepend here as we lose the sign when we split
				}
				if strings.Contains(normalizedText, "+") { // Check if a plus exists
					totalCellSplit := strings.Split(normalizedText, "+") //Same logic as above for the minus sign
					teamOddsLines.totalLine = strings.TrimSpace(totalCellSplit[0])
					teamOddsLines.totalOdds = "+" + strings.TrimSpace(totalCellSplit[1])
				}
			}
			if cellNumber == 4 {
				teamOddsLines.Moneyline = strings.TrimSpace(cellText)
			}
		})
		if homeAwayCheck%2 == 0 { //even numbered teams are always the home team
			game.Home = teamOddsLines
		}
		if homeAwayCheck%2 != 0 { //odd numbered teams are always the away team
			game.Away = teamOddsLines
		}
		if game.Home != nil && game.Away != nil { // Append only when both teams are set
			games = append(games, *game)
			game = &Game{} // Reset game for the next match
		}

	})

	var sportsbookURL string = "https://sportsbook.draftkings.com/leagues/football/nba"

	c.Visit(sportsbookURL)

	for _, game := range games {
		fmt.Printf("Game:\n  Home: %+v\n  Away: %+v\n", game.Home, game.Away)
	}
}

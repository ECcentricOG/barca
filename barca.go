package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/go-resty/resty/v2"
)

type Score struct {
	HomeTeam  int `json:"homeTeam"`
	AwayTeam  int `json:"awayTeam"`
}

type Team struct {
	Name string `json:"name"`
}

type Match struct {
	Status  string  `json:"status"`
	HomeTeam Team   `json:"homeTeam"`
	AwayTeam Team   `json:"awayTeam"`
	Score   Score  `json:"score"`
	UtcDate string `json:"utcDate"`
}

type MatchResponse struct {
	Matches []Match `json:"matches"`
}

const apiURL = "https://api.football-data.org/v4/teams/81/matches"
const apiKey = "Enter API Key" // Replace with your football-data.org API key

func getFCBarcelonaMatches() {
	client := resty.New()

	resp, err := client.R().
		SetHeader("X-Auth-Token", apiKey).
		Get(apiURL)

	if err != nil {
		log.Fatalf("Error fetching data: %v\n", err)
	}

	if resp.IsError() {
		log.Fatalf("Error: %v\n", resp.Status())
	}

	var data MatchResponse
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		log.Fatalf("Error unmarshaling response: %v\n", err)
	}

	var upcomingMatches []Match
	today := time.Now()

	for _, match := range data.Matches {
		parsedDate, err := time.Parse(time.RFC3339, match.UtcDate)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			continue
		}

		if parsedDate.After(today) || parsedDate.Equal(today) {
			upcomingMatches = append(upcomingMatches, match)
		}
	}

	sort.Slice(upcomingMatches, func(i, j int) bool {
		dateI, _ := time.Parse(time.RFC3339, upcomingMatches[i].UtcDate)
		dateJ, _ := time.Parse(time.RFC3339, upcomingMatches[j].UtcDate)
		return dateI.Before(dateJ)
	})

	if len(upcomingMatches) == 0 {
		fmt.Println("No upcoming matches found for FC Barcelona.")
		return
	}

	limit := 3
	if len(upcomingMatches) < 3 {
		limit = len(upcomingMatches)
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(writer, "No\tHome Team\t\tAway Team\t\tDate\t\t\t\tStatus\tScore\n")
	fmt.Fprintf(writer, "--------------------------------------------------------------------------------------------\n")

	for i := 0; i < limit; i++ {
		match := upcomingMatches[i]
		parsedDate, err := time.Parse(time.RFC3339, match.UtcDate)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			continue
		}

		dateFormatted := parsedDate.Format("02/01/2006")

		if match.Status == "LIVE" {
			fmt.Fprintf(writer, "%d\t%-15s\t%-15s\t%-20s\t%-10s\t%-5d:%-5d\n", i+1, match.HomeTeam.Name, match.AwayTeam.Name, dateFormatted, "LIVE", match.Score.HomeTeam, match.Score.AwayTeam)
		} else {
			fmt.Fprintf(writer, "%d\t%-15s\t%-15s\t%-20s\t%-10s\t%-5s\n", i+1, match.HomeTeam.Name, match.AwayTeam.Name, dateFormatted, "Upcoming", "-- : --")
		}
	}
	writer.Flush()
}

func main() {
	getFCBarcelonaMatches()
}


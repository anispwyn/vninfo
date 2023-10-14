package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type VNDBEntry struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Image         Image    `json:"image"`
	AltTitle      string   `json:"alttitle"`
	Official      bool     `json:"titles.official"`
	OriginalLang  string   `json:"olang"`
	DevStatus     int      `json:"devstatus"`
	Released      string   `json:"released"`
	Languages     []string `json:"languages"`
	Platforms     []string `json:"platforms"`
	Description   string   `json:"description"`
	Rating        float64  `json:"rating"`
	VoteCount     int      `json:"votecount"`
	Length        int      `json:"length"`
	LengthMinutes int      `json:"length_minutes"`
	LengthVotes   int      `json:"length_votes"`
}

type Image struct {
	URL string `json:"url"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the title you want to search for: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	apiURL := "https://api.vndb.org/kana/vn"
	query := map[string]interface{}{
		"filters": []interface{}{"search", "=", title},
		"fields": strings.Join([]string{
			"id",
			"title",
			"image.url",
			"alttitle",
			"titles.official",
			"olang",
			"devstatus",
			"released",
			"languages",
			"platforms",
			"description",
			"rating",
			"votecount",
			"length",
			"length_minutes",
			"length_votes",
		}, ","),
	}

	resp, err := fetchData(apiURL, query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(resp) == 1 {
		entry := resp[0]
		printEntry(entry)
	} else if len(resp) > 1 {
		fmt.Println("More than one query was found. Printing names and IDs:")
		for i, entry := range resp {
			fmt.Printf("%d. %s (ID: %s)\n", i+1, entry.Title, entry.ID)
		}
		fmt.Print("Choose a query by entering the corresponding number: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		index := parseChoice(choice, len(resp))
		if index == -1 {
			fmt.Println("Invalid choice.")
			return
		}
		entry := resp[index]
		printEntry(entry)
	} else {
		fmt.Println("No matching query was found.")
	}
}

func fetchData(apiURL string, query map[string]interface{}) ([]VNDBEntry, error) {
	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(apiURL, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse struct {
		Results []VNDBEntry `json:"results"`
	}

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse.Results, nil
}

func printEntry(entry VNDBEntry) {
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println(color.CyanString("Title: " + entry.Title))
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("ID: %s\n", color.YellowString(entry.ID))
	fmt.Printf("Image URL: %s\n", color.BlueString(entry.Image.URL))
	fmt.Printf("Alternate Title: %s\n", entry.AltTitle)
	fmt.Printf("Official Title: %s\n", color.GreenString(fmt.Sprintf("%v", entry.Official)))
	fmt.Printf("Original Language: %s\n", color.MagentaString(entry.OriginalLang))
	fmt.Printf("Development Status: %s\n", formatDevStatus(entry.DevStatus))
	fmt.Printf("Release Date: %s\n", entry.Released)
	fmt.Printf("Languages: %s\n", formatList(entry.Languages))
	fmt.Printf("Platforms: %s\n", formatList(entry.Platforms))
	fmt.Printf("Description: %s\n", entry.Description)
	fmt.Printf("Rating: %s\n", color.YellowString(fmt.Sprintf("%.2f", entry.Rating)))
	fmt.Printf("Vote Count: %d\n", entry.VoteCount)
	fmt.Printf("Length: %d\n", entry.Length)
	fmt.Printf("Length in Minutes: %d\n", entry.LengthMinutes)
	fmt.Printf("Length Votes: %d\n", entry.LengthVotes)
	fmt.Println(strings.Repeat("=", 40))
}

func formatDevStatus(status int) string {
	switch status {
	case 0:
		return color.GreenString("Finished")
	case 1:
		return color.YellowString("In development")
	case 2:
		return color.RedString("Cancelled")
	default:
		return color.RedString("Unknown")
	}
}

func formatList(items []string) string {
	return color.CyanString(strings.Join(items, ", "))
}

func parseChoice(choice string, maxIndex int) int {
	index := -1
	if num, err := strconv.Atoi(choice); err == nil {
		if num >= 1 && num <= maxIndex {
			index = num - 1
		}
	}
	return index
}

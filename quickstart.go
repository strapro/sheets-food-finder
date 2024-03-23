package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type daySelections struct {
	start      int
	end        int
	selections []string
}

var weekDaysMap = map[string]int{
	"ΔΕΥΤΕΡΑ":   1,
	"ΤΡΙΤΗ":     2,
	"ΤΕΤΑΡΤΗ":   3,
	"ΠΕΜΠΤΗ":    4,
	"ΠΑΡΑΣΚΕΥΗ": 5,
	"ΣΑΒΒΑΤΟ":   6,
	"ΚΥΡΙΑΚΗ":   0,
}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	ctx := context.Background()

	spreadsheetId := os.Getenv("SPREADSHEET_ID")

	srv, err := sheets.NewService(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// A map containing the selections for each day
	daysSelections := getDaysSelections(srv, spreadsheetId)

	// Get the row of the user
	row := getUserRow(srv, spreadsheetId, os.Getenv("USER_NAME")) + 1

	// Get the current day
	weekday := int(time.Now().Weekday())

	// Get the selections for the current day for the user
	selections := getUserSelectionsForDay(srv, spreadsheetId, row, *daysSelections[weekday])

	// Print the selections
	for _, selection := range selections {
		fmt.Println(selection)
	}
}

func getDaysSelections(srv *sheets.Service, spreadsheetId string) map[int]*daySelections {
	var daysSelections = make(map[int]*daySelections)

	// Read the first row of the sheet to try and figure out the selections for each day
	// Also read the second row to get the end index for the last day and the selection names
	readRange := "1:2"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		var previousValue int = -1

		for i, cell := range resp.Values[0] {
			// If the cell is not empty, then we have a new day
			if cell != "" {
				// If we have a previous value, then we can set the end index for the previous day
				if previousValue != -1 {
					daysSelections[previousValue].end = i
				}
				var dayIndex = getWeekDayIndex(cell.(string))

				daysSelections[dayIndex] = &daySelections{start: i + 1, selections: []string{resp.Values[1][i].(string)}}

				previousValue = dayIndex
			} else {
				if previousValue != -1 {
					daysSelections[previousValue].selections = append(daysSelections[previousValue].selections, resp.Values[1][i].(string))
				}
			}
		}

		// Set the end index for the last day
		daysSelections[previousValue].end = len(resp.Values[1]) - 1
	}

	return daysSelections
}

func getWeekDayIndex(weekDay string) int {
	for key := range weekDaysMap {
		if strings.Contains(weekDay, key) {
			return weekDaysMap[key]
		}
	}

	return -1
}

func getUserRow(srv *sheets.Service, spreadsheetId string, name string) int {
	readRange := "A:A"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for i, row := range resp.Values {
			if len(row) > 0 && row[0] == name {
				return i
			}
		}
	}

	return -1
}

func getUserSelectionsForDay(srv *sheets.Service, spreadsheetId string, row int, daySelections daySelections) []string {
	selections := make([]string, 0)

	// Read only the columns that contain the selections for the day only for the user
	readRange := fmt.Sprintf("R%dC%d:R%dC%d", row, daySelections.start, row, daySelections.end)

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for i, cell := range resp.Values[0] {
			if cell != "" {
				selections = append(selections, fmt.Sprintf("%s %s", cell.(string), daySelections.selections[i]))
			}
		}
	}

	return selections
}

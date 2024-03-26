package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sheetsFoodFinder/pkg/authhelper"
	"sheetsFoodFinder/pkg/models"
	"sheetsFoodFinder/pkg/sheetshelper"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var SpreadsheetId string

func init() {
	godotenv.Load(".env")
}

func main() {
	ctx := context.Background()

	// start := time.Now()

	client := authhelper.GetClient()

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	if SpreadsheetId == "" {
		SpreadsheetId = os.Getenv("SPREADSHEET_ID")
	}

	username := ""
	if len(os.Args) > 1 {
		username = os.Args[1]
	} else {
		username = os.Getenv("USER_NAME")
	}

	sheetName := models.GetWeekRange()

	// Create channels to receive the results
	daysSelectionsChan := make(chan map[int]*models.DaySelections)
	rowChan := make(chan int)

	// Start goroutine to get the selections for each day
	go func() {
		daysSelectionsChan <- sheetshelper.GetDaysSelections(srv, SpreadsheetId, sheetName)
	}()

	// Start goroutine to get the row of the user
	go func() {
		rowChan <- sheetshelper.GetUserRow(srv, SpreadsheetId, sheetName, username) + 1
	}()

	// Receive the results from the channels
	daysSelections := <-daysSelectionsChan
	row := <-rowChan

	// Get the current day
	weekday := int(time.Now().Weekday())

	// Get the selections for the current day for the user
	selections := sheetshelper.GetUserSelectionsForDay(srv, SpreadsheetId, sheetName, row, *daysSelections[weekday])

	// Print the selections
	for _, selection := range selections {
		fmt.Println(selection)
	}

	// elapsed := time.Since(start)
	// fmt.Printf("Time took %s\n", elapsed)
}

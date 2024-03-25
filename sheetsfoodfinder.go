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

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	ctx := context.Background()

	// start := time.Now()

	client := authhelper.GetClient()

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	username := os.Getenv("USER_NAME")

	// Create channels to receive the results
	daysSelectionsChan := make(chan map[int]*models.DaySelections)
	rowChan := make(chan int)

	// Start goroutine to get the selections for each day
	go func() {
		daysSelectionsChan <- sheetshelper.GetDaysSelections(srv, spreadsheetId)
	}()

	// Start goroutine to get the row of the user
	go func() {
		rowChan <- sheetshelper.GetUserRow(srv, spreadsheetId, username) + 1
	}()

	// Receive the results from the channels
	daysSelections := <-daysSelectionsChan
	row := <-rowChan

	// Get the current day
	weekday := int(time.Now().Weekday())

	// Get the selections for the current day for the user
	selections := sheetshelper.GetUserSelectionsForDay(srv, spreadsheetId, row, *daysSelections[weekday])

	// Print the selections
	for _, selection := range selections {
		fmt.Println(selection)
	}

	// elapsed := time.Since(start)
	// fmt.Printf("Time took %s\n", elapsed)
}

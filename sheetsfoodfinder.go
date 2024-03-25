package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

	spreadsheetId := os.Getenv("SPREADSHEET_ID")

	srv, err := sheets.NewService(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// A map containing the selections for each day
	daysSelections := sheetshelper.GetDaysSelections(srv, spreadsheetId)

	// Get the row of the user
	row := sheetshelper.GetUserRow(srv, spreadsheetId, os.Getenv("USER_NAME")) + 1

	// Get the current day
	weekday := int(time.Now().Weekday())

	// Get the selections for the current day for the user
	selections := sheetshelper.GetUserSelectionsForDay(srv, spreadsheetId, row, *daysSelections[weekday])

	// Print the selections
	for _, selection := range selections {
		fmt.Println(selection)
	}
}

package sheetshelper

import (
	"fmt"
	"log"
	"sheetsFoodFinder/pkg/models"
	"strings"

	"google.golang.org/api/sheets/v4"
)

func GetDaysSelections(srv *sheets.Service, spreadsheetId string, sheetName string) map[int]*models.DaySelections {
	var daysSelections = make(map[int]*models.DaySelections)

	// Read the first row of the sheet to try and figure out the selections for each day
	// Also read the second row to get the end index for the last day and the selection names
	readRange := fmt.Sprintf("%s!1:2", sheetName)

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found when trying to get the selections for each day.")
	} else {
		var previousValue int = -1

		for i, cell := range resp.Values[0] {
			if strings.TrimSpace(cell.(string)) == "-" {
				continue
			}

			// If the cell is not empty, then we have a new day
			if cell != "" {
				// If we have a previous value, then we can set the end index for the previous day
				if previousValue != -1 {
					daysSelections[previousValue].End = i
				}
				var dayIndex = models.GetWeekDayIndex(cell.(string))

				daysSelections[dayIndex] = &models.DaySelections{Start: i + 1, Selections: []string{resp.Values[1][i].(string)}}

				previousValue = dayIndex
			} else {
				if previousValue != -1 {
					daysSelections[previousValue].Selections = append(daysSelections[previousValue].Selections, resp.Values[1][i].(string))
				}
			}
		}

		// Set the end index for the last day
		daysSelections[previousValue].End = len(resp.Values[1]) - 1
	}

	return daysSelections
}

func GetUserRow(srv *sheets.Service, spreadsheetId string, sheetName string, name string) int {
	readRange := fmt.Sprintf("%s!A:A", sheetName)

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found when trying to get the first column containing the users.")
	} else {
		for i, row := range resp.Values {
			if len(row) > 0 && row[0] == name {
				return i
			}
		}
	}

	return -1
}

func GetUserSelectionsForDay(srv *sheets.Service, spreadsheetId string, sheetName string, row int, daySelections models.DaySelections) []string {
	selections := make([]string, 0)

	// Read only the columns that contain the selections for the day only for the user
	readRange := fmt.Sprintf("%s!R%dC%d:R%dC%d", sheetName, row, daySelections.Start, row, daySelections.End)

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found when trying to get the selections of the user for the current day.")
	} else {
		for i, cell := range resp.Values[0] {
			if cell != "" {
				selections = append(selections, fmt.Sprintf("(%s) %s", cell.(string), daySelections.Selections[i]))
			}
		}
	}

	return selections
}

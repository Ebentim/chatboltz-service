package google

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// SheetsService provides methods for interacting with the Google Sheets API.
type SheetsService struct {
	service *sheets.Service
}

// NewSheetsService creates a new SheetsService.
func NewSheetsService(ctx context.Context, ts oauth2.TokenSource) (*SheetsService, error) {
	client := oauth2.NewClient(ctx, ts)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}
	return &SheetsService{service: srv}, nil
}

// CreateSpreadsheet creates a new Google Sheet.
func (s *SheetsService) CreateSpreadsheet(title string) (*sheets.Spreadsheet, error) {
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: title,
		},
	}
	return s.service.Spreadsheets.Create(spreadsheet).Do()
}

// GetValues retrieves values from a spreadsheet.
func (s *SheetsService) GetValues(spreadsheetID, readRange string) (*sheets.ValueRange, error) {
	return s.service.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
}

// UpdateValues updates values in a spreadsheet.
func (s *SheetsService) UpdateValues(spreadsheetID, rangeToUpdate string, values [][]interface{}) (*sheets.UpdateValuesResponse, error) {
	valueRange := &sheets.ValueRange{
		Values: values,
	}
	return s.service.Spreadsheets.Values.Update(spreadsheetID, rangeToUpdate, valueRange).ValueInputOption("USER_ENTERED").Do()
}

// AppendValues appends values to a spreadsheet.
func (s *SheetsService) AppendValues(spreadsheetID, rangeToAppend string, values [][]interface{}) (*sheets.AppendValuesResponse, error) {
	valueRange := &sheets.ValueRange{
		Values: values,
	}
	return s.service.Spreadsheets.Values.Append(spreadsheetID, rangeToAppend, valueRange).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Do()
}

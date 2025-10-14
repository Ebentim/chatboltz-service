package usecase

import (
	gsheets "google.golang.org/api/sheets/v4"
)

// SheetsIntegration defines the interface for sheets operations.
type SheetsIntegration interface {
	CreateSpreadsheet(title string) (*gsheets.Spreadsheet, error)
	GetValues(spreadsheetID, readRange string) (*gsheets.ValueRange, error)
	UpdateValues(spreadsheetID, rangeToUpdate string, values [][]interface{}) (*gsheets.UpdateValuesResponse, error)
	AppendValues(spreadsheetID, rangeToAppend string, values [][]interface{}) (*gsheets.AppendValuesResponse, error)
}

// SheetsUseCase handles the business logic for sheets operations.
type SheetsUseCase struct {
	sheetsIntegration SheetsIntegration
}

// NewSheetsUseCase creates a new SheetsUseCase.
func NewSheetsUseCase(si SheetsIntegration) *SheetsUseCase {
	return &SheetsUseCase{sheetsIntegration: si}
}

// CreateSpreadsheet creates a new Google Sheet.
func (uc *SheetsUseCase) CreateSpreadsheet(title string) (*gsheets.Spreadsheet, error) {
	return uc.sheetsIntegration.CreateSpreadsheet(title)
}

// GetValues retrieves values from a spreadsheet.
func (uc *SheetsUseCase) GetValues(spreadsheetID, readRange string) (*gsheets.ValueRange, error) {
	return uc.sheetsIntegration.GetValues(spreadsheetID, readRange)
}

// UpdateValues updates values in a spreadsheet.
func (uc *SheetsUseCase) UpdateValues(spreadsheetID, rangeToUpdate string, values [][]interface{}) (*gsheets.UpdateValuesResponse, error) {
	return uc.sheetsIntegration.UpdateValues(spreadsheetID, rangeToUpdate, values)
}

// AppendValues appends values to a spreadsheet.
func (uc *SheetsUseCase) AppendValues(spreadsheetID, rangeToAppend string, values [][]interface{}) (*gsheets.AppendValuesResponse, error) {
	return uc.sheetsIntegration.AppendValues(spreadsheetID, rangeToAppend, values)
}

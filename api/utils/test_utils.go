package utils

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/waynekn/tablesync/api/models"
)

// InsertTestData inserts a sample spreadsheet record into the provided test database.
//
// This function is used to populate the database with predictable data for tests
// that involve retrieving spreadsheets. It panics if the insert operation fails,
// as the absence of test data would cause dependent tests to fail.
func InsertTestData(testDb *sql.DB) {
	cols := [][]string{
		{"header1", "header2"},
		{"abc", "def"},
	}
	colsJson, _ := json.Marshal(cols)

	sheet := models.Spreadsheet{
		ID:          GenerateID(),
		Owner:       "test-user",
		Title:       "test sheet",
		Deadline:    time.Time{},
		Description: "Test sheet",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Data:        colsJson,
	}

	_, err := testDb.Exec(`INSERT INTO spreadsheets
	 (id, owner, title, description, deadline, data)
	  VALUES ($1, $2, $3, $4, $5, $6)`,
		sheet.ID, sheet.Owner, sheet.Title, sheet.Description, sheet.Deadline, sheet.Data)

	if err != nil {
		panic("Failed to create spreadsheet: " + err.Error())
	}
}

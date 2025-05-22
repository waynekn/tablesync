package repo

import (
	"database/sql"
	"log/slog"

	"github.com/waynekn/tablesync/api/models"
)

type SpreadsheetRepo interface {
	InsertSpreadsheet(sheet models.SpreadsheetInit, columns []byte, owner, id string) error
}

type spreadsheetRepo struct {
	db *sql.DB
}

// NewSpreadsheetRepo creates a new instance of SpreadsheetRepo
// with the provided database connection.
func NewSpreadsheetRepo(db *sql.DB) SpreadsheetRepo {
	return &spreadsheetRepo{db: db}
}

// InsertSpreadsheet inserts a new spreadsheet into the database.
func (s *spreadsheetRepo) InsertSpreadsheet(sheet models.SpreadsheetInit, columns []byte, owner, id string) error {
	_, err := s.db.Exec(`INSERT INTO spreadsheets
	 (id, owner, title, description, deadline, data)
	  VALUES ($1, $2, $3, $4, $5, $6)`,
		id, owner, sheet.Title, sheet.Description, sheet.Deadline, columns)

	if err != nil {
		slog.Error("Failed to create spreadsheet", "error", err)
		return err
	}
	return nil
}

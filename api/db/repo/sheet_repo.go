package repo

import (
	"database/sql"
	"log/slog"

	"github.com/waynekn/tablesync/api/models"
)

type SpreadsheetRepo interface {
	InsertSpreadsheet(sheet models.SpreadsheetInit, columns []byte, owner, id string) error
	GetByOwner(owner string) (*[]models.Spreadsheet, error)
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

// GetByOwner retrieves spreadsheets created by the `owner` from the db
func (s *spreadsheetRepo) GetByOwner(owner string) (*[]models.Spreadsheet, error) {
	rows, err := s.db.Query(`SELECT *
		FROM spreadsheets WHERE owner = $1`,
		owner)
	if err != nil {
		if err == sql.ErrNoRows {
			emptySlice := []models.Spreadsheet{}
			return &emptySlice, nil
		}
		slog.Error("Failed to query spreadsheets", "error", err)
		return nil, err
	}
	defer rows.Close()

	spreadsheets := make([]models.Spreadsheet, 0, 20)
	for rows.Next() {
		var sheet models.Spreadsheet
		if err := rows.Scan(&sheet.ID, &sheet.Title, &sheet.Description, &sheet.Owner,
			&sheet.CreatedAt, &sheet.UpdatedAt, &sheet.Data, &sheet.Deadline); err != nil {
			slog.Error("Failed to scan spreadsheet row", "error", err)
			return nil, err
		}
		spreadsheets = append(spreadsheets, sheet)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error occurred during row iteration", "error", err)
		return nil, err
	}

	return &spreadsheets, nil
}

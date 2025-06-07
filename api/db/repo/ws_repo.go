package repo

import (
	"database/sql"
	"log/slog"

	"github.com/waynekn/tablesync/api/models"
)

type WsRepo interface {
	GetSheetByID(sheetID string) (*models.Spreadsheet, error)
}

type wsRepo struct {
	db *sql.DB
}

// NewWsRepo creates a new instance of wsRepo with the provided database connection.
func NewWsRepo(db *sql.DB) WsRepo {
	return &wsRepo{
		db: db,
	}
}

// GetSheetByID retrieves a spreadsheet by its ID from the database.
func (ws *wsRepo) GetSheetByID(sheetID string) (*models.Spreadsheet, error) {
	var sheet models.Spreadsheet

	err := ws.db.QueryRow(`SELECT id, title, description, owner, created_at, updated_at, data, deadline 
                           FROM spreadsheets WHERE id = $1`, sheetID).
		Scan(&sheet.ID, &sheet.Title, &sheet.Description, &sheet.Owner,
			&sheet.CreatedAt, &sheet.UpdatedAt, &sheet.Data, &sheet.Deadline)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		slog.Error("Failed to query spreadsheet", "err", err)
		return nil, err
	}

	return &sheet, nil
}

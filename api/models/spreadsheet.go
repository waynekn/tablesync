package models

import "time"

// SpreadsheetInit represents the payload required to create a new spreadsheet.
type SpreadsheetInit struct {
	Title       string    `json:"title" binding:"required,max=255"`
	Description string    `json:"description" binding:"required"`
	Deadline    time.Time `json:"deadline" time_format:"2006-01-02T15:04:05Z07:00" binding:"required"` // Deadline in RFC3339 format
	ColTitles   []string  `json:"colTitles" binding:"required,min=1"`
}

// Spreadsheet represents a spreadsheet stored in the database.
type Spreadsheet struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Owner       string    `json:"owner"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Data        []byte    `json:"data"` // [][]string stored as jsonb
	Deadline    time.Time `json:"deadline"`
}

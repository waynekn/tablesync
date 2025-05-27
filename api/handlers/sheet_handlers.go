package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/waynekn/tablesync/api/db/repo"
	"github.com/waynekn/tablesync/api/models"
	"github.com/waynekn/tablesync/api/utils"
)

type SpreadsheetHandler struct {
	repo repo.SpreadsheetRepo
}

// NewSpreadsheetHandler creates a new instance of SpreadsheetHandler
// with the provided repository.
func NewSpreadsheetHandler(repo repo.SpreadsheetRepo) *SpreadsheetHandler {
	return &SpreadsheetHandler{repo: repo}
}

// CreateSpreadsheetHandler handles the creation of a new spreadsheet.
func (h *SpreadsheetHandler) CreateSpreadsheetHandler(c *gin.Context) {
	var sheet models.SpreadsheetInit

	token, err := utils.TokenFromContext(c)
	if err != nil {
		slog.Error("Unauthorized request gained access to a protected endpoint", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&sheet); err != nil {
		// Handle validation errors
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			detail := make(map[string]string)
			for _, fieldErr := range verr {
				detail[fieldErr.Field()] = utils.GetValidationErrorMessage(fieldErr)
			}
			c.JSON(http.StatusBadRequest, detail)
			return
		}

		// Handle time parsing errors
		var jsonErr *json.UnmarshalTypeError
		var timeErr *time.ParseError
		if (errors.As(err, &jsonErr) && jsonErr.Field == "deadline") ||
			errors.As(err, &timeErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"deadline": "Invalid deadline time format.",
			})
			return
		}

		// Handle other JSON parsing errors
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON syntax",
			})
			return
		}
		// If it's not a validation error, return a generic message
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id := utils.GenerateID()
	// Convert to a 2D array representation
	columns := make([][]string, 1)
	columns[0] = append(columns[0], sheet.ColTitles...)
	columnsJson, err := json.Marshal(columns)
	if err != nil {
		slog.Error("Failed to marshal columns", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create spreadsheet"})
		return
	}

	err = h.repo.InsertSpreadsheet(sheet, columnsJson, token.Subject(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create spreadsheet"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Spreadsheet created successfully"})
}

// GetOwnSpreadsheetsHandler handles requests to retrieve spreadsheets owned by
// the authenticated user.
func (h *SpreadsheetHandler) GetOwnSpreadsheetsHandler(c *gin.Context) {
	token, err := utils.TokenFromContext(c)
	if err != nil {
		slog.Error("Unauthorized request gained access to a protected endpoint", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	spreadsheets, err := h.repo.GetByOwner(token.Subject())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred while retrieving spreadsheets. Please try again later."})
		return
	}

	c.JSON(http.StatusOK, spreadsheets)
}

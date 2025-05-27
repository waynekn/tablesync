package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/waynekn/tablesync/api/models"
)

// CreateTestCtxWithToken creates a test Gin context with a JWT token.
//
// The function initializes a new Gin test context using the provided http.ResponseWriter.
// It creates a JWT token with the subject set to "test-user" and stores the token in the
// context under the "token" key. The returned context is suitable for testing handlers
// that require a JWT token.
func CreateTestCtxWithToken(rec http.ResponseWriter) *gin.Context {
	ctx, _ := gin.CreateTestContext(rec)
	token := jwt.New()
	// set the sub field only as its the only one used by the handlers
	// so far
	token.Set("sub", "test-user")
	ctx.Set("token", token)
	return ctx
}

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

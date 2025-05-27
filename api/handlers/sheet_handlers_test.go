package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/waynekn/tablesync/api"
	"github.com/waynekn/tablesync/api/db"
	"github.com/waynekn/tablesync/api/db/repo"
	"github.com/waynekn/tablesync/api/models"
)

var testDb *sql.DB

// TestMain establishes a connection to the test database
func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env.test")
	conn, err := db.Connect()
	if err != nil {
		panic("Failed to connect to the test database")
	}
	defer conn.Close()
	testDb = conn
	gin.SetMode(gin.TestMode)
	api.RegisterJSONTagNameFormatter()
	// Run the tests
	m.Run()
}

// setupTest initializes a Gin context and a response recorder for testing
// It sets up a mock JWT token and prepares the request body with the provided data.
func setupCreateSpreadsheetTestContext(data models.SpreadsheetInit) (*gin.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	token := jwt.New()
	// set the sub field only as its the only one used by the handler
	// in the test
	token.Set("sub", "test-user")
	ctx.Set("token", token)

	var bodyReader *bytes.Reader

	jsonBytes, _ := json.Marshal(data)
	bodyReader = bytes.NewReader(jsonBytes)

	ctx.Request = httptest.NewRequest("POST", "/spreadsheet/create/", bodyReader)
	return ctx, rec
}

// TestCreateSpreadsheet tests the CreateSpreadsheet function
// It checks for various scenarios including successful creation, missing request body,
// and invalid data formats.
//
// It asserts that the fields with errors are returned in the response and the appropriate
// HTTP status codes are returned for each case.
func TestCreateSpreadsheet(t *testing.T) {
	testRepo := repo.NewSpreadsheetRepo(testDb)
	h := NewSpreadsheetHandler(testRepo)

	data := models.SpreadsheetInit{
		Title:       "Test Spreadsheet",
		Description: "Test Description",
		Deadline:    time.Now().Add(24 * time.Hour),
		ColTitles:   []string{"Column1", "Column2"},
	}

	t.Run("Should create spreadsheet successfully", func(t *testing.T) {
		ctx, rec := setupCreateSpreadsheetTestContext(data)
		h.CreateSpreadsheetHandler(ctx)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	tests := []struct {
		name    string
		modify  func(m *models.SpreadsheetInit)
		wantKey string
	}{
		{"missing title", func(m *models.SpreadsheetInit) { m.Title = "" }, "title"},
		{"long title", func(m *models.SpreadsheetInit) { m.Title = strings.Repeat("a", 256) }, "title"},
		{"missing description", func(m *models.SpreadsheetInit) { m.Description = "" }, "description"},
		{"invalid deadline", func(m *models.SpreadsheetInit) { m.Deadline = time.Time{} }, "deadline"},
		{"no col titles", func(m *models.SpreadsheetInit) { m.ColTitles = []string{} }, "colTitles"},
	}

	for _, tt := range tests {
		t.Run("Should not create spreadsheet with "+tt.name, func(t *testing.T) {
			input := data
			tt.modify(&input)
			ctx, rec := setupCreateSpreadsheetTestContext(input)
			h.CreateSpreadsheetHandler(ctx)
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp map[string]string
			_ = json.Unmarshal(rec.Body.Bytes(), &resp)
			_, ok := resp[tt.wantKey]
			assert.True(t, ok)
		})
	}

}


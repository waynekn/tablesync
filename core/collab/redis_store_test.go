package collab

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/waynekn/tablesync/api/utils"
	"github.com/waynekn/tablesync/core/rdb"
)

var testStore *Store

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env.test")

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	redisClient, err := rdb.Connect(redisAddr, redisPassword, redisDB)

	if err != nil {
		panic(err.Error())
	}

	testStore = NewStore(redisClient)

	m.Run()
}

func TestSheetExists_WhenSheetDoesNotExist(t *testing.T) {
	randomSheetID := utils.GenerateID()

	exists, err := testStore.SheetExists(randomSheetID)

	assert.NoError(t, err, "should not return an error when checking for a non-existent sheet")
	assert.False(t, exists, "should return false for a non-existent sheet")
}

func TestSheetExists_WhenSheetExists(t *testing.T) {
	sheetID := utils.GenerateID()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create a sheet in Redis for testing
	err := testStore.rdb.HSet(ctx, sheetID, "1:1", "value").Err()
	assert.NoError(t, err, "should not return an error when creating a test sheet")

	exists, err := testStore.SheetExists(sheetID)

	assert.NoError(t, err, "should not return an error when checking for an existing sheet")
	assert.True(t, exists, "should return true for an existing sheet")
}

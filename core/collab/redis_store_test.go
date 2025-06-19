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

func TestInitRedisSheet(t *testing.T) {
	sheetID := utils.GenerateID()
	sheetDeadline := time.Now().Add(10 * time.Minute)
	sheetData := &[][]string{
		{"A1", "B1"},
		{"A2", "B2"},
	}

	err := testStore.InitRedisSheet(sheetID, sheetDeadline, sheetData)
	assert.NoError(t, err, "should not return an error when initializing a sheet with a valid deadline")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Check that the key exists
	exists, err := testStore.rdb.Exists(ctx, sheetID).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), exists, "sheet key should exist in Redis")

	// Check TTL is set
	ttl, err := testStore.rdb.TTL(ctx, sheetID).Result()
	assert.NoError(t, err, "should not return an error when checking TTL")
	expectedTTL := 15 * time.Minute
	assert.Equal(t, ttl, expectedTTL)

	// Check that the correct data is stored
	result, err := testStore.rdb.HGetAll(ctx, sheetID).Result()
	assert.NoError(t, err)
	assert.Equal(t, 4, len(result), "should store 4 cells for a 2x2 sheet")

	expected := map[string]string{
		"0:0": "A1",
		"0:1": "B1",
		"1:0": "A2",
		"1:1": "B2",
	}

	for k, v := range expected {
		assert.Equal(t, v, result[k], "value mismatch for cell %s", k)
	}
}

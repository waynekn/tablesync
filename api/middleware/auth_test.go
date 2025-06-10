package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/waynekn/tablesync/core/rdb"
)

var testRdb *redis.Client

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env.test")

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	redisClient, err := rdb.Connect(redisAddr, redisPassword, redisDB)

	if err != nil {
		panic(err.Error())
	}

	testRdb = redisClient

	m.Run()
}

func TestAuthMiddleWare(t *testing.T) {
	// Set up a test JWKS server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
            "keys": [{
                "kty": "RSA",
                "kid": "test-key",
                "n": "test-modulus",
                "e": "AQAB"
            }]
        }`))
	}))
	defer testServer.Close()

	// Set the `PUB_KEY_URL` to our test server
	t.Setenv("PUB_KEY_URL", testServer.URL)

	// Set up test router.
	router := gin.Default()
	router.GET("/protected", RequireAuth(testRdb), func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{})
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Bad Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "BadFormatToken")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer token")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

}

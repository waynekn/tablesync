package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

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
	router.GET("/protected", RequireAuth(), func(ctx *gin.Context) {
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

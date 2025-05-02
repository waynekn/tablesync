package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// RequireAuth is a Gin middleware that validates an access token.
// If the token is valid, the request proceeds to the next handler.
// If the token is missing or invalid, the request is aborted.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		pubKeyUrl := os.Getenv("PUB_KEY_URL")

		authHeader := c.GetHeader("Authorization")
		splitAuthHeader := strings.Split(authHeader, " ")

		if len(splitAuthHeader) != 2 || strings.ToLower(splitAuthHeader[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid Authorization header",
			})

			return
		}

		keySet, err := jwk.Fetch(c.Request.Context(), pubKeyUrl)
		if err != nil {
			slog.Error("Failed to fetch JWK", "error", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}

		token, err := jwt.Parse(
			[]byte(splitAuthHeader[1]),
			jwt.WithKeySet(keySet),
			jwt.WithValidate(true),
		)

		if err != nil {
			slog.Info("Invalid token", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		c.Set("token", token)

		c.Next()
	}
}

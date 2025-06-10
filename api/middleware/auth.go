package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/redis/go-redis/v9"
)

// RequireAuth is a Gin middleware that validates an access token.
// If the token is valid, the request proceeds to the next handler.
// If the token is missing or invalid, the request is aborted.
func RequireAuth(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		pubKeyUrl := os.Getenv("PUB_KEY_URL")

		authHeader := c.GetHeader("Authorization")
		splitAuthHeader := strings.Split(authHeader, " ")

		if len(splitAuthHeader) != 2 || strings.ToLower(splitAuthHeader[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid Authorization header",
			})

			return
		}

		keySet, err := getKeySetFromRedis(ctx, rdb)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// If the key set is not cached, fetch it from the public key URL
		if keySet == nil {
			keySet, err = jwk.Fetch(c.Request.Context(), pubKeyUrl)
			if err != nil {
				slog.Error("Failed to fetch JWK", "error", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
			// Store the fetched key set in Redis for future use
			jsonKeySet, err := json.Marshal(keySet)
			if err != nil {
				slog.Error("Failed to marshal jwk key set", "error", err)
			} else {
				if err := rdb.Set(ctx, "jwk_keySet", jsonKeySet, 24*time.Hour).Err(); err != nil {
					slog.Error("Failed to store pub keys in Redis", "error", err)
				}
			}
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

// getKeySetFromRedis retrieves the JWK Set from Redis.
// It looks for the key "jwk_keySet" in Redis and attempts to parse it as a JWK Set.
// If the key set is not found, it returns nil without an error.
func getKeySetFromRedis(ctx context.Context, rdb *redis.Client) (jwk.Set, error) {
	val, err := rdb.Get(ctx, "jwk_keySet").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		slog.Error("Failed to get jwk key set from Redis", "error", err)
		return nil, err
	}

	keySet, err := jwk.ParseString(val)
	if err != nil {
		slog.Error("Failed to parse keys retrieved from Redis to JWK Set", "error", err)
		return nil, err
	}
	return keySet, nil
}

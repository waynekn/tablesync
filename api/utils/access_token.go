package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
)

// TokenFromContext retrieves the access token from the Gin context.
func TokenFromContext(c *gin.Context) (jwt.Token, error) {
	val, exists := c.Get("token")
	if !exists {
		return nil, errors.New("token missing from context")
	}
	token, ok := val.(jwt.Token)
	if !ok {
		return nil, errors.New("invalid token type in context")
	}
	return token, nil
}

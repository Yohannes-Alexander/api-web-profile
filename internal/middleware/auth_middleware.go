package middleware

import (
	"net/http"
	"strings"

	"apiprofile/internal/helper"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(secret string) gin.HandlerFunc {
	jm := helper.NewJWTManager(secret)
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}
		token := parts[1]
		claims, err := jm.ValidateAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("user_claims", claims)
		c.Next()
	}
}

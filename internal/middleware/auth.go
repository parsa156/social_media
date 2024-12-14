package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"social_media/pkg/jwt"
)

// AuthMiddleware validates the JWT token and sets user info in context.
func AuthMiddleware(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		claims, err := jwtManager.Verify(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the userID and username in context.
		c.Set("userID", claims.UserID)       // In our JWT, we can store the UUID as UserID
		c.Set("username", claims.Username)
		c.Next()
	}
}

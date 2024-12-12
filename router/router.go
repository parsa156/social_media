package router

import (
	"github.com/gin-gonic/gin"
	"social_media/internal/handler"
	"social_media/internal/middleware"
	"social_media/pkg/jwt"
)

// SetupRouter initializes the Gin engine with routes and middleware.
func SetupRouter(authHandler *handler.AuthHandler, jwtManager *jwt.JWTManager) *gin.Engine {
	r := gin.Default()

	// Public routes.
	api := r.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
	}

	// Protected routes example.
	protected := r.Group("/api/protected")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Example protected route.
		protected.GET("/profile", func(c *gin.Context) {
			// For demonstration, we just return the user info from context.
			userID, _ := c.Get("userID")
			username, _ := c.Get("username")
			c.JSON(200, gin.H{
				"message":  "This is a protected route.",
				"userID":   userID,
				"username": username,
			})
		})
	}

	return r
}

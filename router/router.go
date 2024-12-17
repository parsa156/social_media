package router

import (
	"github.com/gin-gonic/gin"
	"social_media/internal/handler"
	"social_media/internal/middleware"
	"social_media/pkg/jwt"
)

// SetupRouter initializes all routes with middleware.
func SetupRouter(authHandler *handler.AuthHandler, profileHandler *handler.ProfileHandler, jwtManager *jwt.JWTManager) *gin.Engine {
	r := gin.Default()

	// Public routes.
	public := r.Group("/api")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	// Protected routes (requires valid JWT).
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		protected.GET("/profile", profileHandler.GetProfile)
		protected.PUT("/profile", profileHandler.UpdateProfile)
		protected.DELETE("/profile", profileHandler.DeleteProfile)
	}

	return r
}

package router

import (
	"github.com/gin-gonic/gin"
	"social_media/internal/handler"
	"social_media/pkg/jwt"
)

// SetupRouter initializes all routes with middleware.
func SetupRouter(authHandler *handler.AuthHandler, jwtManager *jwt.JWTManager) *gin.Engine {
	r := gin.Default()

	// Public routes.
	public := r.Group("/api")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}


	return r
}

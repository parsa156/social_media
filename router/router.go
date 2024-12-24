package router

import (
	"github.com/gin-gonic/gin"
	"social_media/internal/handler"
	"social_media/internal/middleware"
	"social_media/pkg/jwt"
)

// SetupRouter initializes all routes with middleware.
func SetupRouter(authHandler *handler.AuthHandler, profileHandler *handler.ProfileHandler, convoHandler *handler.ConversationHandler, jwtManager *jwt.JWTManager) *gin.Engine {
	r := gin.Default()

	// Public routes.
	public := r.Group("/api")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	// Protected routes.
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Profile endpoints.
		protected.GET("/profile", profileHandler.GetProfile)
		protected.PUT("/profile", profileHandler.UpdateProfile)
		protected.DELETE("/profile", profileHandler.DeleteProfile)

		// Conversation endpoints.
		// Send a message (creates a conversation if needed).
		protected.POST("/conversations/send", convoHandler.SendMessageEndpoint)
		// List conversations for the authenticated user.
		protected.GET("/conversations", convoHandler.ListConversations)
		// Get messages for a specific conversation.
		protected.GET("/conversations/:id/messages", convoHandler.GetMessages)
		// Update a message (only sender can update).
		protected.PUT("/messages/:id", convoHandler.UpdateMessage)
		// Delete a message (only sender can delete).
		protected.DELETE("/messages/:id", convoHandler.DeleteMessage)
	}

	return r
}

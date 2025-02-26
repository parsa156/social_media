package router

import (
	"social_media/internal/handler"
	"social_media/internal/middleware"
	"social_media/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *handler.AuthHandler, profileHandler *handler.ProfileHandler, convoHandler *handler.ConversationHandler, roomHandler *handler.RoomHandler, jwtManager *jwt.JWTManager) *gin.Engine {
	r := gin.Default()

	// Public routes.
	public := r.Group("/api")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	// Protected routes.
	protected := r.Group("/api")
	{
		protected.Use(middleware.AuthMiddleware(jwtManager))
		// Profile endpoints.
		protected.GET("/profile", profileHandler.GetProfile)
		protected.PUT("/profile", profileHandler.UpdateProfile)
		protected.DELETE("/profile", profileHandler.DeleteProfile)

		// Conversation endpoints.
		protected.POST("/conversations/send", convoHandler.SendMessageEndpoint)
		protected.GET("/conversations", convoHandler.ListConversations)
		protected.GET("/conversations/:id/messages", convoHandler.GetMessages)
		protected.PUT("/messages/:id", convoHandler.UpdateMessage)
		protected.DELETE("/messages/:id", convoHandler.DeleteMessage)

		// Room endpoints.
		protected.POST("/rooms", roomHandler.CreateRoom)
		protected.PUT("/rooms", roomHandler.UpdateRoom)
		protected.DELETE("/rooms", roomHandler.DeleteRoom)
		protected.POST("/rooms/add-member", roomHandler.AddMember)
		protected.POST("/rooms/remove-member", roomHandler.RemoveMember)
		protected.POST("/rooms/promote-member", roomHandler.PromoteMember)
		protected.POST("/rooms/ban-member", roomHandler.BanMember)
		protected.POST("/rooms/unban-member", roomHandler.UnbanMember)
		protected.POST("/rooms/send-message", roomHandler.SendMessage)
		protected.DELETE("/rooms/delete-message", roomHandler.DeleteMessage)
		protected.GET("/rooms/:roomID/messages", roomHandler.GetMessages)
		protected.GET("/rooms/:roomID/members", roomHandler.GetMembers)
	}

	return r
}

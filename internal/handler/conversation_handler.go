package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"social_media/internal/service"
)
//ConversationHandler : struct
type ConversationHandler struct {
	convoService service.ConversationService
}

//NewConversationHandler : function
func NewConversationHandler(convoService service.ConversationService) *ConversationHandler {
	return &ConversationHandler{convoService: convoService}
}

// SendMessageEndpoint creates a conversation if needed and sends a message.
// Expects JSON: { "recipient_id": "uuid-of-recipient", "content": "message text" }
func (h *ConversationHandler) SendMessageEndpoint(c *gin.Context) {
	var req struct {
		RecipientID string `json:"recipient_id" binding:"required"`
		Content     string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Assume AuthMiddleware sets "userID" in context.
	senderID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	message, err := h.convoService.SendMessage(senderID.(string), req.RecipientID, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, message)
}

// ListConversations returns a list of conversations for the authenticated user.
func (h *ConversationHandler) ListConversations(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	convos, err := h.convoService.GetConversations(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, convos)
}

// GetMessages returns messages in a conversation.
// Endpoint: GET /api/conversations/:id/messages
func (h *ConversationHandler) GetMessages(c *gin.Context) {
	convoID := c.Param("id")
	messages, err := h.convoService.GetMessages(convoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// UpdateMessage allows the sender to update their message.
// Endpoint: PUT /api/messages/:id
func (h *ConversationHandler) UpdateMessage(c *gin.Context) {
	messageID := c.Param("id")
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	senderID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	message, err := h.convoService.UpdateMessage(senderID.(string), messageID, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, message)
}

// DeleteMessage allows the sender to delete their message.
// Endpoint: DELETE /api/messages/:id
func (h *ConversationHandler) DeleteMessage(c *gin.Context) {
	messageID := c.Param("id")
	senderID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	if err := h.convoService.DeleteMessage(senderID.(string), messageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "message deleted"})
}

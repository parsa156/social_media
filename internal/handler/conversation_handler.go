package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"social_media/internal/service"
)

type ConversationHandler struct {
	convoService service.ConversationService
}

// NewConversationHandler creates a new ConversationHandler.
func NewConversationHandler(convoService service.ConversationService) *ConversationHandler {
	return &ConversationHandler{convoService: convoService}
}

// SendMessageEndpoint accepts a recipient identifier (phone or username) and message content.
// Expected JSON:
// {
//    "recipient": "identifier", // either "1234567890" or "@johndoe"
//    "content": "Hello, how are you?"
// }
func (h *ConversationHandler) SendMessageEndpoint(c *gin.Context) {
	var req struct {
		Recipient string `json:"recipient" binding:"required"`
		Content   string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Retrieve senderID from the context (set by AuthMiddleware).
	senderID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	message, err := h.convoService.SendMessage(senderID.(string), req.Recipient, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, message)
}

// ListConversations returns all conversations for the authenticated user.
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

// GetMessages returns all messages for a specific conversation.
// The conversation ID is taken from the URL parameter.
func (h *ConversationHandler) GetMessages(c *gin.Context) {
	convoID := c.Param("id")
	messages, err := h.convoService.GetMessages(convoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// UpdateMessage allows the sender to update a message.
// The message ID is taken from the URL parameter.
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

// DeleteMessage allows the sender to delete a message.
// The message ID is taken from the URL parameter.
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

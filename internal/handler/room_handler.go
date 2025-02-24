package handler

import (
	"net/http"
	"social_media/internal/domain"
	"github.com/gin-gonic/gin"
	"social_media/internal/service"
)

type RoomHandler struct {
	roomService service.RoomService
}

func NewRoomHandler(roomService service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

type CreateRoomRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username"`
	Type     string `json:"type" binding:"required"` // "group" or "channel"
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	roomType := domain.RoomType(req.Type) // Convert string to domain.RoomType
	room, err := h.roomService.CreateRoom(userID.(string), req.Name, req.Username, roomType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, room)
}

type UpdateRoomRequest struct {
	RoomID      string `json:"room_id" binding:"required"`
	NewName     string `json:"new_name" binding:"required"`
	NewUsername string `json:"new_username"`
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	var req UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	room, err := h.roomService.UpdateRoom(req.RoomID, userID.(string), req.NewName, req.NewUsername)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, room)
}

type DeleteRoomRequest struct {
	RoomID string `json:"room_id" binding:"required"`
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	var req DeleteRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	if err := h.roomService.DeleteRoom(req.RoomID, userID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "room deleted"})
}

type AddMemberRequest struct {
	RoomID string `json:"room_id" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

func (h *RoomHandler) AddMember(c *gin.Context) {
	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requesterID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	if err := h.roomService.AddMember(req.RoomID, requesterID.(string), req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "member added"})
}

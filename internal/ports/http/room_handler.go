package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	roomapp "github.com/jorgerr9011/cartas-game-backend/internal/app/room"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/room"
)

type RoomHandler struct {
	usecase roomapp.UseCase
}

func NewRoomHandler(uc roomapp.UseCase) *RoomHandler {
	return &RoomHandler{usecase: uc}
}

func (h *RoomHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/rooms", h.createRoom)
	rg.POST("/rooms/:id/join", h.joinRoom)
	rg.POST("/rooms/:id/start", h.startGame)
	rg.POST("/rooms/:id/nextturn", h.nextTurn)
	rg.GET("/rooms/:id/currentplayer", h.currentPlayer)
}

type createRoomReq struct {
	Name string `json:"name" binding:"required"`
}

func (h *RoomHandler) createRoom(c *gin.Context) {
	var req createRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	room, err := h.usecase.CreateRoom(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, room)
}

func (h *RoomHandler) joinRoom(c *gin.Context) {
	id := room.RoomID(c.Param("id"))
	playerID := room.PlayerID(c.Query("playerID"))
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "playerID is required"})
		return
	}
	err := h.usecase.JoinRoom(id, playerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *RoomHandler) startGame(c *gin.Context) {
	id := room.RoomID(c.Param("id"))
	err := h.usecase.StartGame(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *RoomHandler) nextTurn(c *gin.Context) {
	id := room.RoomID(c.Param("id"))
	err := h.usecase.NextTurn(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *RoomHandler) currentPlayer(c *gin.Context) {
	id := room.RoomID(c.Param("id"))
	playerID, err := h.usecase.CurrentPlayer(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"currentPlayer": playerID})
}

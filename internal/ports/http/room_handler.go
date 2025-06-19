package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	roomapp "github.com/jorgerr9011/cartas-game-backend/internal/app/room"

	"github.com/jorgerr9011/cartas-game-backend/internal/domain/game"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/room"
)

type RoomHandler struct {
	roomusecase roomapp.UseCase
}

func NewRoomHandler(uc roomapp.UseCase) *RoomHandler {
	return &RoomHandler{
		roomusecase: uc,
	}
}

func (h *RoomHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/rooms", h.createRoom)
	rg.POST("/rooms/:id/join", h.joinRoom)
	rg.POST("/rooms/:id/start", h.startGame)
	rg.POST("/rooms/:id/nextturn", h.nextTurn)
	rg.GET("/rooms/:id/currentplayer", h.currentPlayer)
}

type createRoomReq struct {
	Id       room.RoomID `json:"id" binding:"required"`
	Name     string      `json:"name" binding:"required"`
	GameName string      `json:"game_name" binding:"required"`
}

func (h *RoomHandler) createRoom(c *gin.Context) {
	var req createRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gamefactory := game.NewGameFactory()
	newGame := gamefactory.NewGame(req.GameName)

	room, err := h.roomusecase.CreateRoom(req.Id, req.Name, newGame)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, room)
}

func (h *RoomHandler) joinRoom(c *gin.Context) {
	id := room.RoomID(c.Param("id"))
	playerID := player.PlayerID(c.Query("playerID"))
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "playerID is required"})
		return
	}
	err := h.roomusecase.JoinRoom(id, playerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *RoomHandler) startGame(c *gin.Context) {
	id := room.RoomID(c.Param("id"))
	_, err := h.roomusecase.StartGame(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *RoomHandler) nextTurn(c *gin.Context) {
	id := room.RoomID(c.Param("id"))
	err := h.roomusecase.NextTurn(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *RoomHandler) currentPlayer(c *gin.Context) {
	id := room.RoomID(c.Param("id"))
	playerID, err := h.roomusecase.CurrentPlayer(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"currentPlayer": playerID})
}

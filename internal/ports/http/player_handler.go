package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jorgerr9011/cartas-game-backend/internal/app/player"
)

type PlayerHandler struct {
	usecase player.UseCase
}

func NewPlayerHandler(uc player.UseCase) *PlayerHandler {
	return &PlayerHandler{usecase: uc}
}

func (h *PlayerHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/players", h.createPlayer)
}

type createPlayerReq struct {
	Name string `json:"name" binding:"required"`
}

func (h *PlayerHandler) createPlayer(c *gin.Context) {
	var req createPlayerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	player, err := h.usecase.CreatePlayer(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, player)
}

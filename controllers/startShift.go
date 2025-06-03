package controllers

import (
	"dinerFlow/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// StartShift godoc
// @Summary Start user shift
// @Description Starts a new shift for the authenticated user
// @Tags shift
// @Produce json
// @Success 201 {object} models.MessageResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /startShift [post]
func StartShift(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao converter userID"})
		return
	}

	startTime := time.Now()

	_, err := config.DB.Exec("INSERT INTO start_shift (user_id, start_time) VALUES (?, ?)", userIDUint, startTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar turno"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Turno iniciado com sucesso"})
}

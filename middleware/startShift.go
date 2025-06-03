package middleware

import (
	"dinerFlow/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckStartShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
			c.Abort()
			return
		}

		userIDUint, ok := userID.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao converter userID"})
			c.Abort()
			return
		}

		var existsFlag bool
		err := config.DB.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM start_shift
				WHERE user_id = ? AND end_time IS NULL
			)
		`, userIDUint).Scan(&existsFlag)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar turno"})
			c.Abort()
			return
		}

		if !existsFlag {
			c.JSON(http.StatusForbidden, gin.H{"error": "Funcionário não iniciou o turno"})
			c.Abort()
			return
		}

		c.Next()
	}
}

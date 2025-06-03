package controllers

import (
	"dinerFlow/config"
	"dinerFlow/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

// CloseDay godoc
// @Summary Gera o relatório de fechamento do dia
// @Description Recupera todos os pagamentos fechados na data atual, gera um arquivo de relatório e salva em disco
// @Tags Diner Reports
// @Produce json
// @Success 200 {object} map[string]interface{} "Fechamento do dia gerado com sucesso"
// @Failure 401 {object} map[string]string "Usuário não autenticado"
// @Failure 500 {object} map[string]interface{} "Erro ao gerar relatório de fechamento"
// @Router /closeDay [get]
func CloseDay(c *gin.Context) {
	// Fetch today's date in YYYY-MM-DD format
	today := time.Now().Format("2006-01-02")

	// Prepare the query to get all payments closed today
	rows, err := config.DB.Query(`
		SELECT u.firstname, u.lastname, p.total_price
		FROM diner_payments p
		JOIN users u ON u.id = p.closed_by
		WHERE DATE(p.closed_at) = ?
	`, today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar pagamentos", "details": err.Error()})
		return
	}
	defer rows.Close()

	var summaries []utils.PaymentSummary
	for rows.Next() {
		var firstName, lastName string
		var value float64
		if err := rows.Scan(&firstName, &lastName, &value); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar linha", "details": err.Error()})
			return
		}
		summaries = append(summaries, utils.PaymentSummary{
			User:  firstName + " " + lastName,
			Value: value,
		})
	}

	// Get the closing user name from database using userId from context
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var closedBy string
	err = config.DB.QueryRow("SELECT CONCAT(firstname, ' ', lastname) FROM users WHERE id = ?", userId).Scan(&closedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar nome do usuário", "details": err.Error()})
		return
	}

	// Generate file path
	now := time.Now()
	dir := fmt.Sprintf("public/%d/%02d/%02d", now.Year(), now.Month(), now.Day())
	filename := fmt.Sprintf("%s/closing_%s.txt", dir, now.Format("2006-01-02"))

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar diretório", "details": err.Error()})
		return
	}

	if err = utils.GenerateClosingReport(filename, summaries, closedBy, now); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar relatório de fechamento", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Fechamento do dia gerado com sucesso",
		"file":    filename,
	})
}

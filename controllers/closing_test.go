package controllers

import (
	"dinerFlow/config"
	"dinerFlow/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouterWithContext(handler gin.HandlerFunc, userId int) *gin.Engine {
	r := gin.Default()
	r.GET("/closeDay", func(c *gin.Context) {
		c.Set("userId", userId)
		handler(c)
	})
	return r
}

func TestCloseDay_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	config.DB = db

	today := time.Now().Format("2006-01-02")

	// Mock SELECT diner_payments
	mock.ExpectQuery("SELECT u.firstname, u.lastname, p.total_price FROM diner_payments p JOIN users u ON u.id = p.closed_by WHERE DATE\\(p.closed_at\\) = \\?").
		WithArgs(today).
		WillReturnRows(sqlmock.NewRows([]string{"firstname", "lastname", "total_price"}).
			AddRow("Nicolai", "Furtado", 54.60))

	// Mock SELECT CONCAT(firstname, lastname)
	mock.ExpectQuery("SELECT CONCAT\\(firstname, ' ', lastname\\) FROM users WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).
			AddRow("Nicolai Furtado"))

	// Stub report generation using GenerateClosingReportMock
	originalGenerateClosingReportMock := utils.GenerateClosingReportMock
	defer func() { utils.GenerateClosingReportMock = originalGenerateClosingReportMock }()
	utils.GenerateClosingReportMock = func(filename string, summaries []utils.PaymentSummary, closedBy string, now time.Time) error {
		_ = os.MkdirAll(filepath.Dir(filename), 0755)
		return os.WriteFile(filename, []byte("Fechamento OK"), 0644)
	}

	r := setupTestRouterWithContext(CloseDay, 1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/closeDay", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Fechamento do dia gerado com sucesso", response["message"])
	assert.NotEmpty(t, response["file"])

	// Cleanup generated file
	if filePath, ok := response["file"].(string); ok {
		_ = os.Remove(filePath)
	}
}

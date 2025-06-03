package controllers

import (
	"bytes"
	"database/sql"
	"dinerFlow/config"
	"dinerFlow/models"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func setupTabRouter(handler gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userId", 1)
		c.Next()
	})
	r.POST("/tab", handler)
	return r
}

func TestOpenTab_TableNotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	config.DB = db

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM diner_tables WHERE id = ?")).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	body, _ := json.Marshal(models.Tab{
		TableId:    99,
		ClientName: "John",
		Order:      models.OrderData{Items: []models.OrderItem{}},
	})
	req := httptest.NewRequest(http.MethodPost, "/tab", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := setupTabRouter(OpenTab)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Mesa não encontrada")
}

func TestOpenTab_ConflictOpenTab(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	config.DB = db

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM diner_tables WHERE id = ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT table_id FROM diner_tab WHERE table_id = ? AND status = 'open'")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"table_id"}).AddRow(1))

	body, _ := json.Marshal(models.Tab{
		TableId:    1,
		ClientName: "John",
		Order:      models.OrderData{Items: []models.OrderItem{}},
	})
	req := httptest.NewRequest(http.MethodPost, "/tab", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := setupTabRouter(OpenTab)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "comanda aberta")
}

func TestOpenTab_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/tab", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := setupTabRouter(OpenTab)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Dados inválidos")
}

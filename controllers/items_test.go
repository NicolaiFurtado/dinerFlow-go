package controllers

import (
	"bytes"
	"database/sql"
	"dinerFlow/config"
	"dinerFlow/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouterWithUserID(handler gin.HandlerFunc) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("userId", 1)
		c.Next()
	})
	r.POST("/items", handler)
	r.PUT("/items", handler)
	r.DELETE("/items", handler)
	r.GET("/items", handler)
	return r
}

func TestGetItems_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	config.DB = db

	mock.ExpectQuery("SELECT id, cod, name, description, price, category, created_by, updated_info FROM diner_items").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "cod", "name", "description", "price", "category", "created_by", "updated_info",
		}).AddRow(1, "12345", "Burger", "Delicious", 10.5, "Fast Food", 1, sql.NullString{String: "", Valid: false}))

	r := gin.Default()
	r.GET("/items", GetItems)
	req, _ := http.NewRequest("GET", "/items", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateItem_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	config.DB = db

	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM diner_items WHERE cod = \\?\\)").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec("INSERT INTO diner_items").
		WithArgs(sqlmock.AnyArg(), "Burger", "Delicious", 10.5, "Fast Food", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := models.Items{
		Name:        "Burger",
		Description: "Delicious",
		Price:       10.5,
		Category:    "Fast Food",
	}
	jsonBody, _ := json.Marshal(body)

	r := setupTestRouterWithUserID(CreateItem)
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateItem_InvalidBody(t *testing.T) {
	r := setupTestRouterWithUserID(CreateItem)
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEditItem_InvalidBody(t *testing.T) {
	r := setupTestRouterWithUserID(EditItem)
	req, _ := http.NewRequest("PUT", "/items", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteItem_InvalidBody(t *testing.T) {
	r := setupTestRouterWithUserID(DeleteItem)
	req, _ := http.NewRequest("DELETE", "/items", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

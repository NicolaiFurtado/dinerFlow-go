package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/tables", GetTables)
	r.POST("/tables", CreateTable)
	r.PUT("/tables", EditTable)
	r.DELETE("/tables", DeleteTable)
	return r
}

func TestCreateTable_InvalidBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/tables", bytes.NewBuffer([]byte(`invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router := setupTestRouter()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestEditTable_InvalidBody(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/tables", bytes.NewBuffer([]byte(`invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router := setupTestRouter()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestDeleteTable_InvalidBody(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/tables", bytes.NewBuffer([]byte(`invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router := setupTestRouter()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

package controllers

import (
	"database/sql"
	"dinerFlow/config"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStartShift_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	config.DB = db

	mock.ExpectExec("INSERT INTO start_shift").
		WithArgs(uint(1), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	router := gin.Default()
	router.POST("/startShift", func(c *gin.Context) {
		c.Set("userID", uint(1))
		StartShift(c)
	})

	req, _ := http.NewRequest("POST", "/startShift", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "Turno iniciado com sucesso")
}

func TestStartShift_Unauthenticated(t *testing.T) {
	router := gin.Default()
	router.POST("/startShift", StartShift)

	req, _ := http.NewRequest("POST", "/startShift", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Usuário não autenticado")
}

func TestStartShift_InvalidUserID(t *testing.T) {
	router := gin.Default()
	router.POST("/startShift", func(c *gin.Context) {
		c.Set("userID", "not-an-uint")
		StartShift(c)
	})

	req, _ := http.NewRequest("POST", "/startShift", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), "Erro ao converter userID")
}

func TestStartShift_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	config.DB = db

	mock.ExpectExec("INSERT INTO start_shift").
		WithArgs(uint(1), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	router := gin.Default()
	router.POST("/startShift", func(c *gin.Context) {
		c.Set("userID", uint(1))
		StartShift(c)
	})

	req, _ := http.NewRequest("POST", "/startShift", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), "Erro ao iniciar turno")
}

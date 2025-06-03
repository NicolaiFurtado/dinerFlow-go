package controllers_test

import (
	"bytes"
	"dinerFlow/config"
	"dinerFlow/controllers"
	"dinerFlow/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)
	return r
}

func TestSignUp_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	config.DB = db

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM users WHERE username = \?\)`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec("INSERT INTO users").
		WithArgs("testuser", "test@example.com", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := models.Users{Username: "testuser", Email: "test@example.com", Password: "pass123"}
	jsonBody, _ := json.Marshal(body)

	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Usuário criado com sucesso")
}

func TestLogin_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	config.DB = db

	hashed, _ := controllers.HashPassword("pass123")
	mock.ExpectQuery(`SELECT id, password_hash FROM users WHERE username = \?`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).AddRow(1, hashed))

	body := models.Users{Username: "testuser", Password: "pass123"}
	jsonBody, _ := json.Marshal(body)

	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestLogin_InvalidPassword(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	config.DB = db

	hashed, _ := controllers.HashPassword("pass123")
	mock.ExpectQuery(`SELECT id, password_hash FROM users WHERE username = \?`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).AddRow(1, hashed))

	body := models.Users{Username: "testuser", Password: "wrongpass"}
	jsonBody, _ := json.Marshal(body)

	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Credenciais inválidas")
}

func TestLogout(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Logout realizado com sucesso")
}

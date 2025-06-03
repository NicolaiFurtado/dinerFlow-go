package controllers

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword godoc
// @Summary Hash a plain-text password
// @Description Uses bcrypt to securely hash a plain-text password
// @Tags auth
// @Param password body string true "Plain-text password to hash"
// @Success 200 {string} string "Hashed password"
// @Failure 500 {object} models.ErrorResponse
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// CheckPasswordHash godoc
// @Summary Verify password hash
// @Description Compares a plain-text password with a hashed password
// @Tags auth
// @Param password body string true "Plain-text password"
// @Param hash body string true "Bcrypt hashed password"
// @Success 200 {boolean} boolean "True if password matches hash"
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

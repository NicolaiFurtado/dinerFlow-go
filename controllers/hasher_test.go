package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword_Success(t *testing.T) {
	password := "mySecret123"
	hashed, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed, "Hashed password should not match the original")
}

func TestCheckPasswordHash_ValidMatch(t *testing.T) {
	password := "mySecret123"
	hashed, _ := HashPassword(password)

	isValid := CheckPasswordHash(password, hashed)
	assert.True(t, isValid, "Expected password to match the hash")
}

func TestCheckPasswordHash_InvalidMatch(t *testing.T) {
	password := "mySecret123"
	wrongPassword := "wrongPass123"
	hashed, _ := HashPassword(password)

	isValid := CheckPasswordHash(wrongPassword, hashed)
	assert.False(t, isValid, "Expected wrong password to not match the hash")
}

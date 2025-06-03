package models

type AuthRequest struct {
	Username string `json:"username" example:"john"`
	Password string `json:"password" example:"123456"`
}

type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOi..."`
}

type MessageResponse struct {
	Message string `json:"message" example:"Usuário criado com sucesso"`
}

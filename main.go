// @title DinerFlow API
// @version 1.0
// @description API for managing diner tables
// @host localhost:8080
// @BasePath /
// @schemes http

package main

import (
	"dinerFlow/config"
	"dinerFlow/middleware"
	"dinerFlow/routes"

	_ "dinerFlow/docs" // ✅ importação do Swagger gerado

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env file not found. Using system environment variables.")
	}
}

func main() {
	loadEnv()
	config.Connect()

	r := gin.Default()

	// ✅ Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rotas públicas
	routes.PublicRoutes(r)

	// Rotas protegidas por JWT
	protected := r.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	routes.ProtectedRoutes(protected)

	r.Run(":8080")
}

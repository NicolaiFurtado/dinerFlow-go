package routes

import (
	"dinerFlow/controllers"
	"dinerFlow/middleware"

	"github.com/gin-gonic/gin"
)

func PublicRoutes(router *gin.Engine) {
	router.POST("/signup", controllers.SignUp)
	router.POST("/login", controllers.Login)
	router.POST("/logout", controllers.Logout)
}

func ProtectedRoutes(router *gin.RouterGroup) {
	/**
	* Rotas para gerenciamento de mesas - SEM NECESSIDADE DE TURNO
	 */
	router.GET("/tables", controllers.GetTables)

	/**
	* Rotas para gerenciamento de itens - SEM NECESSIDADE DE TURNO
	 */
	router.GET("/items", controllers.GetItems)

	router.Use(middleware.CheckIfLogged())
	router.POST("/startShift", controllers.StartShift)

	protectedWithShift := router.Group("/")
	protectedWithShift.Use(middleware.CheckStartShift())
	{
		/**
		* Rotas para gerenciamento de mesas - COM NECESSIDADE DE TURNO
		 */
		protectedWithShift.POST("/createTable", controllers.CreateTable)
		protectedWithShift.PUT("/editTable", controllers.EditTable)
		protectedWithShift.DELETE("/deleteTable", controllers.DeleteTable)

		/**
		* Rotas para gerenciamento de itens - COM NECESSIDADE DE TURNO
		 */
		protectedWithShift.POST("/createItem", controllers.CreateItem)
		protectedWithShift.PUT("/editItem", controllers.EditItem)
		protectedWithShift.DELETE("/deleteItem", controllers.DeleteItem)

		/**
		* Rotas para gerenciamento de comandas - COM NECESSIDADE DE TURNO
		 */
		protectedWithShift.POST("/openTab", controllers.OpenTab)
		protectedWithShift.PUT("/updateOrder", controllers.UpdateOrder)
		protectedWithShift.PUT("/removeOrder", controllers.RemoveOrder)
		protectedWithShift.PUT("/closeTab", controllers.CloseTab)
		protectedWithShift.PUT("/finishPayment", controllers.FinishPayment)

		/**
		* Rotas para gerenciamento de fechamento de dia - COM NECESSIDADE DE TURNO
		 */
		protectedWithShift.GET("/closeDay", controllers.CloseDay)
	}
}

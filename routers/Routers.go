package routers

import (
	"API-Books/controllers"
	"API-Books/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://ferytell.github.io"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}
	router.Use(cors.New(config))

	router.GET("/api/ping", controllers.Hellow)
	// User Auth Routes
	router.POST("/api/signup", controllers.SignUp)
	router.POST("/api/login", controllers.Login)
	router.GET("/api/validate", middleware.RequireAuth, controllers.Validate)
	router.POST("/api/logout", controllers.Logout)

	// Books Routes
	router.GET("/api/books", controllers.GetAllBooks)
	router.GET("/api/books/:id", controllers.GetBook)
	router.POST("/api/books", controllers.CreateBook)
	router.PUT("/api/books/:id", controllers.UpdateBook)
	router.DELETE("/api/books/:id", controllers.DeleteBook)

	//Villager Routes
	router.POST("/api/villagers", controllers.CreateVillager)
	router.GET("/api/villagers", controllers.GetVillagers)
	router.GET("/api/villagers/:id", controllers.GetVillager)
	router.PUT("/api/villagers/:id", controllers.UpdateVillager)
	router.DELETE("/api/villagers/:id", controllers.DeleteVillager)

	return router
}

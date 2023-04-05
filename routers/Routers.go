package routers

import (
	"API-Books/controlers"
	"API-Books/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	// Enable CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	router.Use(cors.New(config))

	router.POST("/api/signup", controlers.SignUp)
	router.POST("/api/login", controlers.Login)
	router.GET("/api/validate", middleware.RequireAuth, controlers.Validate)
	router.POST("/api/logout", controlers.Logout)

	router.GET("/api/ping", controlers.Hellow)

	router.GET("/api/books", controlers.GetAllBooks)
	router.GET("/api/books/:id", controlers.GetBook)
	router.POST("/api/books", controlers.CreateBook)
	router.PUT("/api/books/:id", controlers.UpdateBook)
	router.DELETE("/api/books/:id", controlers.DeleteBook)

	//router.DELETE("/books", controlers.TestDel)

	return router
}

package routers

import (
	"API-Books/controlers"
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


	// Enable CORS
	// router.Use(func(c *gin.Context) {
	// 	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// 	c.Next()
	// })
	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost:3000", "https://ferytell.github.io"}

	// router.Use(cors.New(config))

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

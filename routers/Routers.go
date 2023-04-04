package routers

import (
	"API-Books/controlers"
	"API-Books/middleware"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.POST("/signup", controlers.SignUp)
	router.POST("/login", controlers.Login)
	router.GET("/validate", middleware.RequireAuth, controlers.Validate)

	router.GET("/ping", controlers.Hellow)

	router.GET("/books", controlers.GetAllBooks)

	router.GET("/books/:id", controlers.GetBook)

	router.POST("/books", controlers.CreateBook)

	router.PUT("/books/:id", controlers.UpdateBook)

	router.DELETE("/books/:id", controlers.DeleteBook)

	//router.DELETE("/books", controlers.TestDel)

	return router
}

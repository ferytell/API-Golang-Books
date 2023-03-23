package routers

import (
	"API-Books/controlers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", controlers.Hellow)

	router.GET("/books", controlers.GetAllBooks)

	router.GET("/books/:bookId", controlers.GetBook)

	router.POST("/books", controlers.CreateBook)

	router.PUT("/books/:bookId", controlers.UpdateBook)

	router.DELETE("/books/:bookId", controlers.DeleteBook)

	router.DELETE("/books", controlers.TestDel)

	return router
}

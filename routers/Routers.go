package routers

import (
	"API-Books/controlers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.POST("/books", controlers.CreateBook)

	router.GET("/ping", controlers.Hellow)

	router.PUT("/books/:bookId", controlers.UpdateBook)

	router.GET("/books/:bookId", controlers.GetBook)

	return router
}

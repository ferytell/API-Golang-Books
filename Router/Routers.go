package Routers

import (
	"./Controlers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.POST("/books", Controlers.CreateBook)

	return router
}

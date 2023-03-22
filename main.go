package main

import (
	"API-Books/initializer"
	"API-Books/routers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	initializer.LoadEnvVar()
	initializer.ConnectToDB()
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	routers.StartServer().Run()
}

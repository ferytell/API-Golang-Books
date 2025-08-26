package main

import (
	_ "API-Books/docs"
	"API-Books/initializer"
	"API-Books/routers"
	"log"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializer.LoadEnvVar()
	initializer.ConnectToDB()
	initializer.SyncDatabase()
}

// @title UEM Syariah API
// @version 1.0
// @description API for managing villagers, donations, loans, and repayments.
// @host localhost:8080
// @BasePath /
func main() {

	

	r := routers.StartServer()
	// routers.StartServer().Run()
	// Swagger
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local dev
	}

	log.Printf("Starting server on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}

}

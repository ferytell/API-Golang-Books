package main

import (
	"API-Books/initializer"
	"API-Books/routers"
	"log"
	"os"
)

func init() {
	initializer.LoadEnvVar()
	initializer.ConnectToDB()
	initializer.SyncDatabase()
}

func main() {
	r := routers.StartServer()
	// routers.StartServer().Run()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local dev
	}

	log.Printf("Starting server on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}

}

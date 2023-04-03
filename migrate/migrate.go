package main

import (
	"API-Books/initializer"
	"API-Books/models"
)

func init() {
	initializer.LoadEnvVar()
	initializer.ConnectToDB()
}

func main() {
	initializer.DB.AutoMigrate(&models.Post{})
	initializer.DB.AutoMigrate(&models.User{})
}

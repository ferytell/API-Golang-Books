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
	initializer.DB.AutoMigrate(&models.Post{})
	initializer.DB.AutoMigrate(&models.Comment{})
	initializer.DB.AutoMigrate(&models.Photo{})
	initializer.DB.AutoMigrate(&models.SocialMedia{})
}

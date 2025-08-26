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
	initializer.DB.AutoMigrate(&models.Villager{})
	initializer.DB.AutoMigrate(&models.Infaq{})
	initializer.DB.AutoMigrate(&models.Loan{})
	initializer.DB.AutoMigrate(&models.Repayment{})
	initializer.DB.AutoMigrate(&models.Fund{})
}

package initializer

import "API-Books/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}

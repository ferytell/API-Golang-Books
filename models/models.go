package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title       string
	Author      string
	Description string
	Body        string
}

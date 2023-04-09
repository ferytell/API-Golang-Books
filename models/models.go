package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title       string
	Author      string
	Description string
	Body        string
}

type Photo struct {
	gorm.Model
	Title     string `json:"title" validate:"required"`
	Caption   string `json:"caption" validate:"required"`
	PhotoURL  string `json:"photo_url" validate:"required,url"`
	UserID    uint   `json:"user_id" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Comment struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	PhotoID   uint      `json:"photo_id"`
	Message   string    `json:"message" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SocialMedia struct {
	gorm.Model
	Name           string `json:"name"`
	SocialMediaURL string `json:"social_media_url"`
	UserID         uint   `json:"user_id"`
}

func (c *Comment) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

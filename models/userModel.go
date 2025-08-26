package models

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" gorm:"unique" validate:"required,email"`
	Password string `json:"-" validate:"required,min=6"`
	Age      int    `json:"age" validate:"min=8"`
}


func (u *User) Validate() error {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		var errMsgs []string
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, err.Error())
		}
		return fmt.Errorf(strings.Join(errMsgs, ", "))
	}
	return nil
}

// err := validate.Struct(mystruct)
// validationErrors := err.(validator.ValidationErrors)

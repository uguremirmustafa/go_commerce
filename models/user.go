package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=32"`
}

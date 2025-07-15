package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email           string
	UserName        string
	Password        string
}

func (User) TableName() string {
	return "users"
}

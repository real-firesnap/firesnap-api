package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	User    User
	UserID  int
	Content string
}

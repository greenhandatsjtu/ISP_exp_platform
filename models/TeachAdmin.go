package models

import "github.com/jinzhu/gorm"

//teach admin
type TeachAdmin struct {
	gorm.Model
	User   User `json:"user"`
	UserID uint `json:"user_id"`
}

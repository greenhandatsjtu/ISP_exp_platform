package models

import "github.com/jinzhu/gorm"

//system admin
type SysAdmin struct {
	gorm.Model
	User   User `json:"user"`
	UserID uint `json:"user_id"`
}

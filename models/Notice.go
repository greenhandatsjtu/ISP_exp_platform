package models

import "github.com/jinzhu/gorm"

type Notice struct {
	gorm.Model
	Title  string `gorm:"not_null" form:"title" json:"title" binding:"required" example:"test title"`
	Body   string `gorm:"size:1000" form:"body" json:"body" binding:"required" example:"here is some content"`
	Author string `form:"author" json:"author" example:"Sun Hengke"`
}

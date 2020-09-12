package models

import (
	"github.com/jinzhu/gorm"
)

type UserResource struct {
	gorm.Model
	User         User       `json:"user"`
	UserID       uint       `gorm:"not_null" json:"user_id"`
	Experiment   Experiment `json:"experiment"`
	ExperimentID uint       `gorm:"not_null" json:"experiment_id"`
	Port         int32      `gorm:"not_null" json:"port"`
}

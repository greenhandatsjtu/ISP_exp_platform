package models

import (
	"backend/database"
	"github.com/jinzhu/gorm"
	"time"
)

type Experiment struct {
	gorm.Model
	CourseID   uint      `json:"course_id" form:"course_id" binding:"required"`
	EName      string    `gorm:"type:varchar(100);not_null" json:"name" form:"name" binding:"required"`
	ETime      time.Time `json:"time" form:"time" binding:"required"`
	Assignment string    `gorm:"type:varchar(500);not_null" form:"assignment" json:"assignment" binding:"required"`
	Files      []File    `form:"files" json:"files"` // 实验指导书、预习报告、实验报告模板
	Upload     bool      `json:"upload"`             //yaml uploaded
	Enable     bool      `json:"enable"`             // 实验启用
}

type File struct {
	gorm.Model
	ExperimentId uint   `json:"experiment_id" form:"experiment_id"`
	FileName     string `gorm:"type:varchar(200);not_null" json:"file_name" form:"file_name" binding:"required"`
}

func (experiment Experiment) GetCourse() Course {
	var course Course
	database.Db.Model(&experiment).Related(&course)
	return course
}

func (experiment *Experiment) AddFile(file *File) error {
	return database.Db.Model(&experiment).Association("Files").Append(file).Error
}

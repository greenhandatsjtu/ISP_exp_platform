package models

import (
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/database"
	"github.com/jinzhu/gorm"
)

type Student struct {
	gorm.Model
	User    User      `json:"user"`
	UserID  uint      `json:"user_id"`
	SNo     string    `gorm:"type:varchar(50);not_null;unique" json:"student_number"`
	CNo     string    `gorm:"type:varchar(50);not_null" json:"class_number"`
	Courses []*Course `gorm:"many2many:student_courses" json:"courses"`
}

func (student Student) GetCourses() []Course {
	var courses []Course
	database.Db.Model(&student).Preload("Teachers").Preload("Assistants").Preload("Experiments").Related(&courses, "Courses")
	return courses
}

func (student Student) GetExperimentGrades() []StudentExperiment {
	var grades []StudentExperiment
	database.Db.Model(&student).Preload("Experiment").Related(&grades)
	return grades
}

func (student Student) GetCourseGrades() []StudentCourse {
	var grades []StudentCourse
	database.Db.Model(&student).Preload("Course").Related(&grades)
	return grades
}

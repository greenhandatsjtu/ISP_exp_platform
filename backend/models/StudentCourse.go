package models

import (
	"backend/database"
)

//CourseGrade
type StudentCourse struct {
	//gorm.Model
	Student   Student `json:"student"`
	StudentID uint    `gorm:"not_null" json:"student_id"`
	Course    Course  `json:"course"`
	CourseID  uint    `gorm:"not_null" json:"course_id"`
	Grade     int     `json:"grade"`
}

func (grade *StudentCourse) Update(newGrade int) {
	database.Db.Model(&grade).Update(&StudentCourse{Grade: newGrade})
}

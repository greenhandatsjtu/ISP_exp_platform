package models

import (
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/database"
	"github.com/jinzhu/gorm"
)

type Teacher struct {
	gorm.Model
	User    User      `json:"user"`
	UserID  uint      `json:"user_id"`
	Courses []*Course `gorm:"many2many:teacher_courses" json:"courses"`
}

func (teacher Teacher) GetCourses() []Course {
	var courses []Course
	database.Db.Model(&teacher).Related(&courses, "Courses")
	return courses
}

func (teacher Teacher) GetCourse(cno int) (Course, error) {
	var courses []Course
	var course Course
	if err := database.Db.Model(&teacher).Related(&courses, "Courses").Where(cno).First(&course).Error; err != nil {
		return Course{}, err
	}
	return course, nil
}

func (teacher Teacher) GetExperiment(eno int) (Experiment, error) {
	var experiments []Experiment
	var experiment Experiment
	if err := database.Db.Model(&teacher).Related(&experiments, "Experiments").Where(eno).First(&experiment).Error; err != nil {
		return Experiment{}, err
	}
	return experiment, nil
}

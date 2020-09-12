package models

import (
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/database"
	"github.com/jinzhu/gorm"
)

type Assistant struct {
	gorm.Model
	User   User `json:"user"`
	UserID uint `json:"user_id"`
	//SNo     string    `gorm:"type:varchar(50);not_null;unique" json:"s_no"`
	Courses []*Course `gorm:"many2many:assistant_courses" json:"courses"`
}

func (assistant Assistant) GetCourses() []Course {
	var courses []Course
	database.Db.Model(&assistant).Related(&courses, "Courses")
	return courses
}

func (assistant Assistant) GetExperiments(eno int) (Experiment, error) {
	var experiments []Experiment
	var experiment Experiment
	if err := database.Db.Model(&assistant).Related(&experiments, "Experiments").Where(eno).First(&experiment).Error; err != nil {
		return Experiment{}, err
	}
	return experiment, nil
}

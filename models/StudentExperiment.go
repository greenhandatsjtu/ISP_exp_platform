package models

import (
	"backend/database"
	"github.com/jinzhu/gorm"
	"time"
)

//ExperimentGrade
type StudentExperiment struct {
	gorm.Model
	Student      Student    `json:"student"`
	StudentID    uint       `gorm:"not_null" json:"student_id"`
	Experiment   Experiment `json:"experiment"`
	ExperimentID uint       `gorm:"not_null" json:"experiment_id"`
	Grade        *int       `json:"grade" json:"grade"`
	FileName     string     `gorm:"type:varchar(100);not_null" json:"file_name"`
	UploadAt     time.Time  `json:"upload_at"`
}

func (grade *StudentExperiment) UpdateGrade(newGrade int) {
	database.Db.Model(&grade).Update(&StudentExperiment{Grade: &newGrade})
}

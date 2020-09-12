package models

import (
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/database"
	"github.com/jinzhu/gorm"
)

type Course struct {
	gorm.Model
	CName       string       `gorm:"type:varchar(100);not_null" json:"course_name" form:"course_name" binding:"required"`
	Experiments []Experiment `form:"experiments" json:"experiments"`
	Assistants  []*Assistant `gorm:"many2many:assistant_courses" form:"assistants" json:"assistants"`
	Teachers    []*Teacher   `gorm:"many2many:teacher_courses" form:"teachers" json:"teachers"`
	Students    []*Student   `gorm:"many2many:student_courses" form:"students" json:"students"`
}

type UpdateCourseStruct struct {
	Tno   []uint `json:"teachers" form:"teachers" example:"1,3,7"`
	Ano   []uint `json:"assistants" form:"assistants" example:"3,4"`
	Sno   []int  `json:"students" form:"students" example:"517021910443,517021910444"`
	CName string `json:"course_name" form:"course_name" example:"test course name"`
}

type Assignment struct {
	Tno []uint `json:"teachers" form:"teachers" example:"1,3,7"`
	Ano []uint `json:"assistants" form:"assistants" example:"3,4"`
	Sno []int  `json:"students" form:"students" example:"517021910443,517021910444"`
	Cno uint   `json:"course" form:"course" example:"10" binding:"required"`
}

type TeacherAssignment struct {
	Tno []uint `json:"teachers" form:"teachers" example:"1,3,7"`
}

type StudentAssignment struct {
	Sno []int `json:"students" form:"students" example:"517021910443,517021910444"`
}

type AssistantAssignment struct {
	Ano []uint `json:"assistants" form:"assistants" example:"3,4"`
}

func (course Course) GetStudents() []Student {
	var students []Student
	database.Db.Model(&course).Preload("User").Related(&students, "Students")
	return students
}

func (course Course) GetExperiments() []Experiment {
	var experiments []Experiment
	database.Db.Model(&course).Related(&experiments, "Experiments")
	return experiments
}

func (course *Course) AddExperiment(experiment *Experiment) {
	database.Db.Model(&course).Association("Experiments").Append(experiment)
}

func (course *Course) Update(updateInfo UpdateCourseStruct) {
	var teachers []Teacher
	var students []Student
	var assistants []Assistant
	database.Db.Where(updateInfo.Tno).Find(&teachers)
	database.Db.Where(updateInfo.Ano).Find(&assistants)
	database.Db.Where("s_no in (?)", updateInfo.Sno).Find(&students)
	database.Db.Model(&course).Update(Course{CName: updateInfo.CName})
	database.Db.Model(&course).Association("Students").Clear().Append(students)
	database.Db.Model(&course).Association("Teachers").Clear().Append(teachers)
	database.Db.Model(&course).Association("Assistants").Clear().Append(assistants)
	database.Db.Save(&course)
}

func (course *Course) Assign(assignment Assignment) {
	var teachers []Teacher
	var students []Student
	var assistants []Assistant
	database.Db.Where(assignment.Tno).Find(&teachers)
	database.Db.Where(assignment.Ano).Find(&assistants)
	database.Db.Where("s_no in (?)", assignment.Sno).Find(&students)
	database.Db.Model(&course).Association("Students").Append(students)
	database.Db.Model(&course).Association("Teachers").Append(teachers)
	database.Db.Model(&course).Association("Assistants").Append(assistants)
}

func (course *Course) AssignAssistant(assignment AssistantAssignment) {
	var assistants []Assistant
	database.Db.Where(assignment.Ano).Find(&assistants)
	database.Db.Model(&course).Association("Assistants").Clear().Append(assistants)
}

func (course *Course) AssignTeacher(assignment TeacherAssignment) {
	var teachers []Teacher
	database.Db.Where(assignment.Tno).Find(&teachers)
	database.Db.Model(&course).Association("Teachers").Clear().Append(teachers)
}

func (course *Course) AssignStudent(assignment StudentAssignment) {
	var students []Student
	database.Db.Where(assignment.Sno).Find(&students)
	database.Db.Model(&course).Association("Students").Clear().Append(students)
}

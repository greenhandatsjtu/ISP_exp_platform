package controllers

import (
	"backend/models"
	"github.com/gin-gonic/gin"
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/database"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

//GetAllStudents godoc
// @Tags Student
// @Summary 获取所有学生
//@Description 以json格式返回所有学生信息
// @Produce  json
// @Success 200 {object} Response
// @Router /admin/student [get]
func GetAllStudents(c *gin.Context) {
	var students []models.Student
	if err := database.Db.Preload("User").Find(&students).Error; err != nil {
		NotFound(c, "student not found")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    students,
	})
}

//GetStudent godoc
// @Tags Student
// @Summary 获取单个学生
// @Produce json
// @Param id path int true "学号" default(517021910443)
// @Success 200 {object} Response
// @Router /admin/student/{id} [get]
func GetStudent(c *gin.Context) {
	id := c.Params.ByName("sno")
	var student models.Student
	if err := database.Db.Where(&models.Student{SNo: id}).Preload("User").First(&student).Error; err != nil {
		NotFound(c, "student not found")
	} else {
		c.JSON(http.StatusOK, Response{
			Success: true,
			Code:    http.StatusOK,
			Message: "",
			Data:    student,
		})
	}
}

//GetStudentCourses godoc
// @Tags Student
// @Summary 获取学生课程
// @Produce json
// @Param id path int true "学号" default(517021910443)
// @Success 200 {object} Response
// @Router /admin/student/{id}/course [get]
func GetStudentCourses(c *gin.Context) {
	id := c.Params.ByName("sno")
	var student models.Student
	if err := database.Db.Where(&models.Student{SNo: id}).First(&student).Error; err != nil {
		NotFound(c, "student not found")
	} else {
		c.JSON(http.StatusOK, Response{
			Success: true,
			Code:    http.StatusOK,
			Message: "",
			Data:    student.GetCourses(),
		})
	}
}

//GetStudentGrade godoc
// @Tags Student
// @Summary 获取学生成绩
// @Description 根据学号获取该生成绩
// @Produce json
// @Param id path int true "学号" default(517021910443)
// @Success 200 {object} Response
// @Router /admin/student/{id}/grade [get]
func GetStudentGrade(c *gin.Context) {
	sno := c.Params.ByName("sno")
	var student models.Student
	if err = database.Db.Where(models.Student{SNo: sno}).First(&student).Error; err != nil {
		NotFound(c, "student not found")
		return
	}
	//if err = database.Db.Where(models.StudentExperiment{StudentID:student.ID}).Preload("Experiment").Find(&grades).Error;err!=nil{
	//	NotFound(c,"Grade not found")
	//	return
	//}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    student.GetExperimentGrades(),
	})
}

//GetStudentCourseGrade godoc
// @Tags Student
// @Summary 获取学生某实验成绩
// @Description 根据学号和实验ID获取学生实验成绩
// @Produce json
// @Param sno path int true "学号" default(517021910443)
// @Param eno path int true "实验ID" default(1)
// @Success 200 {object} Response
// @Router /admin/student/{sno}/experiment/{eno}/grade [get]
func GetStudentExperimentGrade(c *gin.Context) {
	sno := c.Params.ByName("sno")
	eno := c.Params.ByName("eno")
	var student models.Student
	if err = database.Db.Where(models.Student{SNo: sno}).First(&student).Error; err != nil {
		NotFound(c, "student not found")
		return
	}
	var grades models.StudentExperiment
	if err = database.Db.Where("student_id=? AND experiment_id=?", student.ID, eno).Find(&grades).Error; err != nil {
		NotFound(c, "Grade not found")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    grades,
	})
}

// GetCoursesGrades godoc
// @Tags Student
// @Summary 获取学生本人所有实验成绩，仅学生能调用
// @Produce json
// @Success 200 {object} Response
// @Router /student/experiments/grades [get]
func GetExperimentGrades(c *gin.Context) {
	user, _ := c.Get("user")
	student, _ := user.(models.User).Student()
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    student.GetExperimentGrades(),
	})
}

// GetCourseGrade godoc
// @Tags Student
// @Summary 获取本人某实验成绩，仅学生能调用
// @Description 根据实验ID获取学生实验成绩
// @Produce json
// @Param eno path int true "实验ID" default(1)
// @Success 200 {object} Response
// @Router /student/experiment/{eno}/grade [get]
func GetExperimentGrade(c *gin.Context) {
	eno := c.Params.ByName("eno")
	user, _ := c.Get("user")
	student, _ := user.(models.User).Student()
	var grades models.StudentExperiment
	if err = database.Db.Where("student_id=? AND experiment_id=?", student.ID, eno).Find(&grades).Error; err != nil {
		NotFound(c, "Grade not found")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    grades,
	})
}

//GetAssignedCourses godoc
// @Tags Student
// @Summary 获取学生本人课程
// @Produce json
// @Success 200 {object} Response
// @Router /student/courses [get]
func GetAssignedCourses(c *gin.Context) {
	user, _ := c.Get("user")
	student, _ := user.(models.User).Student()
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    student.GetCourses(),
	})
}

// GetStudentCourseGrade godoc
// @Tags Student
// @Summary 获取学生某课程成绩，管理员调用
// @Description 根据学号和课程ID获取学生课程成绩
// @Produce json
// @Param sno path int true "学号" default(517021910443)
// @Param cno path int true "课程ID" default(1)
// @Success 200 {object} Response
// @Router /admin/student/{sno}/course/{cno}/grade [get]
func GetStudentCourseGrade(c *gin.Context) {
	sno := c.Params.ByName("sno")
	cno := c.Params.ByName("cno")
	var student models.Student
	if err = database.Db.Where(models.Student{SNo: sno}).First(&student).Error; err != nil {
		NotFound(c, "student not found")
		return
	}
	var grades models.StudentCourse
	if err = database.Db.Where("student_id=? AND course_id=?", student.ID, cno).Find(&grades).Error; err != nil {
		NotFound(c, "Grade not found")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    grades,
	})
}

// GetCoursesGrades godoc
// @Tags Student
// @Summary 获取学生本人所有课程成绩，仅学生能调用
// @Produce json
// @Success 200 {object} Response
// @Router /student/courses/grades [get]
func GetCoursesGrades(c *gin.Context) {
	user, _ := c.Get("user")
	student, _ := user.(models.User).Student()
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    student.GetCourseGrades(),
	})
}

// GetCourseGrade godoc
// @Tags Student
// @Summary 获取本人某课程成绩，仅学生能调用
// @Description 根据课程ID获取学生课程成绩
// @Produce json
// @Param cno path int true "课程ID" default(1)
// @Success 200 {object} Response
// @Router /student/course/{cno}/grade [get]
func GetCourseGrade(c *gin.Context) {
	cno := c.Params.ByName("cno")
	user, _ := c.Get("user")
	student, _ := user.(models.User).Student()
	var grade models.StudentCourse
	if err = database.Db.Where("student_id=? AND course_id=?", student.ID, cno).Find(&grade).Error; err != nil {
		NotFound(c, "Grade not found")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    grade,
	})
}

//批量\单个添加学生
func PostAddStudents(c *gin.Context) {
	var students []NewStudentStruct
	if err = c.Bind(&students); err != nil {
		log.Println(err)
		badRequest(c)
		return
	}

	var roles []models.Role
	database.Db.Where(models.Role{Description: "student"}).Find(&roles)

	for _, student := range students {
		if len(student.SNo) < 6 {
			continue
		}

		//generate password
		password := student.SNo[len(student.SNo)-6:]
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

		//insert into database
		database.Db.Create(&models.Student{
			User: models.User{
				Email:    student.Email,
				Password: string(hash),
				Name:     student.Name,
				Roles:    roles,
			},
			SNo: student.SNo,
			CNo: student.CNo,
		})
	}
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Students Created",
		Data:    nil,
	})
}

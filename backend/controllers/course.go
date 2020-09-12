package controllers

import (
	"backend/database"
	"backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
)

//GetAllCourses godoc
// @Tags Course
// @Summary 获取所有课程
//@Description 以json格式返回所有课程信息
// @Produce  json
// @Success 200 {object} Response
// @Router /course [get]
func GetAllCourses(c *gin.Context) {
	var courses []models.Course
	if err = database.Db.Find(&courses).Error; err != nil {
		NotFound(c, "course not found")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    courses,
	})
}

//GetCourse godoc
// @Tags Course
// @Summary 获取单个课程
// @Description 根据课程ID获取单个课程信息，同时返回课程的教师、助教、学生；但对于学生角色不返回学生信息。
// @Produce json
// @Param id path int true "课程ID" default(1)
// @Success 200 {object} Response
// @Router /course/{id} [get]
func GetCourse(c *gin.Context) {
	id := c.Params.ByName("id")
	roles, _ := c.Get("roles")
	var course models.Course
	if err = database.Db.Where(id).First(&course).Error; err != nil {
		NotFound(c, "course not found")
		return
	}
	auth := regexp.MustCompile("teacher|assistant|teach_admin|sys_admin")
	if auth.MatchString(strings.Join(roles.([]string), " ")) {
		database.Db.Where(id).Preload("Teachers.User").Preload("Assistants.User").Preload("Experiments").Preload("Students.User").First(&course)
	} else {
		//student role
		database.Db.Where(id).Preload("Teachers.User").Preload("Assistants.User").Preload("Experiments").First(&course)
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    course,
	})
}

//GetCourseStudents godoc
// @Tags Course
// @Summary 获取课程学生
// @Description 根据课程ID获取某一课程的所有学生
// @Produce json
// @Param id path int true "课程ID" default(1)
// @Success 200 {object} Response
// @Router /admin/course/{id}/student [get]
func GetCourseStudents(c *gin.Context) {
	id := c.Params.ByName("id")
	var course models.Course
	if err = database.Db.Where(id).First(&course).Error; err != nil {
		NotFound(c, "course not found")
	} else {
		c.JSON(http.StatusOK, Response{
			Success: true,
			Code:    http.StatusOK,
			Message: "",
			Data:    course.GetStudents(),
		})
	}
}

// PostCourse godoc
// @Tags Course
// @Summary 教务老师上传课程
// @Description 仅接受课程名。分配老师、助教、学生等请之后调用AssignCourse
// @Accept json
// @Produce  json
// @Param course_name formData string true "课程名"
// @Success 201 {object} Response
// @Router /admin/course [post]
func PostCourse(c *gin.Context) {
	var course models.Course
	if err = c.Bind(&course); err != nil || course.CName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "Bad request.",
			Data:    nil,
		})
	} else {
		database.Db.Create(&course)
		c.JSON(http.StatusCreated, Response{
			Success: true,
			Code:    http.StatusCreated,
			Message: "Course created.",
			Data:    course,
		})
	}
}

// PostUpdateCourse godoc
// @Tags Course
// @Summary 更新课程
// @Description 根据课程ID更新该课程，以json传入课程名、教师ID、助教ID、学号即可
// @Accept json
// @Produce  json
// @Param id path int true "课程ID" default(8)
// @Param updateInfo body models.UpdateCourseStruct true "更新信息"
// @Success 201 {object} Response
// @Router /admin/course/{id}/update [post]
func PostUpdateCourse(c *gin.Context) {
	id := c.Params.ByName("id")
	var course models.Course
	if err = database.Db.Where(id).First(&course).Error; err != nil {
		NotFound(c, "course not found")
		return
	}
	var updateInfo models.UpdateCourseStruct
	if err = c.Bind(&updateInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "Bad request.",
			Data:    nil,
		})
		return
	}

	//update course
	course.Update(updateInfo)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Assigned successfully.",
		Data:    course,
	})
}

//DeleteCourse godoc
// @Tags Course
// @Summary 删除课程
// @Description 根据课程ID删除课程
// @Produce json
// @Param id path int true "课程ID" default(1)
// @Success 200 {object} Response
// @Router /admin/course/{id}/delete [get]
func DeleteCourse(c *gin.Context) {
	id := c.Params.ByName("id")
	var course models.Course
	if err = database.Db.Where(id).First(&course).Error; err != nil {
		NotFound(c, "course not found")
		return
	}
	database.Db.Delete(&course)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "course " + id + " deleted",
		Data:    nil,
	})
}

// AssignStudents godoc
// @Tags Course
// @Summary 教务老师分配课程
// @Description 为该课程分配教师、助教和学生，可以分别分配；只能添加，若需删改教师，请调用UpdateCourse
// @Accept  json
// @Produce  json
// @Param assignment body models.Assignment true "分配信息"
// @Success 201 {object} models.Course
// @Router /admin/assign [post]
func AssignCourse(c *gin.Context) {
	var assignment models.Assignment
	var course models.Course
	if err = c.Bind(&assignment); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if err = database.Db.Where(assignment.Cno).Preload("Teachers").Preload("Students").Preload("Assistants").First(&course).Error; err != nil {
		NotFound(c, "course not found")
		return
	}
	course.Assign(assignment)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Assigned successfully.",
		Data:    course,
	})
}

//GetCourseExperiments godoc
// @Tags Course
// @Summary 获取课程实验
// @Description 根据课程ID获取该课程所有实验，以json格式返回
// @Produce json
// @Param id path int true "课程ID" default(1)
// @Success 200 {object} Response
// @Router /course/{id}/experiment [get]
func GetCourseExperiments(c *gin.Context) {
	id := c.Params.ByName("id") //course id
	var course models.Course
	if err = database.Db.Where(id).First(&course).Error; err != nil {
		NotFound(c, "course not found")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    course.GetExperiments(),
	})
}

// PostCourseExperiment godoc
// @Tags Course
// @Summary 上传课程实验
// @Description 上传实验，需要的参数有课程ID；实验名，实验时间，实验内容
// @Accept json
// @Produce  json
// @Param id path int true "课程ID"
// @Param experiment body ExperimentInfo true "实验"
// @Success 201 {object} Response
// @Router /admin/course/{id}/experiment [post]
func PostCourseExperiment(c *gin.Context) {
	id := c.Params.ByName("id")
	var course models.Course
	if err = database.Db.Where(id).First(&course).Error; err != nil {
		NotFound(c, "course not found")
		return
	}
	var experimentInfo ExperimentInfo
	if err = c.Bind(&experimentInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "Bad request.",
			Data:    nil,
		})
	} else {
		//etime, _ := now.Parse(experimentInfo.Time)
		experiment := models.Experiment{
			EName:      experimentInfo.Name,
			ETime:      experimentInfo.Time,
			Assignment: experimentInfo.Assignment,
		}
		course.AddExperiment(&experiment)
		c.JSON(http.StatusCreated, Response{
			Success: true,
			Code:    http.StatusCreated,
			Message: "Experiment created.",
			Data:    experiment,
		})
	}
}

//NOTE: new
// PostCourseGrade godoc
// @Tags Course
// @Summary 上传课程成绩
// @Description 上传课程成绩，需要的参数有学号、课程ID、成绩
// @Accept json
// @Produce  json
// @Param grade body CourseGradeInfo true "成绩信息"
// @Success 201 {object} Response
// @Router /admin/grade/course/upload [post]
func PostCourseGrade(c *gin.Context) {
	var gradeInfo CourseGradeInfo
	if err = c.Bind(&gradeInfo); err != nil {
		badRequest(c)
		return
	}
	var student models.Student
	if err = database.Db.Where(models.Student{SNo: gradeInfo.SNo}).First(&student).Error; err != nil {
		NotFound(c, "Student not found.")
		return
	}
	var course models.Course
	if err = database.Db.Where(gradeInfo.CNo).First(&course).Error; err != nil {
		NotFound(c, "Course not found.")
		return
	}
	grade := models.StudentCourse{
		Student: student,
		Course:  course,
		Grade:   gradeInfo.Grade,
	}
	database.Db.Create(&grade)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Grade uploaded.",
		Data:    grade,
	})
}

//NOTE: new
// UpdateCourseGrade godoc
// @Tags Course
// @Summary 更新课程成绩
// @Description 更新课程成绩，需要的参数有学号、实验ID、成绩
// @Accept json
// @Produce  json
// @Param grade body CourseGradeInfo true "成绩信息"
// @Success 201 {object} Response
// @Router /admin/grade/course/update [post]
func UpdateCourseGrade(c *gin.Context) {
	var gradeInfo CourseGradeInfo
	if err = c.Bind(&gradeInfo); err != nil {
		badRequest(c)
		return
	}
	var student models.Student
	if err = database.Db.Where(models.Student{SNo: gradeInfo.SNo}).First(&student).Error; err != nil {
		NotFound(c, "Student not found.")
		return
	}
	var course models.Course
	if err = database.Db.Where(gradeInfo.CNo).First(&course).Error; err != nil {
		NotFound(c, "Course not found.")
		return
	}
	var grade models.StudentCourse
	if err = database.Db.Where("student_id=? AND course_id=?", student.ID, course.ID).First(&grade).Error; err != nil {
		NotFound(c, "Grade not found.")
		return
	}
	grade.Update(gradeInfo.Grade)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Grade updated.",
		Data:    grade,
	})
}

//NOTE: new
// AssignStudents godoc
// @Tags Course
// @Summary 为课程分配教师
// @Description 为该课程分配教师（增、删）以json传入教师ID即可；会覆盖之前的分配，因此若想新增教师，需要将之前分配的教师ID传入
// @Accept  json
// @Produce  json
// @Param id path int true "课程ID" default(8)
// @Param assignment body models.TeacherAssignment true "分配教师信息"
// @Success 201 {object} models.Course
// @Router /admin/course/{id}/teacher [post]
func AssignTeacher(c *gin.Context) {
	var assignment models.TeacherAssignment
	var course models.Course
	id := c.Params.ByName("id")
	if err = c.Bind(&assignment); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if err = database.Db.Where(id).First(&course).Error; err != nil {
		NotFound(c, "course not found")
		return
	}
	course.AssignTeacher(assignment)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Assigned successfully.",
		Data:    course,
	})
}

// AssignStudents godoc
// @Tags Course
// @Summary 为课程分配学生
// @Description 为该课程分配学生（增、删）以json传入学生ID即可；会覆盖之前的分配
// @Accept  json
// @Produce  json
// @Param id path int true "课程ID" default(8)
// @Param assignment body models.StudentAssignment true "分配学生信息"
// @Success 201 {object} models.Course
// @Router /admin/course/{id}/students [post]
func AssignStudents(c *gin.Context) {
	var assignment models.StudentAssignment
	var course models.Course
	id := c.Params.ByName("id")
	if err = c.Bind(&assignment); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if err = database.Db.Where(id).First(&course).Error; err != nil {
		NotFound(c, "course not found")
		return
	}
	course.AssignStudent(assignment)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Assigned successfully.",
		Data:    course,
	})
}

//获取助教\教师负责的课程
func GetOwnCourses(c *gin.Context) {
	var teacher models.Teacher
	var assistant models.Assistant
	user, _ := c.Get("user")
	if teacher, err = user.(models.User).Teacher(); err == nil {
		c.JSON(http.StatusOK, Response{
			Success: true,
			Code:    http.StatusOK,
			Message: "",
			Data:    teacher.GetCourses(),
		})
	} else if assistant, err = user.(models.User).Assistant(); err == nil {
		c.JSON(http.StatusOK, Response{
			Success: true,
			Code:    http.StatusOK,
			Message: "",
			Data:    assistant.GetCourses(),
		})
	} else {
		Forbidden(c)
	}
}

//获取某课程所有助教
func GetCourseAssistants(c *gin.Context) {
	id := c.Params.ByName("id")
	var course models.Course
	var assistants []models.Assistant
	if err = database.Db.Where(id).First(&course).Error; err != nil {
		NotFound(c, "course not found")
		return
	}
	database.Db.Model(&course).Association("Assistants").Find(&assistants)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    assistants,
	})
}

//获取某课程所有教师
func GetCourseTeachers(c *gin.Context) {
	id := c.Params.ByName("id")
	var course models.Course
	var teachers []models.Teacher
	if err = database.Db.Where(id).First(&course).Error; err != nil {
		NotFound(c, "course not found")
		return
	}
	database.Db.Model(&course).Association("Teachers").Find(&teachers)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    teachers,
	})
}

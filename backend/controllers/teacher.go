package controllers

import (
	"backend/database"
	"backend/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

//GetAllTeachers godoc
// @Tags Teacher
// @Summary 获取所有老师
//@Description 以json格式返回所有老师信息
// @Produce  json
// @Success 200 {object} Response
// @Router /admin/teacher [get]
func GetAllTeachers(c *gin.Context) {
	var teachers []models.Teacher
	if err = database.Db.Preload("User").Find(&teachers).Error; err != nil {
		NotFound(c, "Teacher not found.")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    teachers,
	})
}

func GetTeacherExperiments(c *gin.Context) {
	id := c.Params.ByName("id")
	var teacher models.Teacher
	var experiments []models.Experiment
	if err = database.Db.Where(id).First(&teacher).Error; err != nil {
		NotFound(c, "Teacher not found.")
		log.Println(err)
	} else {
		database.Db.Model(&teacher).Related(&experiments, "Experiments")
		c.JSON(http.StatusOK, experiments)
	}
}

//GetTeacherCourses godoc
// @Tags Teacher
// @Summary 获取老师课程
// @Description 获取老师教授的所有课程，以json格式返回
// @Produce json
// @Param id path int true "教师ID" default(1)
// @Success 200 {object} Response
// @Router /admin/teacher/{id}/course [get]
func GetTeacherCourses(c *gin.Context) {
	var teacher models.Teacher
	id := c.Params.ByName("id")
	roles, _ := c.Get("roles")
	user, _ := c.Get("user")
	auth := regexp.MustCompile("teach_admin|sys_admin")
	if !auth.MatchString(strings.Join(roles.([]string), " ")) {
		if teacher, err = user.(models.User).Teacher(); err != nil {
			Forbidden(c)
			return
		}
		teacherId := strconv.Itoa(int(teacher.ID))
		if id != teacherId {
			Forbidden(c)
			return
		}
	}
	if err = database.Db.Where(id).First(&teacher).Error; err != nil {
		NotFound(c, "Teacher not found.")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    teacher.GetCourses(),
	})
}

// GetTeacherOwnCourses godoc
// @Tags Teacher
// @Summary 获取教师本人课程
// @Produce json
// @Success 200 {object} Response
// @Router /admin/teachercourses [get]
func GetTeacherOwnCourses(c *gin.Context) {
	var teacher models.Teacher
	user, _ := c.Get("user")
	if teacher, err = user.(models.User).Teacher(); err != nil {
		Forbidden(c)
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    teacher.GetCourses(),
	})
}

//NOTE: new
// @Summary 分配课程助教
// @Description 根据课程ID分配助教，以json传入助教ID即可，需要注意该操作会覆盖之前的助教，所以需要将所有助教ID都传入
// @Tags Teacher
// @Accept json
// @Produce  json
// @Param id path int true "课程ID" default(8)
// @Param assignInfo body models.AssistantAssignment true "助教信息"
// @Success 200 {object} Response
// @Router /admin/course/{id}/assistants [post]
func AssignAssistant(c *gin.Context) {
	var assignment models.AssistantAssignment
	var course models.Course
	id := c.Params.ByName("id")
	if err = c.Bind(&assignment); err != nil {
		badRequest(c)
		return
	}
	if err = database.Db.Where(id).First(&course).Error; err != nil {
		NotFound(c, "course not found")
		return
	}
	course.AssignAssistant(assignment)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Assigned successfully.",
		Data:    course,
	})
}

package controllers

import (
	"backend/models"
	"github.com/gin-gonic/gin"
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/database"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

//GetAssistantCourses godoc
// @Tags Assistant
// @Summary 获取助教课程
// @Description 获取助教负责的所有课程，以json格式返回
// @Produce json
// @Param id path int true "助教ID" default(1)
// @Success 200 {object} Response
// @Router /admin/assistant/{id}/course [get]
func GetAssistantCourses(c *gin.Context) {
	var assistant models.Assistant
	id := c.Params.ByName("id")
	roles, _ := c.Get("roles")
	user, _ := c.Get("user")
	auth := regexp.MustCompile("teach_admin|sys_admin")
	if !auth.MatchString(strings.Join(roles.([]string), " ")) {
		if assistant, err = user.(models.User).Assistant(); err != nil {
			Forbidden(c)
			return
		}
		assistantId := strconv.Itoa(int(assistant.ID))
		if id != assistantId {
			Forbidden(c)
			return
		}
	}
	if err = database.Db.Where(id).First(&assistant).Error; err != nil {
		NotFound(c, "Assistant not found.")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    assistant.GetCourses(),
	})
}

// GetAssistantOwnCourses godoc
// @Tags Assistant
// @Summary 获取助教本人课程
// @Description 获取助教本人负责的所有课程，以json格式返回
// @Produce json
// @Param id path int true "助教ID" default(1)
// @Success 200 {object} Response
// @Router /admin/assistantcourses [get]
func GetAssistantOwnCourses(c *gin.Context) {
	var assistant models.Assistant
	id := c.Params.ByName("id")
	user, _ := c.Get("user")
	if assistant, err = user.(models.User).Assistant(); err != nil {
		Forbidden(c)
		return
	}
	if err = database.Db.Where(id).First(&assistant).Error; err != nil {
		NotFound(c, "Assistant not found.")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    assistant.GetCourses(),
	})
}

//获取所有助教
func GetAllAssistants(c *gin.Context) {
	var assistants []models.Assistant
	if err = database.Db.Preload("User").Find(&assistants).Error; err != nil {
		NotFound(c, "Assistant not found.")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    assistants,
	})
}

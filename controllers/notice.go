package controllers

import (
	"backend/database"
	"backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

//GetAllNotices godoc
// @Tags Notice
// @Summary 获取所有通知
//@Description 以json格式返回所有通知
// @Produce  json
// @Success 200 {object} Response
// @Router /notice [get]
func GetAllNotices(c *gin.Context) {
	var notices []models.Notice
	if err = database.Db.Order("created_at desc").Find(&notices).Error; err != nil {
		NotFound(c, "Notice not found.")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    notices,
	})
}

//GetNoticegodoc
// @Tags Notice
// @Summary 获取单个通知
// @Produce json
// @Param id path int true "通知ID" default(1)
// @Success 200 {object} Response
// @Router /notice/{id} [get]
func GetNotice(c *gin.Context) {
	id := c.Params.ByName("id")
	var notice models.Notice
	if err = database.Db.Where(id).First(&notice).Error; err != nil {
		NotFound(c, "Notice not found.")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    notice,
	})
}

// PostNotice godoc
// @Tags Notice
// @Summary 上传通知
// @Description 上传通知，需要的参数有标题、内容、发布人（非必须）
// @Accpet json
// @Produce  json
// @Param notice body models.Notice true "通知"
// @Success 201 {object} Response
// @Router /admin/notice [post]
func PostNotice(c *gin.Context) {
	var notice models.Notice
	if err = c.Bind(&notice); err != nil {
		badRequest(c)
		return
	}
	user, _ := c.Get("user")
	if notice.Author == "" {
		notice.Author = user.(models.User).Name
	}
	database.Db.Save(&notice)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Notice created.",
		Data:    notice,
	})
}

//DeleteNotice godoc
// @Tags Notice
// @Summary 删除通知
// @Produce json
// @Param id path int true "通知ID" default(1)
// @Success 200 {object} Response
// @Router /admin/notice/{id}/delete [get]
func DeleteNotice(c *gin.Context) {
	id := c.Params.ByName("id")
	var notice models.Notice
	if err = database.Db.Where(id).First(&notice).Error; err != nil {
		NotFound(c, "Notice not found.")
		//log.Println(err)
		return
	}
	database.Db.Where(id).Delete(&notice)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "Notice " + id + " deleted.",
		Data:    notice,
	})
}

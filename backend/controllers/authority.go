package controllers

import (
	"backend/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/database"
	"net/http"
)

func AssignAuthority(c *gin.Context) {
	var auth Auth
	if err = c.Bind(&auth); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	var user models.User
	var roles []models.Role
	if err = database.Db.Where(auth.User).Preload("Roles").First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err = database.Db.Where("description in (?)", auth.Roles).Find(&roles).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	database.Db.Model(&user).Association("Roles").Append(roles)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    user,
	})
}

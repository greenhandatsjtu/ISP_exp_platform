package middlewares

import (
	"backend/database"
	"backend/models"
	"github.com/gin-gonic/gin"
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/controllers"
	"net/http"
	"time"
)

//同一个实验禁止同时开启多个资源
func CheckResourceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		eno := c.Params.ByName("eno")
		user, _ := c.Get("user")
		var resource models.UserResource
		if err := database.Db.Where(map[string]interface{}{"experiment_id": eno, "user_id": user.(models.User).ID}).First(&resource).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, controllers.Response{
				Success: false,
				Code:    http.StatusForbidden,
				Message: "Please end previous experiment.",
				Data:    nil,
			})
			return
		}
		c.Next()
	}
}

//实验开始前或yaml文件上传前不能做实验
func CheckCanExperimentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		eno := c.Params.ByName("eno")
		var experiment models.Experiment
		if err := database.Db.Where(eno).First(&experiment).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, controllers.Response{
				Success: false,
				Code:    http.StatusNotFound,
				Message: "Experiment not found.",
				Data:    nil,
			})
			return
		}
		now := time.Now()
		if !experiment.Upload || !experiment.Enable || now.Before(experiment.ETime) {
			c.AbortWithStatusJSON(http.StatusForbidden, controllers.Response{
				Success: false,
				Code:    http.StatusForbidden,
				Message: "Experiment is not accessible.",
				Data:    nil,
			})
			return
		}
		c.Next()
	}
}

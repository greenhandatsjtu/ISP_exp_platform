package middlewares

import (
	"backend/models"
	"github.com/gin-gonic/gin"
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/controllers"
	"net/http"
)

func CheckIsStudentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user")
		if _, err := user.(models.User).Student(); err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, controllers.Response{
				Success: false,
				Code:    http.StatusForbidden,
				Message: "You are not student.",
				Data:    nil,
			})
			return
		}
	}
}

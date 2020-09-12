package middlewares

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/controllers"
	"log"
	"net/http"
	"strings"
)

var Enforcer *casbin.Enforcer

func InitEnforcer() {
	adapter, _ := gormadapter.NewAdapterByDB(database.Db)
	Enforcer, _ = casbin.NewEnforcer("config/keymatch_model.conf", adapter)
	if err := Enforcer.LoadPolicy(); err != nil {
		log.Fatal(err)
	}

	//Init authority and save to database (table `casbin_role`)
	utils.InitAuthority(Enforcer)
}

func CasbinMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, _ := c.Get("roles")
		roles = strings.Join(roles.([]string), " ")
		p := c.Request.URL.Path
		m := c.Request.Method
		p = strings.TrimPrefix(p, "/api/admin")
		res, err := Enforcer.Enforce(roles, p, m)
		if err != nil || !res {
			c.AbortWithStatusJSON(http.StatusForbidden, controllers.Response{
				Success: false,
				Code:    http.StatusForbidden,
				Message: "You don't have permission.",
				Data:    nil,
			})
			return
		}
		c.Next()
		//for _, role := range roles.([]string) {
		//	res, err := Enforcer.Enforce(role, p, m)
		//	if err != nil {
		//		fmt.Println(err)
		//		continue
		//	}
		//	if !res {
		//		continue
		//	}
		//	c.Next()
		//	return
		//}
		//c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//check session
		session := sessions.Default(c)
		id := session.Get("id")
		if id == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.Response{
				Success: false,
				Code:    http.StatusUnauthorized,
				Message: "Please login.",
				Data:    nil,
			})
			return
		}
		var user models.User
		database.Db.Where(id).Preload("Roles").First(&user)
		//if err := database.Db.Where(id).Preload("Roles").First(&user).Error; err != nil {
		//	c.AbortWithStatusJSON(http.StatusNotFound,controllers.Response{
		//		Success:false,
		//		Code:    http.StatusNotFound,
		//		Message: "you haven't registered.",
		//		Data:    nil,
		//	})
		//	return
		//}
		c.Set("user", user)
		var roles []string
		for _, role := range user.Roles {
			roles = append(roles, role.Description)
		}
		c.Set("roles", roles)
		c.Next()
	}
}

package controllers

import (
	"backend/database"
	"backend/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"path/filepath"
)

// PostLogin godoc
// @Tags User
// @Summary 用户登录
// @Description 用户登录，调用接口必须先登录. Available account: sys_admin: YYGtHYK@Dsgdz.net;teach_admin: MoVNJQq@AYjoO.biz;teacher :HLxEYKX@WBsWQ.org; student: PenvKiV@AoMOE.org
// @Produce  json
// @Param email formData string true "用户邮箱" default(YYGtHYK@Dsgdz.net)
// @Param password formData string true "用户密码" default(test)
// @Success 200 {object} Response
// @Router /login [post]
func PostLogin(c *gin.Context) {
	var loginVals models.Login
	if err = c.Bind(&loginVals); err != nil {
		badRequest(c)
		return
	}
	var user models.User
	//auth
	if err = user.Auth(loginVals); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
			Success: false,
			Code:    http.StatusUnauthorized,
			Message: "Wrong email or password",
			Data:    nil,
		})
		return
	}
	//禁止同时登录同一个用户
	if *user.Online {
		Forbidden(c)
		return
	}
	session := sessions.Default(c)
	session.Set("id", user.ID) // set cookie and session
	if err = session.Save(); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "server error",
			Data:    nil,
		})
		return
	}
	database.Db.Model(&user).Update("Online", true) // set user online status
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "Login successfully.",
		Data:    user,
	})
}

// GetLogout godoc
// @Tags User
// @Summary 用户登出
// @Description 用户登出
// @Produce  json
// @Success 200 {object} Response
// @Router /logout [get]
func GetLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("id")
	if err = session.Save(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "server error",
			Data:    nil,
		})
		return
	}
	user, _ := c.Get("user")
	user.(models.User).Logout() //logout
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "Logout successfully.",
		Data:    nil,
	})
}

//PostChangePassword godoc
// @Tags User
// @Summary 修改密码
// @Description 修改用户密码
// @Accept mpfd
// @Produce json
// @Param pwd formData string true  "密码"
// @Success 201 {object} Response
// @Router /password [post]
func PostChangePassword(c *gin.Context) {
	password := c.PostForm("pwd")
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	user, _ := c.Get("user")
	user.(models.User).ChangePassword(string(hash))
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "change password successfully",
		Data:    nil,
	})
}

//GetPersonalInfo godoc
// @Tags User
// @Summary 获取本人信息
// @Description
// @Produce json
// @Success 200 {object} Response
// @Router /home [get]
func GetPersonalInfo(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
			Success: false,
			Code:    http.StatusUnauthorized,
			Message: "get user info error",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    user.(models.User),
	})
}

//GetUser godoc
// @Tags User
// @Summary 获取单个用户信息
// @Produce json
// @Param id path int true "用户ID" default(1)
// @Success 200 {object} Response
// @Router /admin/user/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user models.User
	if err = database.Db.Preload("Roles").Where(id).First(&user).Error; err != nil {
		NotFound(c, "User not found.")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    user,
	})
}

//GetOnlineUsers godoc
// @Tags User
// @Summary 获取当前在线用户
// @Produce json
// @Success 200 {object} Response
// @Router /admin/onlineuser [get]
func GetOnlineUsers(c *gin.Context) {
	var users []models.User
	if err = database.Db.Preload("Roles").Where("online=?", true).Find(&users).Error; err != nil {
		NotFound(c, "There are no users online.")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: fmt.Sprintf("there are %d users online", len(users)),
		Data:    users,
	})
}

//GetUserRole godoc
// @Tags User
// @Summary 获取用户的角色
// @Produce json
// @Param id path int true "用户ID" default(1)
// @Success 200 {object} Response
// @Router /admin/user/{id}/role [get]
func GetUserRole(c *gin.Context) {
	id := c.Params.ByName("id")
	var user models.User
	if err = database.Db.Where(id).First(&user).Error; err != nil {
		NotFound(c, "User not found.")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    user.GetRoles(),
	})
}

//Ping godoc
// @Summary ping
// @Produce json
// @Success 200 {object} Response
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(200, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "Everything is OK.",
		Data:    nil,
	})
}

//获取所有用户
func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err = database.Db.Preload("Roles").Find(&users).Error; err != nil {
		NotFound(c, "User not found.")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    users,
	})
}

type RoleDescription struct {
	Description string `json:"description" binding:"required"`
}

//为某用户添加角色
func AddRole(c *gin.Context) {
	id := c.Params.ByName("id")
	var user models.User
	var description RoleDescription
	var role models.Role
	if err = c.Bind(&description); err != nil {
		badRequest(c)
		log.Println(err)
		return
	}
	if err = database.Db.Where(models.Role{Description: description.Description}).First(&role).Error; err != nil {
		log.Println(err)
		NotFound(c, "Role not found.")
		return
	}
	if err = database.Db.Where(id).Preload("Roles").First(&user).Error; err != nil {
		log.Println(err)
		NotFound(c, "User not found.")
		return
	}

	//不能添加学生角色
	if role.Description == "student" {
		Forbidden(c)
		return
	}
	//只有学生才能作为助教
	if role.Description == "assistant" {
		if _, err := user.Student(); err != nil {
			Forbidden(c)
			return
		}
	}

	if err = database.Db.Model(&user).Association("Roles").Append(role).Error; err != nil {
		log.Println(err)
		badRequest(c)
		return
	}

	switch role.Description {
	case "sys_admin":
		_ = user.AssignSysAdminRole()
	case "teach_admin":
		_ = user.AssignTeachAdminRole()
	case "assistant":
		_ = user.AssignAssistantRole()
	case "teacher":
		_ = user.AssignTeacherRole()
	}
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Add role successfully",
		Data:    user,
	})
}

// 收回某用户的角色
func RevokeRole(c *gin.Context) {
	id := c.Params.ByName("id")
	description := c.Params.ByName("role")
	var user models.User
	var role models.Role
	var roles []models.Role
	if err = database.Db.Where(id).Preload("Roles").First(&user).Error; err != nil {
		log.Println(err)
		NotFound(c, "User not found.")
		return
	}
	if err = database.Db.Model(&user).Association("Roles").Find(&roles).Error; err != nil {
		log.Println(err)
		NotFound(c, "Roles not found.")
		return
	}
	for _, v := range roles {
		if v.Description == description {
			if err = database.Db.Where(models.Role{Description: description}).First(&role).Error; err != nil {
				log.Println(err)
				NotFound(c, "Role not found.")
				return
			}

			//不能删除学生角色
			if role.Description == "student" {
				Forbidden(c)
				return
			}

			if err = database.Db.Model(&user).Association("Roles").Delete(role).Error; err != nil {
				log.Println(err)
				badRequest(c)
				return
			}
			switch role.Description {
			case "sys_admin":
				database.Db.Where("user_id=?", id).Delete(models.SysAdmin{})
			case "teach_admin":
				database.Db.Where("user_id=?", id).Delete(models.TeachAdmin{})
			case "assistant":
				database.Db.Where("user_id=?", id).Delete(models.Assistant{})
			case "teacher":
				database.Db.Where("user_id=?", id).Delete(models.Teacher{})
			}
			c.JSON(http.StatusCreated, Response{
				Success: true,
				Code:    http.StatusCreated,
				Message: "Revoke role successfully",
				Data:    user,
			})
			return
		}
	}
	NotFound(c, "Role not found.")
}

//删除用户
func DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")
	database.Db.Where(id).Unscoped().Delete(models.User{})
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "",
		Data:    nil,
	})
}

func ChangeUserPassword(c *gin.Context) {
	id := c.Params.ByName("id")
	var user models.User
	if err = database.Db.Where(id).First(&user).Error; err != nil {
		log.Println(err)
		NotFound(c, "User not found.")
		return
	}
	password := c.PostForm("pwd")
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	user.ChangePassword(string(hash))
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Change password successfully",
		Data:    nil,
	})
}

type NewUserStruct struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

//添加用户
func AddUser(c *gin.Context) {
	var newUser NewUserStruct
	if err = c.Bind(&newUser); err != nil {
		log.Println(err)
		badRequest(c)
		return
	}
	log.Println(newUser.Password)
	//generate password
	hash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.MinCost)

	//insert into database
	database.Db.Create(&models.User{
		Email:    newUser.Email,
		Password: string(hash),
		Name:     newUser.Name,
	})
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "User Created",
		Data:    nil,
	})
}

//下载添加学生的excel模板
func DownloadTemplate(c *gin.Context) {
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "模板.xlsx"))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(filepath.Join("static", "template.xlsx"))
}

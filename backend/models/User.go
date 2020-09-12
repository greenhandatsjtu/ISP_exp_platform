package models

import (
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/database"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(50);not_null;unique" json:"email"`
	Password string `gorm:"type:varchar(200);not_null" json:"-"`
	Roles    []Role `gorm:"many2many:user_roles" json:"roles"`
	Name     string `gorm:"type:varchar(50);not_null" json:"name"`
	Online   *bool  `gorm:"default:false" json:"online"`
}

type Login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserInfo struct {
	SysAdmin   []SysAdmin
	TeachAdmin []TeachAdmin
	Teacher    []Teacher
	Assistant  []Assistant
	Student    []Student
}

func (user *User) Auth(loginVals Login) error {
	if err := database.Db.Where(&User{Email: loginVals.Email}).Preload("Roles").First(&user).Error; err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginVals.Password)); err != nil {
		return err
	}
	return nil
}

func (user User) Logout() {
	database.Db.Model(&user).Update("Online", false)
}

func (user User) GetRoles() []Role {
	var roles []Role
	database.Db.Model(&user).Related(&roles, "Roles")
	return roles
}

func (user User) GetInfo() UserInfo {
	var userInfo UserInfo
	database.Db.Model(&user).Preload("User").Related(&userInfo.SysAdmin).Related(&userInfo.TeachAdmin).Related(&userInfo.Teacher).Related(&userInfo.Assistant).Related(&userInfo.Student)
	return userInfo
}

func (user User) ChangePassword(hashedPassword string) {
	database.Db.Model(&user).Update(User{Password: hashedPassword}) // update password
	database.Db.Save(&user)
}

func (user User) Teacher() (Teacher, error) {
	var teacher Teacher
	if err := database.Db.Model(&user).Related(&teacher).Error; err != nil {
		return Teacher{}, err
	}
	return teacher, nil
}

func (user User) Student() (Student, error) {
	var student Student
	if err := database.Db.Model(&user).Preload("User").Related(&student).Error; err != nil {
		return Student{}, err
	}
	return student, nil
}
func (user User) Assistant() (Assistant, error) {
	var assistant Assistant
	if err := database.Db.Model(&user).Related(&assistant).Error; err != nil {
		return Assistant{}, err
	}
	return assistant, nil
}

//这里借助Gorm的软删除实现了一个特性:用户可以复用之前被撤回的角色
func (user User) AssignSysAdminRole() error {
	var sysAdmin SysAdmin
	if err := database.Db.Unscoped().Where(SysAdmin{UserID: user.ID}).First(&sysAdmin).Error; err == nil {
		return database.Db.Unscoped().Model(&sysAdmin).Update("deleted_at", nil).Error
	}
	return database.Db.Create(&SysAdmin{User: user}).Error
}

func (user User) AssignTeachAdminRole() error {
	var teachAdmin TeachAdmin
	if err := database.Db.Unscoped().Where(TeachAdmin{UserID: user.ID}).First(&teachAdmin).Error; err == nil {
		return database.Db.Unscoped().Model(&teachAdmin).Update("deleted_at", nil).Error
	}
	return database.Db.Create(&TeachAdmin{User: user}).Error
}

func (user User) AssignTeacherRole() error {
	var teacher Teacher
	if err := database.Db.Unscoped().Where(Teacher{UserID: user.ID}).First(&teacher).Error; err == nil {
		return database.Db.Unscoped().Model(&teacher).Update("deleted_at", nil).Error
	}
	return database.Db.Create(&Teacher{User: user}).Error
}

func (user User) AssignAssistantRole() error {
	var err error
	//只有学生才能作为助教
	if _, err = user.Student(); err != nil {
		return err
	}
	var assistant Assistant
	if err := database.Db.Unscoped().Where(Assistant{UserID: user.ID}).First(&assistant).Error; err == nil {
		return database.Db.Unscoped().Model(&assistant).Update("deleted_at", nil).Error
	}
	return database.Db.Create(&Assistant{User: user}).Error
}

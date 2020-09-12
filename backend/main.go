package main

import (
	"backend/database"
	"backend/models"
	"backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/controllers"
	_ "github.com/icattlecoder/godaemon"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

func main() {
	controllers.InitClient() // init clientset

	database.Connect() //connect database
	//migrate
	database.Db.AutoMigrate(&models.User{}, &models.Role{}, &models.Student{}, &models.Teacher{}, &models.Assistant{}, &models.TeachAdmin{}, &models.SysAdmin{}, &models.Course{}, &models.Experiment{}, &models.Notice{}, &models.StudentCourse{}, &models.File{}, &models.StudentExperiment{}, &models.UserResource{})
	defer database.Db.Close() //close database when exit
	// init database with fake data
	//utils.InitDatabase(database.Db)

	r := gin.Default()

	//init routes
	routes.InitRoutes(r)
	//
	log.Fatal(r.Run(":18080"))
}
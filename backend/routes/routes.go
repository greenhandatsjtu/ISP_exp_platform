package routes

import (
	"backend/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/controllers"
)

func InitRoutes(router *gin.Engine) {
	router.StaticFS("/file", gin.Dir("./uploads", true))
	//init casbin
	middlewares.InitEnforcer()
	r := router.Group("/api")

	//store := cookie.NewStore([]byte("secret")) // just use session and cookie
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret")) //user redis for session
	r.Use(sessions.Sessions("session_id", store))

	//do not need authority
	r.POST("/login", controllers.PostLogin)

	//need normal user authority
	r.Use(middlewares.UserMiddleware())
	r.GET("/logout", controllers.GetLogout)

	r.GET("/ping", controllers.Ping)

	r.GET("/notice", controllers.GetAllNotices)
	r.GET("/notice/:id", controllers.GetNotice)

	r.GET("/home", controllers.GetPersonalInfo)
	r.POST("/password", controllers.PostChangePassword)

	r.GET("/course", controllers.GetAllCourses)
	r.GET("/course/:id", controllers.GetCourse)
	r.GET("/course/:id/experiment", controllers.GetCourseExperiments)

	r.GET("/experiment/:eno", controllers.GetExperiment)

	student := r.Group("/student")
	student.Use(middlewares.CheckIsStudentMiddleware())

	student.GET("/experiments/grades", controllers.GetExperimentGrades)
	student.GET("/experiment/:eno/grade", controllers.GetExperimentGrade)
	student.GET("/courses", controllers.GetAssignedCourses)
	student.GET("/courses/grades", controllers.GetCoursesGrades)
	student.GET("/course/:cno/grade", controllers.GetCourseGrade)

	//need casbin authority
	admin := r.Group("/admin")
	admin.Use(middlewares.CasbinMiddleWare())

	//system admin
	admin.POST("/authority", controllers.AssignAuthority)
	admin.GET("/onlineuser", controllers.GetOnlineUsers)
	admin.GET("/user/:id", controllers.GetUser)
	admin.GET("/user/:id/role", controllers.GetUserRole)

	admin.POST("/notice", controllers.PostNotice)
	admin.GET("/notice/:id/delete", controllers.DeleteNotice)

	admin.GET("/course/:id/student", controllers.GetCourseStudents)
	admin.POST("/course/:id/experiment", controllers.PostCourseExperiment)
	admin.POST("/course", controllers.PostCourse)
	admin.POST("/assign", controllers.AssignCourse)
	admin.POST("/course/:id/update", controllers.PostUpdateCourse)
	admin.GET("course/:id/delete", controllers.DeleteCourse)

	admin.GET("/teacher", controllers.GetAllTeachers)
	admin.GET("/teacher/:id/course", controllers.GetTeacherCourses)

	admin.GET("/assistant/:id/course", controllers.GetAssistantCourses)

	admin.GET("/student", controllers.GetAllStudents)
	admin.GET("/student/:sno", controllers.GetStudent)
	admin.GET("/student/:sno/grade", controllers.GetStudentGrade)
	admin.GET("/student/:sno/experiment/:eno/grade", controllers.GetStudentExperimentGrade)
	admin.GET("/student/:sno/course", controllers.GetStudentCourses)

	admin.GET("/experiment", controllers.GetAllExperiments)
	admin.POST("/experiment", controllers.PostExperiment)
	admin.POST("/experiment/:eno/update", controllers.PostUpdateExperiment)
	admin.GET("/experiment/:eno/delete", controllers.DeleteExperiment)
	admin.GET("/experiment/:eno/student", controllers.GetExperimentStudents)
	admin.POST("/grade/course/upload", controllers.PostCourseGrade)
	admin.POST("/grade/course/update", controllers.UpdateCourseGrade)
	admin.POST("/experiments/grade/upload", controllers.PostExperimentGrade)
	admin.POST("/experiments/grade/update", controllers.PostUpdateGrade)

	admin.GET("/teachercourses", controllers.GetTeacherOwnCourses)
	admin.POST("/course/:id/assistants", controllers.AssignAssistant)
	admin.GET("/assistantcourses", controllers.GetAssistantOwnCourses)
	admin.GET("/student/:sno/course/:cno/grade", controllers.GetStudentCourseGrade)
	admin.POST("/course/:id/teacher", controllers.AssignTeacher)
	admin.POST("/course/:id/students", controllers.AssignStudents)

	admin.POST("/experiment/:eno/doc", controllers.PostUploadDoc)
	admin.DELETE("/experiment/:eno/doc/:file", controllers.DeleteDoc)

	r.GET("/experiment/:eno/start", middlewares.CheckCanExperimentMiddleware(), middlewares.CheckResourceMiddleware(), controllers.StartExperiment)
	r.GET("/experiment/:eno/end", controllers.EndExperiment)
	r.GET("/experiment/:eno/docs", controllers.GetExperimentDocs)
	r.GET("/experiment/:eno/doc/:file", controllers.GetExperimentDoc)

	student.POST("/experiment/:eno/report", controllers.PostUploadReport)

	admin.POST("/experiment/:eno/yaml", controllers.PostUploadYaml)
	admin.GET("/experiment/:eno/reports/download", controllers.DownloadReports)
	admin.GET("/experiment/:eno/enable", controllers.EnableExperiment)
	admin.GET("/experiment/:eno/disable", controllers.DisableExperiment)
	admin.GET("/experiment/:eno/student/:id/report", controllers.GetReport)

	admin.GET("/resources/all", controllers.GetAllResourceStatus)
	admin.GET("/resources/experiment/:eno", controllers.GetExperimentResourceStatus)
	admin.GET("/resources/course/:cno", controllers.GetCourseResourceStatus)

	admin.GET("/courses", controllers.GetOwnCourses)
	admin.GET("/assistants", controllers.GetAllAssistants)
	r.GET("course/:id/assistants", controllers.GetCourseAssistants)
	r.GET("course/:id/teachers", controllers.GetCourseTeachers)

	admin.GET("users", controllers.GetAllUsers)

	student.GET("/experiment/:eno/report", controllers.GetStudentOwnReport)
	student.GET("/experiment/:eno/report/download", controllers.DownloadStudentOwnReport)

	admin.POST("students", controllers.PostAddStudents)

	admin.GET("/experiment/:eno/reports", controllers.GetReports)

	r.GET("experiment/:eno/status", controllers.GetResourceStatus)

	admin.POST("user/:id/role", controllers.AddRole)
	admin.DELETE("user/:id/role/:role", controllers.RevokeRole)
	admin.DELETE("user/:id", controllers.DeleteUser)
	admin.PUT("user/:id/password", controllers.ChangeUserPassword)
	admin.POST("user", controllers.AddUser)
	admin.GET("/experiment/:eno/start", middlewares.CheckResourceMiddleware(), controllers.StartExperiment)
	admin.DELETE("/experiment/:eno/user/:id", controllers.AdminEndExperiment)
	admin.GET("/experiment/:eno/yaml/:file", controllers.GetExperimentYaml)
	admin.GET("/users/template", controllers.DownloadTemplate)
	admin.GET("/resources/metrics", controllers.GetNodeMetrics)
}

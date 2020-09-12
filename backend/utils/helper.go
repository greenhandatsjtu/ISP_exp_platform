package utils

import (
	"github.com/bxcodec/faker/v3"
	"github.com/casbin/casbin/v2"
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/models"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

//InitDatabase init database with fake data
func InitDatabase(db *gorm.DB) {

	faker.SetGenerateUniqueValues(true) // Enable unique data generation on single fake data functions
	rand.Seed(time.Now().UnixNano())
	SeedRole(db)
	SeedUser(db)
	SeedCourse(4, db)
	SeedExperiment(40, db)
	SeedNotice(10, db)
}

func SeedRole(db *gorm.DB) {
	var role models.Role

	role = models.Role{
		Description: "teach_admin",
	}
	db.Create(&role)

	role = models.Role{
		Description: "sys_admin",
	}
	db.Create(&role)

	role = models.Role{
		Description: "teacher",
	}
	db.Create(&role)

	role = models.Role{
		Description: "assistant",
	}
	db.Create(&role)

	role = models.Role{
		Description: "student",
	}
	db.Create(&role)
}

func SeedUser(db *gorm.DB) {
	//make password
	hash, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.MinCost)

	var role []models.Role

	//seed system admin
	db.Where(models.Role{Description: "sys_admin"}).Find(&role)
	var sysAdmin models.SysAdmin
	for i := 1; i <= 2; i++ {
		sysAdmin = models.SysAdmin{
			User: models.User{
				Email:    faker.Email(),
				Password: string(hash),
				Name:     faker.FirstName() + " " + faker.LastName(),
				Roles:    role,
			},
		}
		db.Create(&sysAdmin)
	}

	//seed teach admin
	var teachAdmin models.TeachAdmin
	db.Where(models.Role{Description: "teach_admin"}).Find(&role)
	for i := 3; i <= 5; i++ {
		teachAdmin = models.TeachAdmin{
			User: models.User{
				Email:    faker.Email(),
				Password: string(hash),
				Name:     faker.FirstName() + " " + faker.LastName(),
				Roles:    role,
			},
		}
		db.Create(&teachAdmin)
	}

	//seed teacher
	var teacher models.Teacher
	for i := 6; i <= 15; i++ {
		db.Where(models.Role{Description: "teacher"}).Find(&role)
		teacher = models.Teacher{
			User: models.User{
				Email:    faker.Email(),
				Password: string(hash),
				Name:     faker.FirstName() + " " + faker.LastName(),
				Roles:    role,
			},
		}
		db.Create(&teacher)
	}

	//seed student
	var student models.Student
	for i := 16; i <= 120; i++ {
		db.Where(models.Role{Description: "student"}).Find(&role)
		student = models.Student{
			User: models.User{
				Email:    faker.Email(),
				Password: string(hash),
				Name:     faker.FirstName() + " " + faker.LastName(),
				Roles:    role,
			},
			SNo: strconv.Itoa(517021910400 + i),
			CNo: "F" + strconv.Itoa(SeedNumber(1702601, 10)),
		}
		db.Create(&student)
	}

	//seed assistant
	var assistant models.Assistant
	for i := 16; i <= 20; i++ {
		db.Where(models.Role{Description: "assistant"}).Find(&role)
		var user models.User
		db.Where(i).First(&user)
		db.Model(&user).Association("Roles").Append(role)
		assistant = models.Assistant{
			User: user,
		}
		db.Create(&assistant)
	}
}

func SeedCourse(num int, db *gorm.DB) {
	var course models.Course
	var teachers []models.Teacher
	var assistants []models.Assistant
	var students []models.Student
	for i := 1; i <= num; i++ {
		var tno, sno, ano []uint
		for j := 0; j < 2; j++ {
			tno = append(tno, uint(SeedNumber(1, 10)))
		}
		if db.Where(tno).Find(&teachers); len(teachers) == 0 {
			continue
		}
		for j := 0; j < 2; j++ {
			ano = append(ano, uint(SeedNumber(1, 5)))
		}
		if db.Where(ano).Find(&assistants); len(assistants) == 0 {
			continue
		}
		for j := 0; j < 20; j++ {
			sno = append(sno, uint(SeedNumber(1, 100)))
		}
		if db.Where(sno).Find(&students); len(students) == 0 {
			continue
		}
		course = models.Course{
			CName: faker.Sentence(),
		}
		db.Create(&course)
		db.Model(&course).Association("Students").Append(students)
		db.Model(&course).Association("Teachers").Append(teachers)
		db.Model(&course).Association("Assistants").Append(assistants)
	}
}

func SeedExperiment(num int, db *gorm.DB) {
	var experiment models.Experiment
	for i := 1; i <= num; i++ {
		etime, _ := now.Parse(faker.Timestamp())
		experiment = models.Experiment{
			CourseID:   uint(SeedNumber(1, 4)),
			EName:      faker.Sentence(),
			ETime:      etime,
			Assignment: faker.Sentence(),
		}
		db.Create(&experiment)
	}
}

func SeedNotice(num int, db *gorm.DB) {
	var notice models.Notice
	for i := 1; i <= num; i++ {
		notice = models.Notice{
			Title:  faker.Sentence(),
			Body:   faker.Paragraph(),
			Author: faker.FirstName() + " " + faker.LastName(),
		}
		db.Create(&notice)
	}
}

func InitAuthority(e *casbin.Enforcer) {
	e.AddPolicy("(teach_admin)|(teacher)|(assistant)", "/student/:sno/course", "GET")
	e.AddPolicy("(teach_admin)|(teacher)|(assistant)", "/student/:sno", "GET")
	e.AddPolicy("(teach_admin)|(teacher)|(assistant)", "/notice", "POST")
	e.AddPolicy("(teach_admin)|(teacher)|(assistant)", "/notice/:id/delete", "GET")
	e.AddPolicy("(teach_admin)|(teacher)|(assistant)", "/course/:id/student", "GET")
	e.AddPolicy("(teach_admin)|(teacher)|(assistant)", "/course/:id/update", "POST")
	e.AddPolicy("(teach_admin)|(teacher)|(assistant)", "/course/:id/students", "POST")

	e.AddPolicy("(teach_admin)|(teacher)|(assistant)", "/student", "GET")
	e.AddPolicy("(teach_admin)|(teacher)", "/teacher/:id/course", "GET")

	e.AddPolicy("(teacher)|(assistant)", "/experiment", "(POST)|(GET)")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/update", "POST")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/upload", "POST")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/delete", "GET")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/student", "GET")
	e.AddPolicy("(teacher)|(assistant)", "/grade/course/upload", "POST")
	e.AddPolicy("(teacher)|(assistant)", "/grade/course/update", "POST")
	e.AddPolicy("(teacher)|(assistant)", "/course/:id/experiment", "POST")
	e.AddPolicy("(teacher)|(assistant)", "/student/:sno/grade", "GET")
	e.AddPolicy("(teacher)|(assistant)", "/student/:sno/experiment/:eno/grade", "GET")
	e.AddPolicy("(teacher)|(assistant)", "/student/:sno/course/:cno/grade", "GET")
	e.AddPolicy("(teacher)|(assistant)", "/experiments/grade/upload", "POST")
	e.AddPolicy("(teacher)|(assistant)", "/experiments/grade/update", "POST")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/doc", "POST")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/yaml", "POST")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/reports", "GET")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/reports/download", "GET")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/student/:id/report", "GET")
	e.AddPolicy("(teacher)|(assistant)", "/resources/all", "GET")
	e.AddPolicy("(teacher)|(assistant)", "/resources/experiment/:eno", "GET")
	e.AddPolicy("(teacher)|(assistant)", "/resources/course/:cno", "GET")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/enable", "GET")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/disable", "GET")

	e.AddPolicy("teach_admin", "/course", "POST")
	e.AddPolicy("teach_admin", "/assign", "POST")
	e.AddPolicy("teach_admin", "/teacher", "GET")
	e.AddPolicy("teach_admin", "/course/:id/teacher", "POST")
	e.AddPolicy("teach_admin", "/course/:id/delete", "GET")
	e.AddPolicy("teach_admin", "/assistant/:id/course", "GET")

	e.AddPolicy("teacher", "/course/:id/assistants", "POST")
	e.AddPolicy("teacher", "/teachercourses", "GET")

	e.AddPolicy("assistant", "/assistantcourses", "GET")

	e.AddPolicy("(teacher)|(assistant)", "/courses", "GET")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/doc/:file", "DELETE")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/start", "GET")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/user/:id", "DELETE")
	e.AddPolicy("(teacher)|(assistant)", "/experiment/:eno/yaml/:file", "GET")
	e.AddPolicy("teacher", "/assistants", "GET")

	e.SavePolicy()
}

//SeedNumber returns a random number, first param is minimum, second param is span
func SeedNumber(min, span int) int {
	i := rand.Intn(span) + min
	return i
}

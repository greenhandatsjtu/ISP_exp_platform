package controllers

import (
	"archive/zip"
	"backend/database"
	"backend/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//GetAllExperiments godoc
// @Tags Experiment
// @Summary 获取所有实验
//@Description 以json格式返回所有实验信息
// @Produce  json
// @Success 200 {object} Response
// @Router /admin/experiment [get]
func GetAllExperiments(c *gin.Context) {
	var experiments []models.Experiment
	if err = database.Db.Find(&experiments).Error; err != nil {
		NotFound(c, "Experiment not found.")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    experiments,
	})
}

//GetExperiment godoc
// @Tags Experiment
// @Summary 获取单个实验信息
// @Description 根据实验ID获取单个实验信息
// @Produce json
// @Param id path int true "实验ID" default(1)
// @Success 200 {object} Response
// @Router /experiment/{id} [get]
func GetExperiment(c *gin.Context) {
	id := c.Params.ByName("eno")
	var experiment models.Experiment
	if err = database.Db.Where(id).Preload("Files").First(&experiment).Error; err != nil {
		NotFound(c, "Experiment not found.")
	} else {
		c.JSON(http.StatusOK, Response{
			Success: true,
			Code:    http.StatusOK,
			Message: "",
			Data:    experiment,
		})
	}
}

//GetExperimentStudents godoc
// @Tags Experiment
// @Summary 获取实验学生
// @Description 根据实验ID获取该实验所有学生，以json格式返回
// @Produce json
// @Param id path int true "实验ID" default(1)
// @Success 200 {object} Response
// @Router /admin/experiment/{id}/student [get]
func GetExperimentStudents(c *gin.Context) {
	id := c.Params.ByName("eno")
	var experiment models.Experiment
	if err = database.Db.Where(id).First(&experiment).Error; err != nil {
		NotFound(c, "Experiment not found.")
	} else {
		c.JSON(http.StatusOK, Response{
			Success: true,
			Code:    http.StatusOK,
			Message: "",
			Data:    experiment.GetCourse().GetStudents(),
		})
	}
}

// PostExperiment godoc
// @Tags Experiment
// @Summary 上传实验
// @Description 上传实验，需要的参数有课程ID,实验名，实验时间，实验内容
// @Accept json
// @Produce  json
// @Param experiment body models.Experiment true "实验" default({"course_id":1,"name":"test","time":2020-3-04T08:00:00+08:00","assignment":"test"})
// @Success 201 {object} Response
// @Router /admin/experiment [post]
func PostExperiment(c *gin.Context) {
	var experiment models.Experiment
	if err = c.Bind(&experiment); err != nil {
		badRequest(c)
	} else {
		database.Db.Create(&experiment)
		c.JSON(http.StatusCreated, Response{
			Success: true,
			Code:    http.StatusCreated,
			Message: "Experiment created.",
			Data:    experiment,
		})
	}
}

//下载某实验所有提交的实验报告
func DownloadReports(c *gin.Context) {
	id := c.Params.ByName("eno")
	var studentExperiments []models.StudentExperiment
	if err = database.Db.Where(map[string]interface{}{"experiment_id": id}).Find(&studentExperiments).Error; err != nil {
		log.Println(err)
		NotFound(c, "Reports not found")
		return
	}
	if len(studentExperiments) == 0 {
		NotFound(c, "Reports not found")
		return
	}
	c.Writer.Header().Set("Content-type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=experiment_%s.zip", id))
	ar := zip.NewWriter(c.Writer) // write zip stream to response writer
	for _, report := range studentExperiments {
		src, _ := os.Open(filepath.Join("uploads", id, "reports", report.FileName))
		dst, _ := ar.Create(report.FileName)
		_, _ = io.Copy(dst, src)
		_ = src.Close()
	}
	_ = ar.Close()
}

//获取学生提交实验报告情况
func GetReports(c *gin.Context) {
	id := c.Params.ByName("eno")
	var reports []models.StudentExperiment
	if err = database.Db.Where(map[string]interface{}{"experiment_id": id}).Preload("Student.User").Find(&reports).Error; err != nil {
		log.Println(err)
		NotFound(c, "Reports not found")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    reports,
	})
}

//获取某同学的实验报告
func GetReport(c *gin.Context) {
	eno := c.Params.ByName("eno")
	sid := c.Params.ByName("id")
	var report models.StudentExperiment
	if err = database.Db.Where(map[string]interface{}{"experiment_id": eno, "student_id": sid}).First(&report).Error; err != nil {
		log.Println(err)
		NotFound(c, "Report not found")
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", report.FileName))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(filepath.Join("uploads", eno, "reports", report.FileName))
}

//上传实验报告
func PostUploadReport(c *gin.Context) {
	user, _ := c.Get("user")
	student, _ := user.(models.User).Student()
	var uploadFile *multipart.FileHeader
	if uploadFile, err = c.FormFile("file"); err != nil {
		log.Println(err)
		badRequest(c)
		return
	}
	id := c.Params.ByName("eno")
	var temp int
	if temp, err = strconv.Atoi(id); err != nil {
		log.Println(err)
		badRequest(c)
		return
	}
	filename := student.SNo + "_" + student.User.Name + filepath.Ext(uploadFile.Filename)
	var studentExperiment models.StudentExperiment
	if err = database.Db.Where("student_id=? AND experiment_id=?", student.ID, id).First(&studentExperiment).Error; err != nil {
		studentExperiment = models.StudentExperiment{
			StudentID:    student.ID,
			ExperimentID: uint(temp),
			UploadAt:     time.Now(),
			FileName:     filename,
		}
		database.Db.Create(&studentExperiment)
	} else {
		_ = os.Remove(filepath.Join("uploads", id, "reports", studentExperiment.FileName))
		database.Db.Model(&studentExperiment).Update(&models.StudentExperiment{UploadAt: time.Now(), FileName: filename})
	}
	mkdirIfNotExists(filepath.Join("uploads", id, "reports"))
	if err = c.SaveUploadedFile(uploadFile, filepath.Join("uploads", id, "reports", filename)); err != nil {
		log.Panic(err)
		return
	}
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "File uploaded.",
		Data:    studentExperiment,
	})
}

//获取学生自己上传的报告
func GetStudentOwnReport(c *gin.Context) {
	user, _ := c.Get("user")
	student, _ := user.(models.User).Student()
	id := c.Params.ByName("eno")
	var report models.StudentExperiment
	if err = database.Db.Where("student_id=? AND experiment_id=?", student.ID, id).First(&report).Error; err != nil {
		NotFound(c, "Report not found.")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    report,
	})
}

//学生下载自己上传的报告
func DownloadStudentOwnReport(c *gin.Context) {
	user, _ := c.Get("user")
	student, _ := user.(models.User).Student()
	id := c.Params.ByName("eno")
	var report models.StudentExperiment
	if err = database.Db.Where("student_id=? AND experiment_id=?", student.ID, id).First(&report).Error; err != nil {
		NotFound(c, "Report not found.")
		return
	}
	c.File(filepath.Join("uploads", id, "reports", report.FileName))
}

//上传实验指导书、预习报告、实验报告模板
func PostUploadDoc(c *gin.Context) {
	id := c.Params.ByName("eno")
	var experiment models.Experiment
	uploadFile, err := c.FormFile("file")
	if err != nil {
		badRequest(c)
		return
	}
	if err = database.Db.Where(id).First(&experiment).Error; err != nil {
		NotFound(c, "Experiment not found.")
		return
	}
	mkdirIfNotExists(filepath.Join("uploads", id, "docs"))
	if err = c.SaveUploadedFile(uploadFile, filepath.Join("uploads", id, "docs", uploadFile.Filename)); err != nil {
		log.Panic(err)
		return
	}

	var file models.File
	if err = database.Db.Where(map[string]interface{}{"experiment_id": id, "file_name": uploadFile.Filename}).First(&file).Error; err == nil {
		//update
		database.Db.Save(&file)
	} else {
		file = models.File{FileName: uploadFile.Filename}
		if err = experiment.AddFile(&file); err != nil {
			log.Println(err)
			badRequest(c)
			return
		}
	}
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "File uploaded.",
		Data:    file,
	})
}

//删除实验指导书、预习报告、实验报告模板
func DeleteDoc(c *gin.Context) {
	id := c.Params.ByName("eno")
	name := c.Params.ByName("file")
	var file models.File
	if err = database.Db.Where(map[string]interface{}{"experiment_id": id, "file_name": name}).First(&file).Error; err != nil {
		NotFound(c, "File not found.")
		return
	}
	if err = os.Remove(filepath.Join("uploads", id, "docs", file.FileName)); err != nil {
		log.Println(err)
		NotFound(c, "File not found.")
		return
	}
	database.Db.Delete(&file)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "File deleted.",
		Data:    file,
	})
}

//获取实验指导书
func GetExperimentDocs(c *gin.Context) {
	id := c.Params.ByName("eno")
	var (
		experiment models.Experiment
		files      []models.File
	)
	if err = database.Db.Where(id).First(&experiment).Error; err != nil {
		NotFound(c, "Experiment not found.")
		return
	}
	if err = database.Db.Model(&experiment).Related(&files).Error; err != nil {
		NotFound(c, "Files not found.")
		return
	}
	if len(files) == 0 {
		NotFound(c, "Files not found.")
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=experiment_%s.zip", id))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	ar := zip.NewWriter(c.Writer) // write zip stream to response writer
	for _, file := range files {
		src, _ := os.Open(filepath.Join("uploads", id, "docs", file.FileName))
		dst, _ := ar.Create(file.FileName)
		_, _ = io.Copy(dst, src)
		_ = src.Close()
	}
	_ = ar.Close()
}

func GetExperimentDoc(c *gin.Context) {
	id := c.Params.ByName("eno")
	name := c.Params.ByName("file")
	var file models.File
	if err = database.Db.Where(map[string]interface{}{"experiment_id": id, "file_name": name}).First(&file).Error; err != nil {
		NotFound(c, "File not found.")
		return
	}
	c.File(filepath.Join("uploads", id, "docs", file.FileName))
}

func GetExperimentYaml(c *gin.Context) {
	id := c.Params.ByName("eno")
	name := c.Params.ByName("file")
	file := filepath.Join("uploads", id, "yaml", name)
	if _, err = os.Stat(file); err != nil {
		NotFound(c, "Yaml not found.")
		return
	}
	c.File(file)
}

//启用实验
func EnableExperiment(c *gin.Context) {
	id := c.Params.ByName("eno")
	var experiment models.Experiment
	if err = database.Db.Where(id).First(&experiment).Error; err != nil {
		NotFound(c, "Experiment not found.")
		return
	}
	database.Db.Model(&experiment).Update("enable", true)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: fmt.Sprintf("Experiment %s enabled.", id),
		Data:    experiment,
	})
}

//禁用实验
func DisableExperiment(c *gin.Context) {
	id := c.Params.ByName("eno")
	var experiment models.Experiment
	if err = database.Db.Where(id).First(&experiment).Error; err != nil {
		NotFound(c, "Experiment not found.")
		return
	}
	database.Db.Model(&experiment).Update("enable", false)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: fmt.Sprintf("Experiment %s disabled.", id),
		Data:    experiment,
	})
}

// PostUpdateExperiment godoc
// @Tags Experiment
// @Summary 更新实验
// @Description 根据实验ID更新该实验
// @Accept json
// @Produce  json
// @Param id path int true "实验ID" default(1)
// @Param experimentInfo body ExperimentInfo true "更新实验信息"
// @Success 201 {object} Response
// @Router /admin/experiment/{id}/update [post]
func PostUpdateExperiment(c *gin.Context) {
	var experiment models.Experiment
	id := c.Params.ByName("eno")

	if err = database.Db.Where(id).First(&experiment).Error; err != nil {
		NotFound(c, "Experiment not found.")
		return
	}
	var experimentInfo ExperimentInfo
	if err = c.Bind(&experimentInfo); err != nil {
		badRequest(c)
		return
	}
	//etime, _ := now.Parse(experimentInfo.Time)
	database.Db.Model(&experiment).Update(models.Experiment{
		EName:      experimentInfo.Name,
		ETime:      experimentInfo.Time,
		Assignment: experimentInfo.Assignment,
	})
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Experiment updated.",
		Data:    experiment,
	})
}

//DeleteExperiment godoc
// @Tags Experiment
// @Summary 删除实验
// @Produce json
// @Param id path int true "实验ID" default(1)
// @Success 200 {object} Response
// @Router /admin/experiment/{id}/delete [get]
func DeleteExperiment(c *gin.Context) {
	id := c.Params.ByName("eno")
	var experiment models.Experiment
	if err = database.Db.Where(id).First(&experiment).Error; err != nil {
		NotFound(c, "Experiment not found")
		return
	}
	database.Db.Delete(&experiment)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "Experiment " + id + " deleted.",
		Data:    experiment,
	})
}

// PostCourseGrade godoc
// @Tags Experiment
// @Summary 上传实验成绩
// @Description 上传实验成绩，需要的参数有学号、实验ID、成绩
// @Accept json
// @Produce  json
// @Param grade body ExperimentGradeInfo true "成绩信息"
// @Success 201 {object} Response
// @Router /admin/experiments/grade/upload [post]
func PostExperimentGrade(c *gin.Context) {
	var gradeInfo ExperimentGradeInfo
	if err = c.Bind(&gradeInfo); err != nil {
		badRequest(c)
		return
	}
	var student models.Student
	if err = database.Db.Where(models.Student{SNo: gradeInfo.SNo}).First(&student).Error; err != nil {
		NotFound(c, "Student not found.")
		return
	}
	var experiment models.Experiment
	if err = database.Db.Where(gradeInfo.ENo).First(&experiment).Error; err != nil {
		NotFound(c, "Experiment not found.")
		return
	}
	grade := models.StudentExperiment{
		Student:    student,
		Experiment: experiment,
		Grade:      &gradeInfo.Grade,
	}
	database.Db.Create(&grade)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Grade uploaded.",
		Data:    grade,
	})
}

// UpdateCourseGrade godoc
// @Tags Experiment
// @Summary 更新实验成绩
// @Description 更新实验成绩，需要的参数有学号、实验ID、成绩
// @Accept json
// @Produce  json
// @Param grade body ExperimentGradeInfo true "成绩信息"
// @Success 201 {object} Response
// @Router /admin/experiments/grade/update [post]
func PostUpdateGrade(c *gin.Context) {
	var gradeInfo ExperimentGradeInfo
	if err = c.Bind(&gradeInfo); err != nil {
		badRequest(c)
		return
	}
	var student models.Student
	if err = database.Db.Where(models.Student{SNo: gradeInfo.SNo}).First(&student).Error; err != nil {
		NotFound(c, "Student not found.")
		return
	}
	var experiment models.Experiment
	if err = database.Db.Where(gradeInfo.ENo).First(&experiment).Error; err != nil {
		NotFound(c, "Experiment not found.")
		return
	}
	var grade models.StudentExperiment
	if err = database.Db.Where("student_id=? AND experiment_id=?", student.ID, experiment.ID).First(&grade).Error; err != nil {
		NotFound(c, "Grade not found.")
		return
	}
	grade.UpdateGrade(gradeInfo.Grade)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Grade updated.",
		Data:    grade,
	})
}

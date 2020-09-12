package controllers

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var err error

type Auth struct {
	User  uint     `form:"user" json:"user" binding:"required"`
	Roles []string `form:"roles" json:"roles" binding:"required"`
}

type Response struct {
	Success bool        `json:"success" example:"true"`
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"example message"`
	Data    interface{} `json:"data"`
}

type ExperimentInfo struct {
	Name       string    `json:"name" form:"name" binding:"required" example:"This is example experiment name"`
	Time       time.Time `json:"time" form:"time" example:"2007-09-30 08:00:00"`
	Assignment string    `form:"assignment" json:"assignment" example:"This is example assignment."`
}

type ExperimentGradeInfo struct {
	SNo   string `json:"student" form:"student" binding:"required" example:"517021910444"`
	ENo   int    `json:"experiment" form:"experiment" binding:"required" example:"1"`
	Grade int    `json:"grade" form:"grade" example:"100"`
}

type CourseGradeInfo struct {
	SNo   string `json:"student" form:"student" binding:"required" example:"517021910444"`
	CNo   int    `json:"course" form:"course" binding:"required" example:"1"`
	Grade int    `json:"grade" form:"grade" example:"100"`
}

type NewStudentStruct struct {
	SNo   string `json:"student_number" binding:"required"`
	CNo   string `json:"class_number" binding:"required"`
	Email string `json:"email" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

func NotFound(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusNotFound, Response{
		Success: false,
		Code:    http.StatusNotFound,
		Message: msg,
		Data:    nil,
	})
}

func Forbidden(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, Response{
		Success: false,
		Code:    http.StatusForbidden,
		Message: "You don't have permission.",
		Data:    nil,
	})
}

func badRequest(c *gin.Context) {
	log.Println("Bad request")
	c.AbortWithStatusJSON(http.StatusBadRequest, Response{
		Success: false,
		Code:    http.StatusBadRequest,
		Message: "Bad request.",
		Data:    nil,
	})
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

// 如果文件夹不存在就创建文件夹
func mkdirIfNotExists(path string) {
	_, err := os.Stat(path)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}
}

package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var Db *gorm.DB
var err error

func Connect() {
	Db, err = gorm.Open("mysql", "test:test@(localhost)/k8s?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
	}
}

package models

import (
	"Music/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Init() {
	var DBNAME = config.Config.Database.Name
	var USER = config.Config.Database.User
	var PASSWORD = config.Config.Database.Password
	dsn := USER + ":" + PASSWORD + "@tcp(127.0.0.1:3306)/" + DBNAME + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Gorm init error: " + err.Error())
	}
	// Create DB
	//err = DB.Exec("CREATE DATABASE IF NOT EXISTS " + DBNAME + " CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci").Error
	//if err != nil {
	//	panic("failed to create database")
	//}
	//fmt.Println("Database created or already exists")
	// Table auto migrate
	err = DB.AutoMigrate(MusicInfo{})
	if err != nil {
		log.Fatal("迁移失败: ", err)
	}
}

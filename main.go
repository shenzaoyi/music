package main

import (
	"Music/config"
	"Music/models"
	"Music/my_utils"
	"Music/router"
	"Music/services"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志
	logFile, err := my_utils.SetupLogFile("app.log")
	if err != nil {
		my_utils.Fatal("日志设置失败:", err)
	}
	defer logFile.Close()

	// 设置日志级别
	my_utils.SetLogLevel(my_utils.LevelInfo)

	// Init Config
	config.InitConfig()

	// Init Database
	models.Init()

	// test
	services := services.MusicService{}
	musicinfo := models.MusicInfo{
		Singer: "jay chou",
		Album:  "fantasy",
		Name:   "yequ",
		Cover:  "www.cover.com",
	}
	path := "src/1.mp3"
	err = services.CreateMusic(&musicinfo, path)
	if err != nil {
		fmt.Println(err.Error())
	}
	// test

	// 初始化路由
	r := gin.Default()
	router.InitRouter(r)

	// 启动服务
	my_utils.Info("启动服务，监听端口：8080")
	if err := r.Run(":8080"); err != nil {
		my_utils.Fatal("服务启动失败:", err)
	}
}

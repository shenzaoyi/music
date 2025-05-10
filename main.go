package main

import (
	"Music/config"
	"Music/models"
	"Music/services"
	"Music/utils"
	"fmt"
)

func main() {
	// 初始化日志
	logFile, err := utils.SetupLogFile("app.log")
	if err != nil {
		utils.Fatal("日志设置失败:", err)
	}
	defer logFile.Close()

	// 设置日志级别
	utils.SetLogLevel(utils.LevelInfo)

	// Init Database
	models.Init()

	// Init Config
	config.InitConfig()

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
	//r := gin.Default()
	//router.InitRouter(r)
	//
	//// 启动服务
	//utils.Info("启动服务，监听端口：8080")
	//if err := r.Run(":8080"); err != nil {
	//	utils.Fatal("服务启动失败:", err)
	//}
}

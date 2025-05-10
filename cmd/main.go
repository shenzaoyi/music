package main

import (
	"Music/config"
	"Music/models"
	"Music/my_utils"
	"Music/services"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Prepare() {
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
}

func main() {
	Prepare()
	rootDir := "./src/"
	singer := "周杰伦"
	service := services.NewMusicService()

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// 只处理 .flac 和 .mp3 文件
		//if !strings.HasSuffix(info.Name(), ".flac") && !strings.HasSuffix(info.Name(), ".mp3") {
		//	return nil
		//}

		ext := filepath.Ext(info.Name())
		baseName := strings.TrimSuffix(info.Name(), ext)

		// 专辑名 = 父目录名
		album := filepath.Base(filepath.Dir(path))

		// 去掉 "周杰伦 - " 前缀
		title := strings.TrimPrefix(baseName, singer+" - ")

		music := &models.MusicInfo{
			Name:   title,
			Album:  album,
			Singer: singer,
		}

		fmt.Printf("开始上传 [%s] - [%s] - [%s]\n", album, title, path)
		if err := service.CreateMusic(music, path); err != nil {
			log.Printf("❌ 上传失败 [%s]: %v\n", info.Name(), err)
		} else {
			log.Printf("✅ 上传成功 [%s]\n", info.Name())
		}

		return nil
	})

	if err != nil {
		log.Fatalf("遍历目录出错: %v", err)
	}
}

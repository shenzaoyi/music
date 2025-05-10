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
	"sync"
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

	var wg sync.WaitGroup
	sem := make(chan struct{}, 3) // 控制最大并发数为 5，可调整

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := filepath.Ext(info.Name())
		if ext != ".flac" && ext != ".mp3" {
			return nil
		}

		baseName := strings.TrimSuffix(info.Name(), ext)
		album := filepath.Base(filepath.Dir(path))
		title := strings.TrimPrefix(baseName, singer+" - ")

		music := &models.MusicInfo{
			Name:   title,
			Album:  album,
			Singer: singer,
		}

		wg.Add(1)
		sem <- struct{}{} // 占用一个槽位

		go func(m *models.MusicInfo, p string, name string, album string) {
			defer wg.Done()
			defer func() { <-sem }() // 释放槽位

			fmt.Printf("开始上传 [%s] - [%s] - [%s]\n", album, m.Name, p)
			if err := service.CreateMusic(m, p); err != nil {
				log.Printf("❌ 上传失败 [%s]: %v\n", name, err)
			} else {
				log.Printf("✅ 上传成功 [%s]\n", name)
			}
		}(music, path, info.Name(), album)

		return nil
	})

	wg.Wait()

	if err != nil {
		log.Fatalf("遍历目录出错: %v", err)
	}
}

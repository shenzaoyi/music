package my_utils

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var (
	logLevel   = LevelInfo
	logFile    *os.File
	logMutex   sync.Mutex
	maxSize    int64 = 10 * 1024 * 1024 // 10MB
	currentLog string
)

// 设置日志级别
func SetLogLevel(level int) {
	logMutex.Lock()
	defer logMutex.Unlock()
	logLevel = level
}

// 检查并轮转日志文件
func rotateLogFile(filename string) error {
	logMutex.Lock()
	defer logMutex.Unlock()

	// 检查文件权限
	if err := checkFilePermissions(filename); err != nil {
		return fmt.Errorf("文件权限检查失败: %v", err)
	}

	// 获取文件信息
	info, err := logFile.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 如果文件大小超过限制
	if info.Size() >= maxSize {
		// 创建备份文件
		backupName := fmt.Sprintf("%s.%s", currentLog, time.Now().Format("20060102-150405"))
		backupFile, err := os.Create(backupName)
		if err != nil {
			return fmt.Errorf("创建备份文件失败: %v", err)
		}
		defer backupFile.Close()

		// 复制当前日志内容到备份文件
		if _, err := logFile.Seek(0, 0); err != nil {
			return fmt.Errorf("文件指针重置失败: %v", err)
		}
		if _, err := io.Copy(backupFile, logFile); err != nil {
			return fmt.Errorf("日志复制失败: %v", err)
		}

		// 关闭当前日志文件
		if err := logFile.Close(); err != nil {
			return fmt.Errorf("关闭日志文件失败: %v", err)
		}

		// 创建新的日志文件
		newFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("创建新日志文件失败: %v", err)
		}

		// 更新文件句柄
		logFile = newFile

		// 设置新的日志输出
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}

	return nil
}

// 检查文件权限
func checkFilePermissions(filename string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// 如果文件不存在，检查目录权限
		dir := filepath.Dir(filename)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %v", err)
		}
		return nil
	}

	// 检查文件权限
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("文件权限不足: %v", err)
	}
	file.Close()

	return nil
}

// 初始化日志轮转
func initLogRotate(filename string) {
	c := cron.New()

	// 添加定时任务：每周日零点轮转日志文件
	_, err := c.AddFunc("0 0 * * 0", func() {
		err := rotateLogFile(filename)
		if err != nil {
			log.Printf("日志轮转失败: %v", err)
		}
	})

	if err != nil {
		log.Printf("添加日志轮转定时任务失败: %v", err)
		return
	}

	c.Start()
}

// 配置日志文件
func SetupLogFile(filename string) (*os.File, error) {
	// 获取当前程序运行的目录
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("获取当前目录失败: %v", err)
	}

	// 构建完整的日志文件路径
	logPath := filepath.Join(currentDir, filename)
	currentLog = logPath

	// 打开日志文件
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %v", err)
	}

	// 设置日志输出
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 初始化日志轮转
	initLogRotate(logPath)

	return logFile, nil
}

// 自定义日志函数
func Debug(format string, v ...interface{}) {
	if logLevel <= LevelDebug {
		msg := fmt.Sprintf("[DEBUG] "+format, v...)
		log.Println(msg)
	}
}

func Info(format string, v ...interface{}) {
	if logLevel <= LevelInfo {
		msg := fmt.Sprintf("[INFO] "+format, v...)
		log.Println(msg)
	}
}

func Warn(format string, v ...interface{}) {
	if logLevel <= LevelWarn {
		msg := fmt.Sprintf("[WARN] "+format, v...)
		log.Println(msg)
	}
}

func Error(format string, v ...interface{}) {
	if logLevel <= LevelError {
		msg := fmt.Sprintf("[ERROR] "+format, v...)
		log.Println(msg)
	}
}

func Fatal(format string, v ...interface{}) {
	if logLevel <= LevelFatal {
		msg := fmt.Sprintf("[FATAL] "+format, v...)
		log.Fatal(msg)
	}
}

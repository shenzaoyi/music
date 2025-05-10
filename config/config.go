package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type TencentCOSConfig struct {
	Region    string `yaml:"region"`
	Bucket    string `yaml:"bucket"`
	SecretID  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
}

type DatabaseConfig struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

type AppConfig struct {
	Database   DatabaseConfig   `yaml:"database"`
	TencentCOS TencentCOSConfig `yaml:"tencent_cos"`
}

var Config AppConfig

func InitConfig() {
	file, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
}

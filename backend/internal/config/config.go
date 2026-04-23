package config

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"os"
)

// Config 全局配置结构体
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	DSN      string `yaml:"dsn"`
	MaxIdle  int    `yaml:"max_idle"`
	MaxOpen  int    `yaml:"max_open"`
	LifeTime int    `yaml:"life_time"` // 单位：小时
}

// ServerConfig 服务器配置结构体
type ServerConfig struct {
	Port    string `yaml:"port"`
	BaseURL string `yaml:"base_url"`
}

var GlobalConfig *Config

// LoadConfig 从指定文件加载配置
func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	GlobalConfig = &cfg
	return nil
}

// InitDefaultConfig 初始化默认配置（作为兜底）
func InitDefaultConfig() {
	GlobalConfig = &Config{
		Database: DatabaseConfig{
			Driver:   "mysql",
			DSN:      "root:231792@tcp(localhost:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local",
			MaxIdle:  10,
			MaxOpen:  100,
			LifeTime: 1,
		},
		Server: ServerConfig{
			Port:    "8080",
			BaseURL: "http://localhost:8080",
		},
	}
}

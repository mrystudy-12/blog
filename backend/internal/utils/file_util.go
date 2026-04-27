package utils

import (
	"GoWork_9/backend/internal/config"
	"fmt"
	"os"
	"path/filepath"
)

func GetSavePath(category string) (string, error) {
	wd, _ := os.Getwd()

	// 灵活判断：无论在 backend 还是根目录运行都能兼容
	var base string
	if filepath.Base(wd) == "backend" {
		base = filepath.Join(wd, "..")
	} else {
		base = wd
	}

	// 最终路径：根目录/frontend/static/images/类别
	finalPath := filepath.Join(base, "frontend", "static", "images", category)

	// 自动创建文件夹
	if err := os.MkdirAll(finalPath, 0755); err != nil {
		return "", err
	}
	return finalPath, nil
}

func GetAccessURL(category string, filename string) string {
	// 这里的 BaseURL 来自你的配置文件 http://localhost:8080
	baseURL := config.GlobalConfig.Server.BaseURL
	return fmt.Sprintf("%s/static/images/%s/%s", baseURL, category, filename)
}

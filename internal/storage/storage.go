package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"goSSH/models"
)

// Storage 提供配置文件的读写功能
type Storage struct {
	configPath string
}

// NewStorage 创建新的存储实例
func NewStorage() (*Storage, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("获取配置目录失败: %v", err)
	}

	gosshDir := filepath.Join(configDir, "gossh")
	if err := os.MkdirAll(gosshDir, 0755); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %v", err)
	}

	configPath := filepath.Join(gosshDir, "servers.json")
	return &Storage{configPath: configPath}, nil
}

// Load 加载配置文件
func (s *Storage) Load() (*models.ServerConfig, error) {
	config := &models.ServerConfig{
		Servers: make([]models.Server, 0),
	}

	// 如果文件不存在，返回空配置
	if _, err := os.Stat(s.configPath); os.IsNotExist(err) {
		return config, nil
	}

	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	if len(data) == 0 {
		return config, nil
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return config, nil
}

// Save 保存配置文件
func (s *Storage) Save(config *models.ServerConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	if err := os.WriteFile(s.configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}


package config

import (
	"fmt"

	"goSSH/internal/storage"
	"goSSH/models"
)

// Manager 管理服务器配置
type Manager struct {
	storage *storage.Storage
}

// NewManager 创建新的配置管理器
func NewManager() (*Manager, error) {
	st, err := storage.NewStorage()
	if err != nil {
		return nil, err
	}
	return &Manager{storage: st}, nil
}

// AddServer 添加服务器
func (m *Manager) AddServer(server models.Server) error {
	config, err := m.storage.Load()
	if err != nil {
		return err
	}

	// 检查是否已存在同名服务器
	for _, s := range config.Servers {
		if s.Name == server.Name {
			return fmt.Errorf("服务器 '%s' 已存在", server.Name)
		}
	}

	// 设置默认端口
	if server.Port == 0 {
		server.Port = 22
	}

	config.Servers = append(config.Servers, server)
	return m.storage.Save(config)
}

// RemoveServer 删除服务器
func (m *Manager) RemoveServer(name string) error {
	config, err := m.storage.Load()
	if err != nil {
		return err
	}

	found := false
	newServers := make([]models.Server, 0, len(config.Servers))
	for _, s := range config.Servers {
		if s.Name != name {
			newServers = append(newServers, s)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("服务器 '%s' 不存在", name)
	}

	config.Servers = newServers
	return m.storage.Save(config)
}

// ListServers 列出所有服务器
func (m *Manager) ListServers() ([]models.Server, error) {
	config, err := m.storage.Load()
	if err != nil {
		return nil, err
	}
	return config.Servers, nil
}

// GetServer 根据名称获取服务器
func (m *Manager) GetServer(name string) (*models.Server, error) {
	config, err := m.storage.Load()
	if err != nil {
		return nil, err
	}

	for _, s := range config.Servers {
		if s.Name == name {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("服务器 '%s' 不存在", name)
}

// UpdateServer 更新服务器信息
func (m *Manager) UpdateServer(server models.Server) error {
	config, err := m.storage.Load()
	if err != nil {
		return err
	}

	found := false
	for i := range config.Servers {
		if config.Servers[i].Name == server.Name {
			config.Servers[i] = server
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("服务器 '%s' 不存在", server.Name)
	}

	return m.storage.Save(config)
}


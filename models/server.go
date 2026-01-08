package models

// Server 表示一个SSH服务器配置信息
type Server struct {
	Name     string `json:"name"`     // 服务器别名/名称
	Host     string `json:"host"`     // IP地址或主机名
	Port     int    `json:"port"`     // SSH端口，默认22
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码（明文存储）
}

// ServerConfig 表示服务器配置文件结构
type ServerConfig struct {
	Servers []Server `json:"servers"` // 服务器列表
}


package ssh

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
	"goSSH/models"
)

// Client SSH客户端封装
type Client struct {
	server *models.Server
	conn   *ssh.Client
}

// NewClient 创建新的SSH客户端
func NewClient(server *models.Server) *Client {
	return &Client{
		server: server,
	}
}

// Connect 建立SSH连接
func (c *Client) Connect() error {
	config := &ssh.ClientConfig{
		User: c.server.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略主机密钥验证（适合内网环境）
		Timeout:         10 * time.Second,
	}

	address := fmt.Sprintf("%s:%d", c.server.Host, c.server.Port)
	conn, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return fmt.Errorf("连接服务器失败: %v", err)
	}

	c.conn = conn
	return nil
}

// Close 关闭SSH连接
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetConnection 获取SSH连接（用于执行命令或文件传输）
func (c *Client) GetConnection() *ssh.Client {
	return c.conn
}

// IsConnected 检查是否已连接
func (c *Client) IsConnected() bool {
	return c.conn != nil
}

// Reconnect 重连SSH服务器
func (c *Client) Reconnect() error {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	return c.Connect()
}

// TestConnection 测试连接（不保持连接）
func TestConnection(server *models.Server) error {
	config := &ssh.ClientConfig{
		User: server.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	address := fmt.Sprintf("%s:%d", server.Host, server.Port)
	conn, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return fmt.Errorf("连接测试失败: %v", err)
	}
	defer conn.Close()

	return nil
}

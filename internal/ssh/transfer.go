package ssh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
)

// Transfer 提供文件传输功能
type Transfer struct {
	client  *Client
	sftpCli *sftp.Client
}

// NewTransfer 创建新的文件传输器
func NewTransfer(client *Client) (*Transfer, error) {
	if !client.IsConnected() {
		if err := client.Connect(); err != nil {
			return nil, err
		}
	}

	sshCli := client.GetConnection()
	sftpCli, err := sftp.NewClient(sshCli)
	if err != nil {
		return nil, fmt.Errorf("创建SFTP客户端失败: %v", err)
	}

	return &Transfer{
		client:  client,
		sftpCli: sftpCli,
	}, nil
}

// Close 关闭SFTP连接
func (t *Transfer) Close() error {
	if t.sftpCli != nil {
		return t.sftpCli.Close()
	}
	return nil
}

// Upload 上传文件到远程服务器
func (t *Transfer) Upload(localPath, remotePath string) error {
	// 打开本地文件
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("打开本地文件失败: %v", err)
	}
	defer localFile.Close()

	// 获取文件信息
	localInfo, err := localFile.Stat()
	if err != nil {
		return fmt.Errorf("获取本地文件信息失败: %v", err)
	}

	// 创建远程目录（如果不存在）
	remoteDir := filepath.Dir(remotePath)
	if err := t.sftpCli.MkdirAll(remoteDir); err != nil {
		return fmt.Errorf("创建远程目录失败: %v", err)
	}

	// 创建远程文件
	remoteFile, err := t.sftpCli.Create(remotePath)
	if err != nil {
		return fmt.Errorf("创建远程文件失败: %v", err)
	}
	defer remoteFile.Close()

	// 设置文件权限
	if err := remoteFile.Chmod(localInfo.Mode().Perm()); err != nil {
		// 忽略权限设置错误（某些系统可能不支持）
	}

	// 复制文件内容
	fmt.Printf("正在上传: %s -> %s\n", localPath, remotePath)
	written, err := io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("上传文件失败: %v", err)
	}

	fmt.Printf("上传完成: %d 字节\n", written)
	return nil
}

// Download 从远程服务器下载文件
func (t *Transfer) Download(remotePath, localPath string) error {
	// 打开远程文件
	remoteFile, err := t.sftpCli.Open(remotePath)
	if err != nil {
		return fmt.Errorf("打开远程文件失败: %v", err)
	}
	defer remoteFile.Close()

	// 获取远程文件信息
	remoteInfo, err := remoteFile.Stat()
	if err != nil {
		return fmt.Errorf("获取远程文件信息失败: %v", err)
	}

	// 创建本地目录（如果不存在）
	localDir := filepath.Dir(localPath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return fmt.Errorf("创建本地目录失败: %v", err)
	}

	// 创建本地文件
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("创建本地文件失败: %v", err)
	}
	defer localFile.Close()

	// 复制文件内容
	fmt.Printf("正在下载: %s -> %s\n", remotePath, localPath)
	written, err := io.Copy(localFile, remoteFile)
	if err != nil {
		return fmt.Errorf("下载文件失败: %v", err)
	}

	// 设置文件权限
	if err := localFile.Chmod(remoteInfo.Mode().Perm()); err != nil {
		// 忽略权限设置错误
	}

	fmt.Printf("下载完成: %d 字节\n", written)
	return nil
}

// UploadDir 上传整个目录到远程服务器
func (t *Transfer) UploadDir(localDir, remoteDir string) error {
	return filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}

		remotePath := filepath.Join(remoteDir, relPath)

		if info.IsDir() {
			// 创建远程目录
			return t.sftpCli.MkdirAll(remotePath)
		}

		// 上传文件
		return t.Upload(path, remotePath)
	})
}

// DownloadDir 从远程服务器下载整个目录
func (t *Transfer) DownloadDir(remoteDir, localDir string) error {
	return t.walkRemoteDir(remoteDir, localDir, "")
}

// walkRemoteDir 递归遍历远程目录
func (t *Transfer) walkRemoteDir(remoteDir, localDir, prefix string) error {
	// 列出远程目录内容
	files, err := t.sftpCli.ReadDir(remoteDir)
	if err != nil {
		return fmt.Errorf("读取远程目录失败: %v", err)
	}

	for _, file := range files {
		remotePath := filepath.Join(remoteDir, file.Name())
		localPath := filepath.Join(localDir, file.Name())

		if file.IsDir() {
			// 创建本地目录
			if err := os.MkdirAll(localPath, file.Mode().Perm()); err != nil {
				return fmt.Errorf("创建本地目录失败: %v", err)
			}

			// 递归处理子目录
			if err := t.walkRemoteDir(remotePath, localPath, prefix+file.Name()+"/"); err != nil {
				return err
			}
		} else {
			// 下载文件
			if err := t.Download(remotePath, localPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// ListRemote 列出远程目录内容
func (t *Transfer) ListRemote(remotePath string) ([]os.FileInfo, error) {
	files, err := t.sftpCli.ReadDir(remotePath)
	if err != nil {
		return nil, fmt.Errorf("读取远程目录失败: %v", err)
	}
	return files, nil
}

// RemoveRemote 删除远程文件或目录
func (t *Transfer) RemoveRemote(remotePath string) error {
	info, err := t.sftpCli.Stat(remotePath)
	if err != nil {
		return fmt.Errorf("获取远程文件信息失败: %v", err)
	}

	if info.IsDir() {
		return t.sftpCli.RemoveDirectory(remotePath)
	}

	return t.sftpCli.Remove(remotePath)
}

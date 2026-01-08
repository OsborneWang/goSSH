package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"goSSH/internal/config"
	"goSSH/internal/ssh"
)

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "文件传输",
	Long:  "在本地和远程服务器之间传输文件",
}

var uploadCmd = &cobra.Command{
	Use:   "upload [name] [local] [remote]",
	Short: "上传文件到远程服务器",
	Long:  "上传本地文件或目录到远程服务器",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := config.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}

		var serverName, localPath, remotePath string

		if len(args) >= 3 {
			serverName = args[0]
			localPath = args[1]
			remotePath = args[2]
		} else {
			// 交互式输入
			servers, err := manager.ListServers()
			if err != nil {
				fmt.Fprintf(os.Stderr, "错误: %v\n", err)
				return
			}

			if len(servers) == 0 {
				fmt.Println("没有配置任何服务器")
				return
			}

			items := make([]string, len(servers))
			for i, s := range servers {
				items[i] = fmt.Sprintf("%s (%s:%d)", s.Name, s.Host, s.Port)
			}

			prompt := promptui.Select{
				Label: "选择服务器",
				Items: items,
			}

			index, _, err := prompt.Run()
			if err != nil {
				fmt.Printf("操作取消: %v\n", err)
				return
			}

			serverName = servers[index].Name

			prompt2 := promptui.Prompt{
				Label: "本地路径",
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("路径不能为空")
					}
					return nil
				},
			}
			localPath, err = prompt2.Run()
			if err != nil {
				fmt.Printf("输入取消: %v\n", err)
				return
			}

			prompt3 := promptui.Prompt{
				Label: "远程路径",
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("路径不能为空")
					}
					return nil
				},
			}
			remotePath, err = prompt3.Run()
			if err != nil {
				fmt.Printf("输入取消: %v\n", err)
				return
			}
		}

		// 验证本地路径
		info, err := os.Stat(localPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}

		// 获取服务器配置
		server, err := manager.GetServer(serverName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}

		// 创建SSH客户端和传输器
		client := ssh.NewClient(server)
		defer client.Close()

		transfer, err := ssh.NewTransfer(client)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}
		defer transfer.Close()

		// 上传文件或目录
		if info.IsDir() {
			if err := transfer.UploadDir(localPath, remotePath); err != nil {
				fmt.Fprintf(os.Stderr, "上传失败: %v\n", err)
				os.Exit(1)
			}
		} else {
			// 如果是文件，确保远程路径是完整路径
			if !filepath.IsAbs(remotePath) && !strings.HasPrefix(remotePath, "/") {
				// 相对路径，使用文件名
				remotePath = filepath.Join(remotePath, filepath.Base(localPath))
			}
			if err := transfer.Upload(localPath, remotePath); err != nil {
				fmt.Fprintf(os.Stderr, "上传失败: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Println("✓ 上传完成")
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download [name] [remote] [local]",
	Short: "从远程服务器下载文件",
	Long:  "从远程服务器下载文件或目录到本地",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := config.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}

		var serverName, remotePath, localPath string

		if len(args) >= 3 {
			serverName = args[0]
			remotePath = args[1]
			localPath = args[2]
		} else {
			// 交互式输入
			servers, err := manager.ListServers()
			if err != nil {
				fmt.Fprintf(os.Stderr, "错误: %v\n", err)
				return
			}

			if len(servers) == 0 {
				fmt.Println("没有配置任何服务器")
				return
			}

			items := make([]string, len(servers))
			for i, s := range servers {
				items[i] = fmt.Sprintf("%s (%s:%d)", s.Name, s.Host, s.Port)
			}

			prompt := promptui.Select{
				Label: "选择服务器",
				Items: items,
			}

			index, _, err := prompt.Run()
			if err != nil {
				fmt.Printf("操作取消: %v\n", err)
				return
			}

			serverName = servers[index].Name

			prompt2 := promptui.Prompt{
				Label: "远程路径",
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("路径不能为空")
					}
					return nil
				},
			}
			remotePath, err = prompt2.Run()
			if err != nil {
				fmt.Printf("输入取消: %v\n", err)
				return
			}

			prompt3 := promptui.Prompt{
				Label: "本地路径",
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("路径不能为空")
					}
					return nil
				},
			}
			localPath, err = prompt3.Run()
			if err != nil {
				fmt.Printf("输入取消: %v\n", err)
				return
			}
		}

		// 获取服务器配置
		server, err := manager.GetServer(serverName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}

		// 创建SSH客户端和传输器
		client := ssh.NewClient(server)
		defer client.Close()

		transfer, err := ssh.NewTransfer(client)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}
		defer transfer.Close()

		// 检查远程路径是文件还是目录
		files, err := transfer.ListRemote(remotePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}

		// 简单判断：如果只有一个文件且是文件，则下载文件；否则下载目录
		if len(files) == 1 && !files[0].IsDir() {
			// 下载文件
			if err := transfer.Download(remotePath, localPath); err != nil {
				fmt.Fprintf(os.Stderr, "下载失败: %v\n", err)
				os.Exit(1)
			}
		} else {
			// 下载目录
			if err := transfer.DownloadDir(remotePath, localPath); err != nil {
				fmt.Fprintf(os.Stderr, "下载失败: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Println("✓ 下载完成")
	},
}

func init() {
	transferCmd.AddCommand(uploadCmd)
	transferCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(transferCmd)
}

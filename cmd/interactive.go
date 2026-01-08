package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"goSSH/internal/config"
	"goSSH/internal/ssh"
	"goSSH/models"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "进入交互式菜单模式",
	Long:  "进入交互式菜单，可以方便地选择服务器和操作",
	Run: func(cmd *cobra.Command, args []string) {
		runInteractiveMode()
	},
}

func runInteractiveMode() {
	manager, err := config.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		return
	}

	banner := color.New(color.FgCyan, color.Bold)
	banner.Println(`
╔══════════════════════════════════════╗
║      GoSSH 交互式菜单模式            ║
╚══════════════════════════════════════╝
`)

	for {
		// 主菜单
		mainMenu := []string{
			"连接服务器 (SSH Shell)",
			"执行远程命令",
			"上传文件/目录",
			"下载文件/目录",
			"添加服务器",
			"删除服务器",
			"列出所有服务器",
			"退出",
		}

		prompt := promptui.Select{
			Label: "请选择操作",
			Items: mainMenu,
			Size:  10,
		}

		index, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("操作取消或错误: %v\n", err)
			return
		}

		switch index {
		case 0: // 连接服务器
			handleConnect(manager)
		case 1: // 执行命令
			handleExecute(manager)
		case 2: // 上传文件
			handleUpload(manager)
		case 3: // 下载文件
			handleDownload(manager)
		case 4: // 添加服务器
			handleAddServer(manager)
		case 5: // 删除服务器
			handleRemoveServer(manager)
		case 6: // 列出服务器
			handleListServers(manager)
		case 7: // 退出
			fmt.Println("再见!")
			return
		}

		if result != "退出" {
			fmt.Println("\n按 Enter 继续...")
			fmt.Scanln()
		}
	}
}

func selectServer(manager *config.Manager, label string) (*models.Server, error) {
	servers, err := manager.ListServers()
	if err != nil {
		return nil, err
	}

	if len(servers) == 0 {
		return nil, fmt.Errorf("没有配置任何服务器")
	}

	items := make([]string, len(servers))
	for i, s := range servers {
		items[i] = fmt.Sprintf("%s (%s:%d - %s)", s.Name, s.Host, s.Port, s.Username)
	}

	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &servers[index], nil
}

func handleConnect(manager *config.Manager) {
	server, err := selectServer(manager, "选择要连接的服务器")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	client := ssh.NewClient(server)
	defer client.Close()

	fmt.Printf("正在连接到 %s (%s:%d)...\n", server.Name, server.Host, server.Port)
	if err := client.Connect(); err != nil {
		fmt.Printf("连接失败: %v\n", err)
		return
	}

	fmt.Printf("✓ 已连接到 %s\n\n", server.Name)

	executor := ssh.NewExecutor(client)
	if err := executor.ExecuteShell(); err != nil {
		fmt.Printf("Shell错误: %v\n", err)
	}
}

func handleExecute(manager *config.Manager) {
	server, err := selectServer(manager, "选择服务器")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	prompt := promptui.Prompt{
		Label: "要执行的命令",
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("命令不能为空")
			}
			return nil
		},
	}

	command, err := prompt.Run()
	if err != nil {
		fmt.Printf("输入取消: %v\n", err)
		return
	}

	client := ssh.NewClient(server)
	defer client.Close()

	executor := ssh.NewExecutor(client)

	fmt.Printf("\n执行命令: %s\n", command)
	fmt.Println("─────────────────────────────────────")
	if err := executor.ExecuteWithStream(command); err != nil {
		fmt.Printf("\n执行失败: %v\n", err)
	}
	fmt.Println("─────────────────────────────────────")
}

func handleUpload(manager *config.Manager) {
	server, err := selectServer(manager, "选择服务器")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	prompt1 := promptui.Prompt{
		Label: "本地路径",
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("路径不能为空")
			}
			return nil
		},
	}

	localPath, err := prompt1.Run()
	if err != nil {
		fmt.Printf("输入取消: %v\n", err)
		return
	}

	prompt2 := promptui.Prompt{
		Label: "远程路径",
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("路径不能为空")
			}
			return nil
		},
	}

	remotePath, err := prompt2.Run()
	if err != nil {
		fmt.Printf("输入取消: %v\n", err)
		return
	}

	client := ssh.NewClient(server)
	defer client.Close()

	transfer, err := ssh.NewTransfer(client)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	defer transfer.Close()

	info, err := os.Stat(localPath)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	if info.IsDir() {
		if err := transfer.UploadDir(localPath, remotePath); err != nil {
			fmt.Printf("上传失败: %v\n", err)
		}
	} else {
		if err := transfer.Upload(localPath, remotePath); err != nil {
			fmt.Printf("上传失败: %v\n", err)
		}
	}
}

func handleDownload(manager *config.Manager) {
	server, err := selectServer(manager, "选择服务器")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	prompt1 := promptui.Prompt{
		Label: "远程路径",
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("路径不能为空")
			}
			return nil
		},
	}

	remotePath, err := prompt1.Run()
	if err != nil {
		fmt.Printf("输入取消: %v\n", err)
		return
	}

	prompt2 := promptui.Prompt{
		Label: "本地路径",
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("路径不能为空")
			}
			return nil
		},
	}

	localPath, err := prompt2.Run()
	if err != nil {
		fmt.Printf("输入取消: %v\n", err)
		return
	}

	client := ssh.NewClient(server)
	defer client.Close()

	transfer, err := ssh.NewTransfer(client)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	defer transfer.Close()

	files, err := transfer.ListRemote(remotePath)
	if err != nil {
		// 可能是文件而不是目录
		if err := transfer.Download(remotePath, localPath); err != nil {
			fmt.Printf("下载失败: %v\n", err)
		}
		return
	}

	if len(files) == 1 && !files[0].IsDir() {
		if err := transfer.Download(remotePath, localPath); err != nil {
			fmt.Printf("下载失败: %v\n", err)
		}
	} else {
		if err := transfer.DownloadDir(remotePath, localPath); err != nil {
			fmt.Printf("下载失败: %v\n", err)
		}
	}
}

func handleAddServer(manager *config.Manager) {
	// 复用add命令的逻辑
	addCmd.Run(nil, []string{})
}

func handleRemoveServer(manager *config.Manager) {
	removeCmd.Run(nil, []string{})
}

func handleListServers(manager *config.Manager) {
	listCmd.Run(nil, []string{})
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}


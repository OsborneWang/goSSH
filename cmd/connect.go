package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"goSSH/internal/config"
	"goSSH/internal/ssh"
)

var connectCmd = &cobra.Command{
	Use:   "connect [name]",
	Short: "连接到SSH服务器",
	Long:  "连接到指定的SSH服务器，如果未提供名称则交互式选择。连接后将启动交互式Shell。",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := config.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}

		var serverName string
		if len(args) > 0 {
			serverName = args[0]
		} else {
			// 交互式选择服务器
			servers, err := manager.ListServers()
			if err != nil {
				fmt.Fprintf(os.Stderr, "错误: %v\n", err)
				return
			}

			if len(servers) == 0 {
				fmt.Println("没有配置任何服务器，请先使用 'goss add' 添加服务器")
				return
			}

			items := make([]string, len(servers))
			for i, s := range servers {
				items[i] = fmt.Sprintf("%s (%s:%d)", s.Name, s.Host, s.Port)
			}

			prompt := promptui.Select{
				Label: "选择要连接的服务器",
				Items: items,
			}

			index, _, err := prompt.Run()
			if err != nil {
				fmt.Printf("操作取消: %v\n", err)
				return
			}

			serverName = servers[index].Name
		}

		// 获取服务器配置
		server, err := manager.GetServer(serverName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}

		// 创建SSH客户端
		client := ssh.NewClient(server)
		defer client.Close()

		// 连接服务器
		fmt.Printf("正在连接到 %s (%s:%d)...\n", server.Name, server.Host, server.Port)
		if err := client.Connect(); err != nil {
			fmt.Fprintf(os.Stderr, "连接失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ 已连接到 %s\n\n", server.Name)

		// 启动交互式Shell
		executor := ssh.NewExecutor(client)
		if err := executor.ExecuteShell(); err != nil {
			fmt.Fprintf(os.Stderr, "Shell错误: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}


package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"goSSH/internal/config"
	"goSSH/internal/ssh"
)

var execCmd = &cobra.Command{
	Use:   "exec [name] [command]",
	Short: "在远程服务器上执行命令",
	Long:  "在指定的远程服务器上执行命令，如果未提供名称则交互式选择",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := config.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}

		var serverName string
		var command string

		if len(args) >= 2 {
			serverName = args[0]
			command = strings.Join(args[1:], " ")
		} else if len(args) == 1 {
			serverName = args[0]
			// 交互式输入命令
			prompt := promptui.Prompt{
				Label: "要执行的命令",
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("命令不能为空")
					}
					return nil
				},
			}
			var err error
			command, err = prompt.Run()
			if err != nil {
				fmt.Printf("输入取消: %v\n", err)
				return
			}
		} else {
			// 交互式选择服务器和命令
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
				Label: "要执行的命令",
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("命令不能为空")
					}
					return nil
				},
			}
			command, err = prompt2.Run()
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

		// 创建SSH客户端和执行器
		client := ssh.NewClient(server)
		defer client.Close()

		executor := ssh.NewExecutor(client)

		// 执行命令（流式输出）
		if err := executor.ExecuteWithStream(command); err != nil {
			fmt.Fprintf(os.Stderr, "执行失败: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}


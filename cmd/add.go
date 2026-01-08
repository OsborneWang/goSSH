package cmd

import (
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"goSSH/internal/config"
	"goSSH/models"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "添加SSH服务器配置",
	Long:  "交互式添加新的SSH服务器配置",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := config.NewManager()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}

		// 交互式输入服务器信息
		prompt := promptui.Prompt{
			Label: "服务器名称/别名",
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("服务器名称不能为空")
				}
				return nil
			},
		}
		name, err := prompt.Run()
		if err != nil {
			fmt.Printf("输入取消: %v\n", err)
			return
		}

		prompt = promptui.Prompt{
			Label: "主机地址 (IP或域名)",
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("主机地址不能为空")
				}
				return nil
			},
		}
		host, err := prompt.Run()
		if err != nil {
			fmt.Printf("输入取消: %v\n", err)
			return
		}

		prompt = promptui.Prompt{
			Label:   "SSH端口",
			Default: "22",
			Validate: func(input string) error {
				port, err := strconv.Atoi(input)
				if err != nil {
					return fmt.Errorf("端口必须是数字")
				}
				if port < 1 || port > 65535 {
					return fmt.Errorf("端口范围必须在1-65535之间")
				}
				return nil
			},
		}
		portStr, err := prompt.Run()
		if err != nil {
			fmt.Printf("输入取消: %v\n", err)
			return
		}
		port, _ := strconv.Atoi(portStr)

		prompt = promptui.Prompt{
			Label: "用户名",
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("用户名不能为空")
				}
				return nil
			},
		}
		username, err := prompt.Run()
		if err != nil {
			fmt.Printf("输入取消: %v\n", err)
			return
		}

		prompt = promptui.Prompt{
			Label: "密码",
			Mask:  '*',
		}
		password, err := prompt.Run()
		if err != nil {
			fmt.Printf("输入取消: %v\n", err)
			return
		}

		server := models.Server{
			Name:     name,
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
		}

		if err := manager.AddServer(server); err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}

		fmt.Printf("✓ 服务器 '%s' 添加成功\n", name)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}


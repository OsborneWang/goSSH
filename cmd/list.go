package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"goSSH/internal/config"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有已配置的SSH服务器",
	Long:  "显示所有已保存的SSH服务器配置信息",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := config.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}

		servers, err := manager.ListServers()
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}

		if len(servers) == 0 {
			fmt.Println("没有配置任何服务器")
			return
		}

		headerColor := color.New(color.FgCyan, color.Bold)
		headerColor.Printf("\n%-20s %-20s %-8s %-15s\n", "名称", "主机", "端口", "用户名")
		fmt.Println("────────────────────────────────────────────────────────────")

		for _, server := range servers {
			fmt.Printf("%-20s %-20s %-8d %-15s\n",
				server.Name,
				server.Host,
				server.Port,
				server.Username)
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}


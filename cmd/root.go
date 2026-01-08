package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goss",
	Short: "GoSSH - 跨平台SSH命令行工具",
	Long: `GoSSH 是一个使用Go语言开发的跨平台SSH命令行工具。
支持服务器管理、SSH连接、远程命令执行和文件传输功能。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 如果没有提供子命令，显示帮助信息或进入交互式模式
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}


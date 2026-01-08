package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"goSSH/internal/config"
)

var removeCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "删除SSH服务器配置",
	Long:  "删除指定的SSH服务器配置，如果未提供名称则交互式选择",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := config.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}

		var name string
		if len(args) > 0 {
			name = args[0]
		} else {
			// 交互式选择服务器
			servers, err := manager.ListServers()
			if err != nil {
				fmt.Fprintf(os.Stderr, "错误: %v\n", err)
				return
			}

			if len(servers) == 0 {
				fmt.Println("没有可删除的服务器")
				return
			}

			items := make([]string, len(servers))
			for i, s := range servers {
				items[i] = fmt.Sprintf("%s (%s:%d)", s.Name, s.Host, s.Port)
			}

			prompt := promptui.Select{
				Label: "选择要删除的服务器",
				Items: items,
			}

			index, _, err := prompt.Run()
			if err != nil {
				fmt.Printf("操作取消: %v\n", err)
				return
			}

			name = servers[index].Name
		}

		// 确认删除
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("确认删除服务器 '%s'? (y/N)", name),
			Default:   "N",
			AllowEdit: true,
		}
		confirm, err := prompt.Run()
		if err != nil || (confirm != "y" && confirm != "Y" && confirm != "yes" && confirm != "YES") {
			fmt.Println("操作已取消")
			return
		}

		if err := manager.RemoveServer(name); err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			return
		}

		fmt.Printf("✓ 服务器 '%s' 已删除\n", name)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}


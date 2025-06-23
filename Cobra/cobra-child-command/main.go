package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	// 根命令
	rootCmd := &cobra.Command{Use: "sysctl", Short: "系统管理工具"}

	// 在根命令添加全局参数
	rootCmd.PersistentFlags().Bool("debug", false, "调试模式")

	// 创建子命令：user（功能模块入口）
	userCmd := &cobra.Command{Use: "user", Short: "用户管理"}

	// 创建 user 的子命令：add（具体操作）
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "添加用户",
		Run: func(cmd *cobra.Command, args []string) {
			// 实际业务
			name, _ := cmd.Flags().GetString("name")
			fmt.Printf("添加用户：%s\n", name)
		},
	}

	// 为子命令添加参数
	addCmd.Flags().StringP("name", "n", "", "用户名（必须）")
	addCmd.MarkFlagRequired("name")

	// 在user命令添加持久参数
	userCmd.PersistentFlags().BoolP("verbose", "v", false, "详细模式")

	// 在user的所有子命令中都可以访问
	addCmd.Run = func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose") // 获取父命令参数
		fmt.Printf("持久参数 —— 获取父命令参数：%v\n", verbose)

		debugMode, _ := cmd.Root().PersistentFlags().GetBool("debug")
		fmt.Printf("全局参数 —— 任何子命令中都可以获取到：%v\n", debugMode)
	}

	// 构建命令树：父子关系绑定
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(addCmd)

	// 添加删除用户子命令
	delCmd := &cobra.Command{
		Use:   "delete",
		Short: "删除用户",
		Run: func(cmd *cobra.Command, args []string) {
			// 使用args而不是flags
			if len(args) == 0 {
				fmt.Println("错误：需要提供用户名")
				return
			}
			fmt.Printf("删除用户: %s\n", args[0])
		},
	}

	// 添加查看用户子命令
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "列出所有用户",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("用户列表：\n- Alice\n- Bob")
		},
	}

	// 绑定到user命令
	userCmd.AddCommand(delCmd)
	userCmd.AddCommand(listCmd)

	// 创建服务管理父命令
	serviceCmd := &cobra.Command{
		Use:   "service",
		Short: "服务管理",
	}

	// 启动服务子命令
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "启动服务",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			fmt.Printf("启动服务: %s\n", name)
		},
	}
	startCmd.Flags().StringP("name", "n", "", "服务名称")
	startCmd.MarkFlagRequired("name")

	// 绑定到根命令
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(startCmd)

	rootCmd.ExecuteC()
}

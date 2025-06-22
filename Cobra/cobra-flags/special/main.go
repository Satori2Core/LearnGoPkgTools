package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "sptest",
		Short: "特殊场景测试",
		Run: func(cmd *cobra.Command, args []string) {
			output, _ := cmd.Flags().GetString("output")
			theme, _ := cmd.Flags().GetString("theme")

			fmt.Println("输出：", output)
			fmt.Println("主题：", theme)

			// username, _ := cmd.Flags().GetString("username")
			// fmt.Printf("用户输入的用户名是：%s\n", username)
		},
	}

	rootCmd.Flags().StringP("output", "o", "text", "输出格式")
	rootCmd.Flags().StringP("theme", "t", "light", "颜色主题")

	// 参数验证（定义在Run函数前）
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		if output != "text" && output != "json" {
			log.Fatal("输出格式必须是text或json")
		}
	}

	// rootCmd.Flags().StringP("username", "u", "匿名", "用户名（必须）")
	// rootCmd.MarkFlagRequired("username")

	rootCmd.Execute()
}

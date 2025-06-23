package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	simpleText()
}

// simpleText 命令参数的简单认识与使用
func simpleText() {
	// 根命令
	rootCmd := &cobra.Command{
		Use:   "repeat",
		Short: "重复输出文本",

		// 参数绑定后自动注入
		Run: func(cmd *cobra.Command, args []string) {
			// 获取参数值（从 cmd.Flags()）
			text, _ := cmd.Flags().GetString("text")
			count, _ := cmd.Flags().GetInt("count")

			// 核心业务逻辑
			result := strings.Repeat(text+" ", count)
			fmt.Println("输出", strings.TrimSpace(result))
		},
	}

	// 添加文本参数（支持长/短两种模式）
	rootCmd.Flags().StringP("text", "t", "Hello", "要重复的文本")

	// 添加次数参数
	rootCmd.Flags().IntP("count", "c", 3, "重复次数")

	rootCmd.Execute()
}

package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		// Use: "simple",
		// Short: "一个简单的例子",
		// Long:  "这个例子演示了根命令的执行",
		// 设置根命令的Run函数
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("你好，这是根命令在执行！")
		},
	}

	rootCmd.Execute()
}

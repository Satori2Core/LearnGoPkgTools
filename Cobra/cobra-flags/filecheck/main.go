package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "filecheck",
		Short: "文件检查工具",
		Run: func(cmd *cobra.Command, args []string) {
			// 获取所有参数
			path, _ := cmd.Flags().GetString("path")
			recursive, _ := cmd.Flags().GetBool("recursive") // 是否启用递归模式
			minSize, _ := cmd.Flags().GetInt64("min-size")
			exts, _ := cmd.Flags().GetStringSlice("ext")

			fmt.Printf("检查路径: %s\n", path)
			fmt.Printf("递归模式: %v\n", recursive)
			fmt.Printf("最小尺寸: %d字节\n", minSize)
			fmt.Printf("扩展名过滤: %v\n", exts)
		},
	}

	// 支持不同数据类型
	rootCmd.Flags().StringP("path", "p", ".", "检查路径")
	rootCmd.Flags().BoolP("recursive", "r", false, "递归检查")
	rootCmd.Flags().Int64P("min-size", "s", 1024, "最小文件大小（字节）")
	rootCmd.Flags().StringSliceP("ext", "e", []string{}, "按扩展名过滤（可多个）")

	rootCmd.Execute()
}

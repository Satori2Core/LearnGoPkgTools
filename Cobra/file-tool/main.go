package main

import (
	"log"

	"github.com/spf13/cobra"
)

func main() {
	// 1. 创建根命令（顶级命令）
	rootCmd := &cobra.Command{}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

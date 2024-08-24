package main

import (
	"cvgo/commands/add"
	"cvgo/commands/crosscompile"
	"cvgo/commands/daemon"
	"cvgo/commands/enable"
	"cvgo/commands/hotcompile"
	"cvgo/console"
	"cvgo/types"
	"github.com/spf13/cobra" // https://github.com/spf13/cobra
)

// go build -o $GOPATH/bin/cvg
func main() {
	console.LoadConfig()
	RunConsole()
}

func RunConsole() {
	var cobraRoot = &cobra.Command{
		// 定义根命令的关键字
		Use: "cvg",
		// 简短介绍
		Short: "Cvgo 配套开发工具",
		// 根命令的执行函数
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
	}
	var command = &types.Command{
		RootCmd: cobraRoot,
	}
	// 绑定指令
	hotcompile.AddCommand(command)
	daemon.AddCommand(command)
	crosscompile.AddCommand(command)
	add.AddCommand(command)
	enable.AddCommand(command)

	// 命令行运行，执行 RootCommand
	command.RootCmd.Execute()
}

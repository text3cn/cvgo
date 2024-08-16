package create

import (
	"bufio"
	"cvgo/console/types"
	"cvgo/kit/filekit"
	"cvgo/provider"
	"cvgo/provider/plog"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var log = provider.Services.NewSingle(plog.Name).(plog.Service)

func AddCommand(command *types.Command) {
	// 一级命令
	lv1 := &cobra.Command{
		Use:   "create",
		Short: "创建相关文件",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("命令不完整，请执行 cvg create --help 查看帮助")
			}
		},
	}

	// 绑定 flags
	var webserver string
	lv1.PersistentFlags().StringVar(&webserver, "webserver", "", "使用 Web 服务框架")

	// 二级命令
	lv1.AddCommand(&cobra.Command{
		Use:     "module",
		Short:   "创建模块",
		Aliases: []string{"mod"},
		Example: "cvg create module <module1> <mudoule2>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("命令不完整，请执行 cvg create module --help 查看帮助")
				return
			}
			if goOn := deleteBeforeCreate(); goOn == false {
				fmt.Println("你取消了操作")
				return
			}
			for _, arg := range args {
				createModule(arg)
			}
		},
	})
	command.RootCmd.AddCommand(lv1)
}

func deleteBeforeCreate() (delete bool) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(plog.YellowSprintf("⚠️ 当前开启了创建模块之前先删除 app 目录和 go.work 文件。是否继续操作？(yes/no) [default:" + plog.BlueSprintf("yes", plog.ColorCyan) + "]:"))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "yes" {
		delete = true
	} else if input == "no" {
		delete = false
	} else {
		delete = true
	}
	if delete {
		filekit.DeleteDirOrFile("./app")
		filekit.DeleteDirOrFile("./go.work")
	}
	return
}

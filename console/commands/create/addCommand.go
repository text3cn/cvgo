package create

import (
	"cvgo/console/types"
	"cvgo/provider"
	"cvgo/provider/clog"
	"fmt"
	"github.com/spf13/cobra"
)

var log = provider.Services.NewSingle(clog.Name).(clog.Service)

//var cfg = provider.Services.NewSingle(config.Name).(config.Service)

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
	var force bool
	var swagger bool
	lv1.PersistentFlags().StringVar(&webserver, "webserver", "", "使用 Web 框架")
	lv1.PersistentFlags().BoolVar(&force, "force", false, "强制创建，如果模块已存在会先删除")
	lv1.PersistentFlags().BoolVar(&swagger, "swagger", false, "Swagger 文档支持")

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
			for _, moduleName := range args {
				createModule(moduleName, webserver, swagger, force)
			}
		},
	})
	command.RootCmd.AddCommand(lv1)
}

package create

import (
	"cvgo/console/internal/types"
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
				if moduleName == "cvgo" {
					log.Error("名称冲突，模块名称不能是 cvgo")
					return
				}
				createModule(moduleName, webserver, swagger, force)
			}
		},
	})

	// mysql entity
	lv1.AddCommand(&cobra.Command{
		Use:     "mysqlEntity",
		Short:   "创建模型 entity",
		Example: "cvg create mysqlEntity 表名称 表注释(可选)",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				clog.CyanPrintln("命令不完整，查看帮助请执行 cvg create mysqlEntity -h")
				return
			}
			tableName := args[0]
			comment := ""
			if len(args) > 1 {
				comment = args[1]
			}
			createMysqlEntity(tableName, comment)
		},
	})

	// docker
	lv1.AddCommand(&cobra.Command{
		Use:     "dockerEnv",
		Short:   "添加一个 docker-compose 模板文件到 app",
		Example: "cvg create dockerEnv",
		Run: func(cmd *cobra.Command, args []string) {
			createDockerEnv()
		},
	})

	command.RootCmd.AddCommand(lv1)
}

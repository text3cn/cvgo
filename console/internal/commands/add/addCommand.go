package add

import (
	"cvgo/console/internal/commands/add/gencvgo"
	"cvgo/console/internal/commands/add/scripts"
	"cvgo/console/internal/console"
	"cvgo/console/internal/paths"
	"cvgo/console/internal/types"
	"cvgo/kit/arrkit"
	"cvgo/kit/filekit"
	"cvgo/kit/strkit"
	"cvgo/provider"
	"cvgo/provider/clog"
	"fmt"
	"github.com/spf13/cobra"
	"unicode"
)

var log = provider.Services.NewSingle(clog.Name).(clog.Service)

//var cfg = provider.Services.NewSingle(config.Name).(config.Service)

func AddCommand(command *types.Command) {
	// 一级命令
	lv1 := &cobra.Command{
		Use:   "add",
		Short: "创建相关文件",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("命令不完整，请执行 cvg add --help 查看帮助")
			}
		},
	}

	// 二级命令
	module := &cobra.Command{
		Use:     "module",
		Short:   "创建模块",
		Aliases: []string{"mod"},
		Example: "cvg add module <module1> <mudoule2>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("命令不完整，请执行 cvg add module --help 查看帮助")
				return
			}
			webserver, _ := cmd.Flags().GetString("webserver")
			swagger, _ := cmd.Flags().GetBool("swagger")
			force, _ := cmd.Flags().GetBool("force")
			for _, moduleName := range args {
				if moduleName == "cvgo" {
					log.Error("名称冲突，模块名称不能是 cvgo")
					return
				}
				createModule(moduleName, webserver, swagger, force)
			}
		},
	}
	// 绑定 flags
	var webserver string
	var force bool
	var swagger bool
	module.Flags().StringVar(&webserver, "webserver", "", "使用 Web 框架")
	module.Flags().BoolVar(&force, "force", false, "强制创建，如果模块已存在会先删除")
	module.Flags().BoolVar(&swagger, "swagger", false, "Swagger 文档支持")
	lv1.AddCommand(module)

	// mysql entity
	lv1.AddCommand(&cobra.Command{
		Use:     "table",
		Short:   "创建模型 entity",
		Example: "cvg add table 表名称 表注释(可选)",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				clog.CyanPrintln("命令不完整，查看帮助请执行 cvg add mysqlEntity -h")
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
		Example: "cvg add dockerEnv",
		Run: func(cmd *cobra.Command, args []string) {
			createDockerEnv()
		},
	})

	// script
	lv1.AddCommand(&cobra.Command{
		Use:     "script",
		Short:   "添加脚本模板",
		Example: "cvg add script <脚本模板>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				clog.CyanPrintln("命令不完整，查看帮助请执行 cvg add mysqlEntity -h")
				return
			}
			scripts.AddScript(args[0])
		},
	})

	api := &cobra.Command{
		Use:     "api",
		Short:   "一键创建 route、api、dto",
		Example: "cvg add api <path> ",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				clog.RedPrintln("命令不完整，查看帮助请执行 cvg add api -h")
				return
			}
			if !paths.CheckRunAtModuleRoot() {
				return
			}
			// 检查 method
			allowMethod := []string{"get", "post", "put", "patch", "delete"}
			if !arrkit.InArray(args[0], allowMethod) {
				clog.RedPrintln("不支持的请求方法", args[0], "，支持的方法为：get / post / put / delete")
				return
			}
			if len(args) < 2 {
				clog.BrownPrintln("请指定 path（http 请求路径）")
				return
			}
			webFramework, _ := console.NewKvStorage(filekit.GetParentDir(3)).GetWebFramework()
			supportSwagger, _ := console.NewKvStorage(filekit.GetParentDir(3)).GetSwagger()
			cvgflag, _ := cmd.Flags().GetString("cvgflag")
			switch webFramework {
			case "cvgo":
				gencvgo.GenApi(args[0], args[1], supportSwagger, cvgflag)
			}
		},
	}
	var cvgflag string
	api.Flags().StringVar(&cvgflag, "cvgflag", "", "路由生成标记")
	lv1.AddCommand(api)

	// service
	svc := &cobra.Command{
		Use:     "svc",
		Short:   "创建 service",
		Example: "cvg add svc <file/func> --curd=user",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				clog.RedPrintln("命令不完整，查看帮助请执行 cvg add api -h")
				return
			}
			pathArr := strkit.Explode("/", args[0])
			if len(pathArr) != 2 {
				clog.RedPrintln("请求路径格式错误，路径示例：user/GetUserinfo")
				return
			}
			webFramework, _ := console.NewKvStorage(filekit.GetParentDir(3)).GetWebFramework()
			table, _ := cmd.Flags().GetString("table")
			curdType := ""
			switch webFramework {
			case "cvgo":
				if len(args) > 1 {
					curdType = args[1]
					for _, v := range curdType {
						u := string(unicode.ToLower(v))
						funcName := ""
						switch u {
						case "c":
							funcName = "Create" + pathArr[1]
						case "u":
							funcName = "Update" + pathArr[1]
						case "r":
							funcName = "Get" + pathArr[1]
						case "d":
							funcName = "Delete" + pathArr[1]
						case "l":
							funcName = "List" + pathArr[1]
						}
						gencvgo.GenService(pathArr[0], funcName, u, table)
					}
				} else {
					gencvgo.GenService(pathArr[0], pathArr[1], curdType, table)
				}

			}
		},
	}
	var table string
	svc.Flags().StringVar(&table, "table", "", "使用指定表名称创建 CURD 代码")
	lv1.AddCommand(svc)

	command.RootCmd.AddCommand(lv1)
}

package add

import (
	"cvgo/commands/codegen/gencvgo"
	"github.com/textthree/cvgokit/strkit"
	"unicode"

	"cvgo/kvs"
	"cvgo/types"
	"github.com/spf13/cobra"
	"github.com/textthree/cvgokit/arrkit"
	"github.com/textthree/provider/clog"
)

func AddCommand(command *types.Command) {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "创建 route、api、dto、service",
		Run: func(cmd *cobra.Command, args []string) {
			clog.RedPrintln("命令不完整，查看帮助请执行 cvg add -h")
		},
	}

	// 创建 api
	api := &cobra.Command{
		Use:     "api",
		Short:   "一键创建 route、api、dto",
		Example: "cvg add api get api/user/info",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				clog.RedPrintln("命令不完整，查看帮助请执行 cvg add api -h")
				return
			}
			kv := kvs.Instance()
			kv.CheckInModuleDir()

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
			webFramework := kv.GetWebFramework()
			supportSwagger, _ := kv.GetSwagger()
			cvgflag, _ := cmd.Flags().GetString("cvgflag")
			switch webFramework {
			case "cvgo":
				gencvgo.GenApi(args[0], args[1], supportSwagger, cvgflag)
			default:
				clog.CyanPrintln(webFramework + " 暂不支持，目前只支持 cvgo")
			}
		},
	}
	var cvgflag string
	api.Flags().StringVar(&cvgflag, "cvgflag", "", "路由生成标记")
	addCmd.AddCommand(api)

	// service
	svc := &cobra.Command{
		Use:     "svc",
		Short:   "添加 service 方法",
		Example: "cvg add svc <file/func> --curd=user",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				clog.RedPrintln("命令不完整，查看帮助请执行 cvg codegen api -h")
				return
			}
			pathArr := strkit.Explode("/", args[0])
			if len(pathArr) != 2 {
				clog.RedPrintln("请求路径格式错误，路径示例：user/GetUserinfo")
				return
			}
			kv := kvs.Instance()
			webFramework := kv.GetWebFramework()
			tableName, _ := cmd.Flags().GetString("table")
			cursorPaging, _ := cmd.Flags().GetBool("cursor")
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
						gencvgo.GenService(pathArr[0], funcName, u, tableName, cursorPaging)
					}
				} else {
					gencvgo.GenService(pathArr[0], pathArr[1], curdType, tableName, cursorPaging)
				}
			default:
				clog.CyanPrintln("目前只支持对 Cvgo Web 生成 service")
			}
		},
	}
	var tableName string
	var cursorPaging bool
	svc.Flags().StringVar(&tableName, "table", "", "使用指定表名称创建 CURD 代码")
	svc.Flags().BoolVar(&cursorPaging, "cursor", false, "列表是否使用游标分页")
	addCmd.AddCommand(svc)

	command.RootCmd.AddCommand(addCmd)
}

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
			pass := false
			useCurdl := false
			allowMethod := []string{"get", "post", "put", "patch", "delete"}
			if arrkit.InArray(args[0], allowMethod) {
				pass = true
			}
			// 检查 curdl 组合
			if !pass {
				pass = true
				useCurdl = true
				allowCurdl := []string{"c", "u", "r", "d", "l"}
				for _, v := range args[0] {
					if !arrkit.InArray(string(v), allowCurdl) {
						pass = false
					}
				}
			}
			if !pass {
				clog.RedPrintln("不支持的请求方法", args[0], "，支持的方法为：get / post / put / delete，或 curdl 组合")
			}
			var table string
			if useCurdl {
				table, _ = cmd.Flags().GetString("table")
				if table == "" {
					clog.RedPrintln("请使用 --table 选项指定表名称")
					return
				}
			}
			if len(args) < 2 {
				clog.BrownPrintln("请指定 path（http 请求路径）")
				return
			}
			webFramework := kv.GetWebFramework()
			supportSwagger, _ := kv.GetSwagger()
			cvgflag, _ := cmd.Flags().GetString("cvgflag")
			cursorPaging, _ := cmd.Flags().GetBool("cursor")

			switch webFramework {
			case "cvgo":
				if table != "" {
					for _, v := range args[0] {
						pathArr := strkit.Explode("/", args[1])
						switch string(v) {
						case "c":
							svcFuncName := "Create" + strkit.Ucfirst(pathArr[1]) + strkit.Ucfirst(pathArr[2])
							svcName := gencvgo.GenService(pathArr[0], svcFuncName, "c", table, cursorPaging)
							if svcName == "" {
								table = ""
							}
							gencvgo.GenApi("post", args[1], supportSwagger, cursorPaging, cvgflag, table, svcName, svcFuncName, "c")
						case "u":
							svcFuncName := "Update" + strkit.Ucfirst(pathArr[1]) + strkit.Ucfirst(pathArr[2])
							svcName := gencvgo.GenService(pathArr[1], svcFuncName, "u", table, cursorPaging)
							if svcName == "" {
								table = ""
							}
							gencvgo.GenApi("put", args[1], supportSwagger, cursorPaging, cvgflag, table, svcName, svcFuncName, "u")
						case "r":
							svcFuncName := "Get" + strkit.Ucfirst(pathArr[1]) + strkit.Ucfirst(pathArr[2])
							svcName := gencvgo.GenService(pathArr[1], svcFuncName, "r", table, cursorPaging)
							if svcName == "" {
								table = ""
							}
							gencvgo.GenApi("get", args[1], supportSwagger, cursorPaging, cvgflag, table, svcName, svcFuncName, "r")
						case "d":
							svcFuncName := "Delete" + strkit.Ucfirst(pathArr[1]) + strkit.Ucfirst(pathArr[2])
							svcName := gencvgo.GenService(pathArr[1], svcFuncName, "d", table, cursorPaging)
							if svcName == "" {
								table = ""
							}
							gencvgo.GenApi("delete", args[1], supportSwagger, cursorPaging, cvgflag, table, svcName, svcFuncName, "d")
						case "l":
							svcFuncName := "List" + strkit.Ucfirst(pathArr[1]) + strkit.Ucfirst(pathArr[2])
							svcName := gencvgo.GenService(pathArr[1], svcFuncName, "l", table, cursorPaging)
							if svcName == "" {
								table = ""
							}
							requestPath := args[1] + "List"
							gencvgo.GenApi("post", requestPath, supportSwagger, cursorPaging, cvgflag, table, svcName, svcFuncName, "l")
						}
					}
				} else {
					gencvgo.GenApi(args[0], args[1], supportSwagger, cursorPaging, cvgflag)
				}

			default:
				clog.CyanPrintln(webFramework + " 暂不支持，目前只支持 cvgo")
			}
		},
	}
	var cvgflag string
	var table string
	var cursor bool
	api.Flags().StringVar(&cvgflag, "cvgflag", "", "路由生成标记")
	api.Flags().StringVar(&table, "table", "", "同时生成 service 方法")
	api.Flags().BoolVar(&cursor, "cursor", false, "列表是否使用游标分页")
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
							funcName = "Create" + strkit.Ucfirst(pathArr[1])
						case "u":
							funcName = "Update" + strkit.Ucfirst(pathArr[1])
						case "r":
							funcName = "Get" + strkit.Ucfirst(pathArr[1])
						case "d":
							funcName = "Delete" + strkit.Ucfirst(pathArr[1])
						case "l":
							funcName = "List" + strkit.Ucfirst(pathArr[1])
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

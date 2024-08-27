package create

import (
	"cvgo/commands/create/docker"
	"cvgo/commands/create/gitlabci"
	"cvgo/commands/create/module"
	"cvgo/commands/create/scripts"
	"cvgo/commands/create/table"
	"cvgo/commands/create/work"
	"cvgo/kvs"
	"cvgo/kvs/kvsKey"
	"cvgo/types"
	"github.com/spf13/cobra"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/provider/clog"
)

func AddCommand(command *types.Command) {
	// Create workspace
	work := &cobra.Command{
		Use:     "create-work",
		Aliases: []string{"cw"},
		Example: "cvg create-work myProject",
		Short:   "Create workspace",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				help := clog.CyanSprintf("cvg create-work -h", clog.ColorRed)
				clog.RedPrintln("The command is incomplete. To view the help, run ", help)
			}
			workspaceDirName := args[0]
			if filekit.DirExists(workspaceDirName) {
				clog.RedPrintln(workspaceDirName + " is not an empty directory.")
				return
			}
			work.CreateWork(workspaceDirName)
		},
	}
	command.RootCmd.AddCommand(work)

	// Create module
	module := &cobra.Command{
		Use:     "create-module",
		Short:   "Create module",
		Aliases: []string{"cm"},
		Example: "cvg codegen module <module1> <mudoule2>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				help := clog.CyanSprintf("cvg create-module -h", clog.ColorRed)
				clog.RedPrintln("The command is incomplete. To view the help, run ", help)
				return
			}
			webserver, _ := cmd.Flags().GetString("webserver")
			swagger, _ := cmd.Flags().GetBool("swagger")
			force, _ := cmd.Flags().GetBool("force")
			workspaceName, _ := kvs.Instance().GetString(kvsKey.WorkspaceName)
			for _, moduleName := range args {
				if moduleName == workspaceName {
					clog.RedPrintln("The module name cannot be cvgo because of workspace name conflict")
					return
				}
				module.CreateModule(moduleName, webserver, swagger, force)
			}
		},
	}
	var webserver string
	var force bool
	var swagger bool
	module.Flags().StringVar(&webserver, "webserver", "", "使用 Web 框架")
	module.Flags().BoolVar(&force, "force", false, "强制创建，如果模块已存在会先删除")
	module.Flags().BoolVar(&swagger, "swagger", false, "Swagger 文档支持")
	command.RootCmd.AddCommand(module)

	// Create table
	table := &cobra.Command{
		Use:     "create-table",
		Aliases: []string{"ct"},
		Example: "cvg create-table user 用户表",
		Short:   "创建 Gorm 表实体",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				help := clog.CyanSprintf("cvg create-table -h", clog.ColorRed)
				clog.RedPrintln("The command is incomplete. To view the help, run ", help)
			}
			tableName := args[0]
			comment := ""
			if len(args) > 1 {
				comment = args[1]
			}
			table.CreateMysqlEntity(tableName, comment)
		},
	}
	command.RootCmd.AddCommand(table)

	// docker-compose
	dockerCompose := &cobra.Command{
		Use:     "create-docker",
		Aliases: []string{"cd"},
		Short:   "添加一个 docker-compose 模板到 app",
		Example: "cvg create-docker",
		Run: func(cmd *cobra.Command, args []string) {
			docker.CreateDocker()
		},
	}
	command.RootCmd.AddCommand(dockerCompose)

	// script
	script := &cobra.Command{
		Use:     "create-script",
		Short:   "添加脚本模板",
		Example: "cvg create-script <脚本模板>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				clog.CyanPrintln("命令不完整，查看帮助请执行 cvg codegen mysqlEntity -h")
				return
			}
			scripts.CreateScript(args[0])
		},
	}
	command.RootCmd.AddCommand(script)

	// gitlab-ci.yml
	gitlabCI := &cobra.Command{
		Use:     "create-gitlab-ci",
		Short:   "添加一个 gitlab-ci.yml 到工作区",
		Example: "cvg create-gitlab-ci",
		Run: func(cmd *cobra.Command, args []string) {
			gitlabci.CreateGitlabCiYml()
		},
	}
	command.RootCmd.AddCommand(gitlabCI)

}

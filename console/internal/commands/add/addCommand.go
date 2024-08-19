package add

import (
	"cvgo/console/internal/types"
	"cvgo/kit/filekit"
	"cvgo/provider"
	"cvgo/provider/clog"
	"github.com/spf13/cobra"
)

var log = provider.Services.NewSingle(clog.Name).(clog.Service)
var pwd string

// 在模块目录中执行
// cd ../../../ && go build -o $GOPATH/bin/cvg ./console && cd app/modules/client && cvg add i18n
func AddCommand(command *types.Command) {
	pwd = filekit.Getwd()
	// 一级命令
	cmd := &cobra.Command{
		Use:   "add",
		Short: "创建相关文件",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				clog.CyanPrintln("命令不完整，请执行 cvg add --help 查看帮助")
			}
		},
	}

	// 二级命令
	cmd.AddCommand(&cobra.Command{
		Use:     "i18n",
		Short:   "添加多语言支持",
		Example: "cvg add i18n",
		Run: func(cmd *cobra.Command, args []string) {
			addI18n()
		},
	})
	command.RootCmd.AddCommand(cmd)
}

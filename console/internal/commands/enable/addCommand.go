package enable

import (
	"cvgo/console/internal/types"
	"cvgo/kit/filekit"
	"cvgo/provider/clog"
	"github.com/spf13/cobra"
)

// 在模块目录中执行
// cd ../../../ && go build -o $GOPATH/bin/cvg ./console && cd app/modules/chord && cvg add i18n
func AddCommand(command *types.Command) {
	pwd = filekit.Getwd()
	// 一级命令
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "开启相关功能支持",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				clog.CyanPrintln("命令不完整，请执行 cvg add --help 查看帮助")
			}
		},
	}

	// 二级命令
	cmd.AddCommand(&cobra.Command{
		Use:     "i18n",
		Short:   "开启多语言支持",
		Example: "cvg add i18n",
		Run: func(cmd *cobra.Command, args []string) {
			addI18n()
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:     "mysql",
		Short:   "开启 mysql 支持",
		Example: "cvg add mysql",
		Run: func(cmd *cobra.Command, args []string) {
			addGorm()
		},
	})

	command.RootCmd.AddCommand(cmd)
}

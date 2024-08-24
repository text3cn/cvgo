package enable

import (
	"cvgo/types"
	"github.com/spf13/cobra"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/provider/clog"
)

func AddCommand(command *types.Command) {
	pwd = filekit.Getwd()

	// 一级命令
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "开启相关功能支持",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				clog.CyanPrintln("命令不完整，请执行 cvg enable --h 查看帮助")
			}
		},
	}

	// 二级命令
	cmd.AddCommand(&cobra.Command{
		Use:     "i18n",
		Short:   "开启多语言支持",
		Example: "cvg enable i18n",
		Run: func(cmd *cobra.Command, args []string) {
			addI18n()
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:     "mysql",
		Short:   "开启 mysql 支持",
		Example: "cvg enable mysql",
		Run: func(cmd *cobra.Command, args []string) {
			addMysql()
		},
	})

	command.RootCmd.AddCommand(cmd)
}

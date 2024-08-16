package crosscompile

import (
	"cvgo/console/console"
	"cvgo/console/types"
	"cvgo/kit/filekit"
	"cvgo/kit/gokit"
	"github.com/silenceper/log"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
)

var buildFileFullPath string
var moduleName string

func AddCommand(command *types.Command) {
	dev := &cobra.Command{
		Use:   "build",
		Short: "自动编译",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				platform := args[0]
				if platform != "linux" && platform != "windows" && platform != "mac" {
					log.Error("Platform not support")
					return
				}

				// 创建目录
				if console.CrossCompileCfg.OutputDir != "./" {
					filekit.MkDir(console.CrossCompileCfg.OutputDir, 0777)
				}
				var err error
				moduleName, err = gokit.GetModuleName()
				if err != nil {
					log.Error(err.Error())
				}
				if err != nil {
					log.Errorf("Get module name fail: %s\n", err)
				}

				buildFileFullPath = filepath.Join(console.CrossCompileCfg.OutputDir, moduleName, moduleName)

				// 构建
				pkg := ""
				if len(args) > 1 {
					pkg = args[1]
				}

				switch platform {
				case "linux":
					build("linux", "amd64", pkg)
				case "windows":
					build("windows", "amd64", pkg)
				case "mac":
					build("darwin", "arm64", pkg)
				}
			}
		},
	}
	command.RootCmd.AddCommand(dev)
}

func build(goods, goarch, pkg string) {
	// 设置命令和参数
	if goods == "windows" {
		buildFileFullPath += ".exe"
	}
	args := []string{"build", "-ldflags=-s -w", "-o", buildFileFullPath}
	if pkg != "" {
		args = append(args, pkg)
	}
	cmd := exec.Command("go", args...)

	// 设置环境变量
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS="+goods, "GOARCH="+goarch)

	// 设置输出到标准输出和标准错误
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	err := cmd.Run()
	if err != nil {
		log.Errorf("Build fail: %s\n", err)
		return
	}
	copyConfig()
	log.Infof("Build %s %s successed.", goods, goarch)
}

// TODO 对 internal 里面的配置文件提取出来进行合并到最外层
// TODO 如果存在 i18n 目录，则将其拷贝过去
// viper 提供了 viper.WriteConfigAs("new_config.yaml") 将内存中的配置再存为文件。
func copyConfig() {
	// 公共配置
	src := filepath.Join(console.RootPath, "app", "config")
	tirgetDir := filepath.Join(console.CrossCompileCfg.OutputDir, moduleName)
	filekit.CopyFiles(src, tirgetDir)
	filekit.Rename(filepath.Join(tirgetDir, "alpha"), filepath.Join(tirgetDir, "config"))
	// 模块配置
	internalAppYaml := filepath.Join(console.RootPath, "app", "modules", moduleName, "internal", "config", "app.yaml")
	distAppYaml := filepath.Join(console.CrossCompileCfg.OutputDir, moduleName, "config", "local", "app.yaml")
	filekit.CopyFile(internalAppYaml, distAppYaml)
	// 将 config 用 internal 目录包装
	filekit.MoveDir(filepath.Join(tirgetDir, "config"), filepath.Join(tirgetDir, "internal"))
}

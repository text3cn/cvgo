package paths

import (
	"cvgo/kvs"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/cvgokit/gokit"
	"os"
	"path/filepath"
)

var commandsDir = filepath.Join("console", "internal", "commands")

// 判断是否在模块根目录执行命令
func CheckRunAtModuleRoot() {
	modName, err := gokit.GetModuleName()
	if err != nil || modName == kvs.Instance().GetWorkspaceName() {
		panic("请在模目录下执行此命令（go.mod 文件所在目录）")
	}
}

// 判断是否在工程根目录执行命令
func CheckRunAtProjectRoot() bool {
	if exists, _ := filekit.PathExists(filepath.Join(filekit.Getwd(), "go.work")); exists {
		return true
	}
	return false
}

// 支持在工程任何目录下切换到工程根目录
func CdToWorkspacePath() {
	if !CheckRunAtProjectRoot() {
		err := os.Chdir(kvs.Instance().GetWorkspacePath())
		if err != nil {
			panic(err)
		}
	}
}

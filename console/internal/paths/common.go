package paths

import (
	"cvgo/console/internal/console"
	"cvgo/kit/filekit"
	"cvgo/kit/gokit"
	"cvgo/provider"
	"cvgo/provider/clog"
	"os"
	"path/filepath"
)

var log = provider.Clog()

var commandsDir = filepath.Join("console", "internal", "commands")

var FiberTplForRoot = filepath.Join("console", "internal", "commands", "add", "tpl", "module_tpl", "fiber")

// 当前路径为模块目录时，获取 fiber 模板
func FiberTplForModule() string {
	rootPath := filekit.GetParentDir(3)
	ret := filepath.Join(rootPath, FiberTplForRoot)
	return ret
}

// 判断是否在模块根目录执行命令
func CheckRunAtModuleRoot() bool {
	modName, err := gokit.GetModuleName()
	if err != nil || modName == "cvgo" {
		clog.CyanPrintln("请在模块根目录下执行此命令（go.mod 文件所在目录）")
		return false
	}
	return true
}

// 判断是否在工程根目录执行命令
func CheckRunAtProjectRoot() bool {
	if exists, _ := filekit.PathExists(filepath.Join(filekit.Getwd(), "go.work")); exists {
		return true
	}
	return false
}

func GetProjectRootPathFromKv() string {
	findGoWork := func(path string) bool {
		if exists, _ := filekit.PathExists(filepath.Join(path, "go.work")); exists {
			return true
		}
		return false
	}
	root := filekit.Getwd()
	for i := 0; i < 5; i++ {
		if !findGoWork(root) {
			root = filepath.Dir(filekit.Getwd())
		} else {
			break
		}
	}
	kv := console.NewKvStorage(root)
	ret, err := kv.GetString("projectRootPath")
	if err != nil {
		log.Error("从 runtime.json 获取根路径失败", err.Error())
	}
	return ret
}

// 支持在工程任何目录下执行工具命令时先切换到工程根目录
func CdToProjectRootPath() {
	if !CheckRunAtProjectRoot() {
		err := os.Chdir(GetProjectRootPathFromKv())
		if err != nil {
			panic(err)
		}
	}
}

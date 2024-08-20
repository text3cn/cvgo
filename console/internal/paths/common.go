package paths

import (
	"cvgo/kit/filekit"
	"cvgo/kit/gokit"
	"path/filepath"
)

var commandsDir = filepath.Join("console", "internal", "commands")

var FiberTplForRoot = filepath.Join("console", "internal", "commands", "create", "sampletpl", "module_tpl", "fiber")

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
		return false
	}
	return true
}

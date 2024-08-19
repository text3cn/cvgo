package consolepath

import (
	"cvgo/kit/filekit"
	"path/filepath"
)

var FiberTplForRoot = filepath.Join("console", "internal", "commands", "create", "sampletpl", "module_tpl", "fiber")

// 当前路径为模块目录时，获取 fiber 模板
func FiberTplForModule() string {
	rootPath := filekit.GetParentDir(3)
	ret := filepath.Join(rootPath, FiberTplForRoot)
	return ret
}

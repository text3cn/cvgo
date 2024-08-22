package paths

import (
	"cvgo/kit/filekit"
	"path/filepath"
)

// 当前路径为工程根目录
type pathForProjectRoot struct {
	projectRootPath string // 工程根目录
	commandsPath    string // console 根目录
}

func NewPathForProjectRoot() *pathForProjectRoot {
	projectRootPath := filekit.Getwd()
	return &pathForProjectRoot{
		projectRootPath: filekit.Getwd(),
		commandsPath:    filepath.Join(projectRootPath, commandsDir),
	}
}

func (p *pathForProjectRoot) AppEntityMysql() string {
	ret := filepath.Join(p.projectRootPath, "app", "entity", "mysql")
	return ret
}

// autoMigrate.go
func (p *pathForProjectRoot) AppAutoMigrate() string {
	ret := filepath.Join(p.projectRootPath, "app", "entity", "autoMigrate.go")
	return ret
}

// scripts tpl
func (p *pathForProjectRoot) ScriptsTplDir() string {
	ret := filepath.Join(p.commandsPath, "add", "tpl", "scripts")
	return ret
}

// app scripts
func (p *pathForProjectRoot) AppScriptsDir() string {
	ret := filepath.Join(p.projectRootPath, "app", "scripts")
	return ret
}

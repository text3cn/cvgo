package paths

import (
	"cvgo/kit/filekit"
	"path/filepath"
)

// 当前路径为模块目录
type pathForApp struct {
	rootPath     string // 工程根目录
	commandsPath string // console 根目录
	moduleDir    string // 模块目录
}

func NewPathForApp() *pathForApp {
	rootPath := filekit.GetParentDir(3)
	return &pathForApp{
		rootPath:     rootPath,
		commandsPath: commandsDir,
		moduleDir:    filekit.Getwd(),
	}
}

func (p *pathForApp) MysqlBaseEntityTpl() string {
	ret := filepath.Join(p.rootPath, p.commandsPath, "enable", "tpl", "gorm", "entity", "mysql", "base.go")
	return ret
}

func (p *pathForApp) AppEntityMysqlBaseGoFile() string {
	ret := filepath.Join(p.rootPath, "app", "entity", "mysql", "base.go")
	return ret
}

// instance.go
func (p *pathForApp) InstanceGo() string {
	ret := filepath.Join(p.rootPath, "app", "instance.go")
	return ret
}

// boot -> init.go
func (p *pathForApp) BootInitGo() string {
	ret := filepath.Join(p.moduleDir, "internal", "boot", "init.go")
	return ret
}

// app entity mysql
func (p *pathForApp) AppEntityMysql() string {
	ret := filepath.Join(p.rootPath, "app", "entity", "mysql")
	return ret
}

// autoMigrate.go tpl
func (p *pathForApp) AutoMigrateTpl() string {
	ret := filepath.Join(p.rootPath, p.commandsPath, "enable", "tpl", "gorm", "entity", "autoMigrate.go")
	return ret
}

// autoMigrate.go
func (p *pathForApp) AppAutoMigrate() string {
	ret := filepath.Join(p.rootPath, "app", "entity", "autoMigrate.go")
	return ret
}

// database.yaml tpl
func (p *pathForApp) DatabaseYamlTpl() string {
	ret := filepath.Join(p.rootPath, p.commandsPath, "enable", "tpl", "gorm", "database.yaml")
	return ret
}

// database.yaml
func (p *pathForApp) DatabaseYaml() string {
	ret := filepath.Join(p.rootPath, "app", "config", "database.yaml")
	return ret
}

// database-alpha.yaml tpl
func (p *pathForApp) DatabaseAlphaYamlTpl() string {
	ret := filepath.Join(p.rootPath, p.commandsPath, "enable", "tpl", "gorm", "database-alpha.yaml")
	return ret
}

// database-alpha.yaml
func (p *pathForApp) DatabaseAlphaYaml() string {
	ret := filepath.Join(p.rootPath, "app", "config", "alpha", "database.yaml")
	return ret
}

// database-release.yaml tpl
func (p *pathForApp) DatabaseReleaseYamlTpl() string {
	ret := filepath.Join(p.rootPath, p.commandsPath, "add", "tpl", "gorm", "database-release.yaml")
	return ret
}

// database-release.yaml
func (p *pathForApp) DatabaseReleaseYaml() string {
	ret := filepath.Join(p.rootPath, "app", "config", "release", "database-release.yaml")
	return ret
}

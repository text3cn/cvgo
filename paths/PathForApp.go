package paths

import (
	"github.com/textthree/cvgokit/filekit"
	"path/filepath"
)

// 当前路径为模块目录
type pathForModule struct {
	rootPath     string // 工程根目录
	commandsPath string // console 根目录
	moduleDir    string // 模块目录
}

func NewPathForModule() *pathForModule {
	rootPath := filekit.GetParentDir(3)
	return &pathForModule{
		rootPath:     rootPath,
		commandsPath: commandsDir,
		moduleDir:    filekit.Getwd(),
	}
}

func (p *pathForModule) MysqlBaseEntityTpl() string {
	ret := filepath.Join(p.rootPath, p.commandsPath, "enable", "tpl", "gorm", "entity", "mysql", "base.go")
	return ret
}

func (p *pathForModule) AppEntityMysqlBaseGoFile() string {
	ret := filepath.Join(p.rootPath, "app", "entity", "mysql", "base.go")
	return ret
}

// instance.go
func (p *pathForModule) InstanceGo() string {
	ret := filepath.Join(p.rootPath, "app", "instance.go")
	return ret
}

// boot -> init.go
func (p *pathForModule) BootInitGo() string {
	ret := filepath.Join(p.moduleDir, "internal", "boot", "init.go")
	return ret
}

// app entity mysql
func (p *pathForModule) AppEntityMysql() string {
	ret := filepath.Join(p.rootPath, "app", "entity", "mysql")
	return ret
}

// autoMigrate.go tpl
func (p *pathForModule) AutoMigrateTpl() string {
	ret := filepath.Join(p.rootPath, p.commandsPath, "enable", "tpl", "gorm", "entity", "autoMigrate.go")
	return ret
}

// autoMigrate.go
func (p *pathForModule) AppAutoMigrate() string {
	ret := filepath.Join(p.rootPath, "app", "entity", "autoMigrate.go")
	return ret
}

// database.yaml tpl
func (p *pathForModule) DatabaseYamlTpl() string {
	ret := filepath.Join(p.rootPath, p.commandsPath, "enable", "tpl", "gorm", "database.yaml")
	return ret
}

// database.yaml
func (p *pathForModule) DatabaseYaml() string {
	ret := filepath.Join(p.rootPath, "app", "config", "database.yaml")
	return ret
}

// database-alpha.yaml tpl
func (p *pathForModule) DatabaseAlphaYamlTpl() string {
	ret := filepath.Join(p.rootPath, p.commandsPath, "enable", "tpl", "gorm", "database-alpha.yaml")
	return ret
}

// database-alpha.yaml
func (p *pathForModule) DatabaseAlphaYaml() string {
	ret := filepath.Join(p.rootPath, "app", "config", "alpha", "database.yaml")
	return ret
}

// database-release.yaml tpl
func (p *pathForModule) DatabaseReleaseYamlTpl() string {
	ret := filepath.Join(p.rootPath, p.commandsPath, "enable", "tpl", "gorm", "database-release.yaml")
	return ret
}

// database-release.yaml
func (p *pathForModule) DatabaseReleaseYaml() string {
	ret := filepath.Join(p.rootPath, "app", "config", "release", "database-release.yaml")
	return ret
}

// 模块 api 目录
func (p *pathForModule) ModuleApiDir() string {
	ret := filepath.Join(p.moduleDir, "internal", "api")
	return ret
}

// 模块 service 目录
func (p *pathForModule) ModuleServiceDir() string {
	ret := filepath.Join(p.moduleDir, "internal", "service")
	return ret
}

// 模块 dto 目录
func (p *pathForModule) ModuleDtoDir() string {
	ret := filepath.Join(p.moduleDir, "internal", "dto")
	return ret
}

// 模块 routting.go 文件
func (p *pathForModule) ModuleRoutingFile() string {
	ret := filepath.Join(p.moduleDir, "internal", "routing", "routing.go")
	return ret
}

// apidebug 目录
func (p *pathForModule) ModuleApiDebugDir() string {
	ret := filepath.Join(p.moduleDir, "apidebug")
	return ret
}

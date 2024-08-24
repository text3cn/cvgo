package paths

import (
	"cvgo/kvs"
	"path/filepath"
)

type ModulePath struct {
	ModulePath string
}

func NewModulePath(moduleName ...string) *ModulePath {
	var modulePath string
	if len(moduleName) > 0 {
		root := kvs.Instance().GetRootPath()
		modulePath = filepath.Join(root, "app", moduleName[0])
	}

	return &ModulePath{
		ModulePath: modulePath,
	}
}

func (p *ModulePath) IndeApiGo() string {
	apiFilePath := filepath.Join(p.ModulePath, "internal", "api", "index.api.go")
	return apiFilePath
}

func (p *ModulePath) IndeApiGitkeep() string {
	apiFilePath := filepath.Join(p.ModulePath, "internal", "api", ".gitkeep")
	return apiFilePath
}

func (p *ModulePath) BootInitGo() string {
	apiFilePath := filepath.Join(p.ModulePath, "internal", "boot", "init.go")
	return apiFilePath
}

func (p *ModulePath) BootInitGitkeep() string {
	apiFilePath := filepath.Join(p.ModulePath, "internal", "boot", ".gitkeep")
	return apiFilePath
}

func (p *ModulePath) BootSwaggerGo() string {
	apiFilePath := filepath.Join(p.ModulePath, "internal", "boot", "swagger.go")
	return apiFilePath
}

func (p *ModulePath) RoutingGo() string {
	return filepath.Join(p.ModulePath, "internal", "routing", "routing.go")
}

func (p *ModulePath) RoutingGitkeep() string {
	return filepath.Join(p.ModulePath, "internal", "routing", ".gitkeep")
}

func (p *ModulePath) ConfigAppYaml() string {
	return filepath.Join(p.ModulePath, "internal", "config", "app.yaml")
}

func (p *ModulePath) AuthMiddlewareGo() string {
	return filepath.Join(p.ModulePath, "internal", "middleware", "auth.go")
}

func (p *ModulePath) I18nMiddlewareGo() string {
	return filepath.Join(p.ModulePath, "internal", "middleware", "i18n.go")
}

func (p *ModulePath) MiddlewareGitkeep() string {
	return filepath.Join(p.ModulePath, "internal", "middleware", ".gitkeep")
}

//
//// instance.go
//func (p *pathForModule) InstanceGo() string {
//	ret := filepath.Join(p.rootPath, "app", "instance.go")
//	return ret
//}
//
//// boot -> init.go
//func (p *pathForModule) BootInitGo() string {
//	ret := filepath.Join(p.moduleDir, "internal", "boot", "init.go")
//	return ret
//}
//
//// app entity mysql
//func (p *pathForModule) AppEntityMysql() string {
//	ret := filepath.Join(p.rootPath, "app", "entity", "mysql")
//	return ret
//}
//

// 模块 api 目录
func (p *ModulePath) ModuleApiDir() string {
	ret := filepath.Join(p.ModulePath, "internal", "api")
	return ret
}

// 模块 service 目录
func (p *ModulePath) ModuleServiceDir() string {
	ret := filepath.Join(p.ModulePath, "internal", "service")
	return ret
}

// 模块 dto 目录
func (p *ModulePath) ModuleDtoDir() string {
	ret := filepath.Join(p.ModulePath, "internal", "dto")
	return ret
}

// apidebug 目录
func (p *ModulePath) ModuleApiDebugDir() string {
	ret := filepath.Join(p.ModulePath, "apidebug")
	return ret
}

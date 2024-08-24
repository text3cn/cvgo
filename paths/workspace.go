package paths

import (
	"cvgo/kvs"
	"path/filepath"
)

type WorkPath struct {
	workspacePath string
}

func NewWorkPath() *WorkPath {
	return &WorkPath{
		workspacePath: kvs.Instance().GetRootPath(),
	}
}

func (p *WorkPath) InstanceGo() string {
	return filepath.Join(p.workspacePath, "app", "instance.go")
}

func (p *WorkPath) AppEntityMysqlBaseGoFile() string {
	return filepath.Join(p.workspacePath, "entity", "mysql", "base.go")
}

// autoMigrate.go
func (p *WorkPath) AppAutoMigrate() string {
	return filepath.Join(p.workspacePath, "entity", "autoMigrate.go")
}

func (p *WorkPath) EntityRegistryTpl() string {
	return filepath.Join(p.workspacePath, "tpl", "enable", "gorm", "entity", "entityRegistry.go")
}

func (p *WorkPath) EntityRegistry() string {
	return filepath.Join(p.workspacePath, "scripts", "cvgo", "codegen", "entityregistry", "entityRegistry.go")
}

func (p *WorkPath) CurdGenScript() string {
	return filepath.Join(p.workspacePath, "scripts", "cvgo", "codegen", "curdl.go")
}

// database.yaml
func (p *WorkPath) DatabaseYaml() string {
	ret := filepath.Join(p.workspacePath, "config", "database.yaml")
	return ret
}

// database-alpha.yaml
func (p *WorkPath) DatabaseAlphaYaml() string {
	ret := filepath.Join(p.workspacePath, "config", "alpha", "database.yaml")
	return ret
}

// database-release.yaml
func (p *WorkPath) DatabaseReleaseYaml() string {
	ret := filepath.Join(p.workspacePath, "config", "release", "database-release.yaml")
	return ret
}

func (p *WorkPath) EntityMysqlDir() string {
	ret := filepath.Join(p.workspacePath, "entity", "mysql")
	return ret
}

func (p *WorkPath) DockerComposeEnv() string {
	ret := filepath.Join(p.workspacePath, "docker-compose-env.yml")
	return ret
}

func (p *WorkPath) DockerDir() string {
	ret := filepath.Join(p.workspacePath, "docker")
	return ret
}

func (p *WorkPath) ScriptsDir() string {
	ret := filepath.Join(p.workspacePath, "scripts")
	return ret
}

func (p *WorkPath) MigrateSwaggerScript() string {
	ret := filepath.Join(p.workspacePath, "scripts", "migrate_swagger.go")
	return ret
}

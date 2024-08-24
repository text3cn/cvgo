package paths

type CvgoPath struct {
	//cvgoPath string
}

func NewCvgoPath() *CvgoPath {
	return &CvgoPath{
		//	cvgoPath: kvs.Instance().GetRootPath(),
	}
}

func (p *CvgoPath) MysqlBaseEntityTpl() string {
	return "enable/gorm/entity/mysql/base.go"
}

// autoMigrate.go tpl
func (p *CvgoPath) AutoMigrateTpl() string {
	return "enable/gorm/entity/autoMigrate.go"
}

func (p *CvgoPath) EntiryRegistryTpl() string {
	return "enable/gorm/entityregistry/entityRegistry.go"
}

// database.yaml tpl
func (p *CvgoPath) DatabaseYamlTpl() string {
	return "enable/gorm/database.yaml"
}

// database-alpha.yaml tpl
func (p *CvgoPath) DatabaseAlphaYamlTpl() string {
	return "enable/gorm/database-alpha.yaml"
}

// database-release.yaml tpl
func (p *CvgoPath) DatabaseReleaseYamlTpl() string {
	return "enable/gorm/database-release.yaml"
}

func (p *CvgoPath) CurdGenScript() string {
	return "enable/gorm/gen_curdl.go.tpl"
}

func (p *CvgoPath) DockerComposeEnv() string {
	return "docker/docker-compose-env.yml"
}

func (p *CvgoPath) DockerDir() string {
	return "docker/docker"
}

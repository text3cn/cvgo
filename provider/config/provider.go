package config

import (
	"cvgo/provider/core"
	"os"
)

var instance *ConfigService

const Name = "config"

type ConfigProvider struct {
	core.ServiceProvider
}

func (self *ConfigProvider) Name() string {
	return Name
}

func (*ConfigProvider) InitOnBind() bool {
	return true
}

func (*ConfigProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (*ConfigProvider) RegisterProviderInstance(c core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		file, _ := os.Getwd()
		instance = &ConfigService{container: c, currentPath: file + "/"}
		return instance, nil
	}
}

func (*ConfigProvider) BeforeInit(c core.Container) error {
	//goodlog.Trace("BeforeInit Config Provider")
	return nil
}

func (*ConfigProvider) AfterInit(instance any) error {
	//goodlog.Trace("AfterInit Config Provider")
	return nil
}

package i18n

import (
	"cvgo/provider/core"
	"github.com/spf13/viper"
	"sync"
)

const Name = "i18n"

var instance *I18nService

type I18nProvider struct {
	core.ServiceProvider
}

func (self *I18nProvider) Name() string {
	return Name
}

func (*I18nProvider) InitOnBind() bool {
	return true
}

func (*I18nProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (*I18nProvider) RegisterProviderInstance(c core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		instance = &I18nService{c: c, pkgs: make(map[string]*viper.Viper), lock: sync.RWMutex{}}
		return instance, nil
	}
}

func (*I18nProvider) BeforeInit(c core.Container) error {
	return nil
}

func (*I18nProvider) AfterInit(instance any) error {
	return nil
}

package orm

import (
	plog2 "cvgo/provider/clog"
	"cvgo/provider/config"
	"cvgo/provider/core"
	"sync"
)

const Name = "orm"

var instance *OrmService

type OrmProvider struct {
	core.ServiceProvider
}

func (self *OrmProvider) Name() string {
	return Name
}

func (*OrmProvider) RegisterProviderInstance(c core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		instance = &OrmService{
			c: c, lock: sync.Mutex{},
			plog:   c.NewSingle(plog2.Name).(plog2.Service),
			cfgSvc: c.NewSingle(config.Name).(config.Service),
		}
		return instance, nil
	}
}

func (*OrmProvider) InitOnBind() bool {
	return false
}

func (*OrmProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (*OrmProvider) BeforeInit(c core.Container) error {
	return nil
}

func (*OrmProvider) AfterInit(instance any) error {
	return nil
}

package etcd

import (
	"cvgo/provider"
	"cvgo/provider/core"
)

const Name = "etcd"

var etcdInstance *etcdService

type ReidsProvider struct {
	core.ServiceProvider
}

func (self *ReidsProvider) Name() string {
	return Name
}

func (*ReidsProvider) RegisterProviderInstance(box core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		etcdInstance = &etcdService{c: box}
		return etcdInstance, nil
	}
}

func (*ReidsProvider) InitOnBind() bool {
	return false
}

func (*ReidsProvider) BeforeInit(c core.Container) error {
	provider.Clog().Trace("BeforeInit Etcd Provider")
	return nil
}

func (*ReidsProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (*ReidsProvider) AfterInit(instance any) error {
	return nil
}

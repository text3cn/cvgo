// 实现服务中心规定的服务注册要求，遵循注册协议 engine.Container
package localcache

import (
	"cvgo/provider/core"
	"fmt"
)

const Name = "localcache"

var instance *CacheService

type LocalCacheProvider struct {
	core.ServiceProvider
}

func (self *LocalCacheProvider) Name() string {
	return Name
}

func (sp *LocalCacheProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (*LocalCacheProvider) InitOnBind() bool {
	return false
}

func (sp *LocalCacheProvider) RegisterProviderInstance(c core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		instance = &CacheService{holder: c}
		return instance, nil
	}
}

func (sp *LocalCacheProvider) BeforeInit(c core.Container) error {
	fmt.Println("BeforeInit Cache Provider")
	return nil
}

func (*LocalCacheProvider) AfterInit(instance any) error {
	return nil
}

package redis

import (
	"cvgo/provider/config"
	"cvgo/provider/core"
	"sync"
)

const Name = "redis"

var instance *RedisService

type RedisProvider struct {
	core.ServiceProvider
}

func (self *RedisProvider) Name() string {
	return Name
}

func (*RedisProvider) RegisterProviderInstance(c core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		instance = &RedisService{
			c:      c,
			lock:   sync.Mutex{},
			cfgSvc: c.NewSingle(config.Name).(config.Service),
		}
		return instance, nil
	}
}

func (*RedisProvider) InitOnBind() bool {
	return false
}

func (*RedisProvider) BeforeInit(c core.Container) error {
	return nil
}

func (*RedisProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (*RedisProvider) AfterInit(instance any) error {
	return nil
}

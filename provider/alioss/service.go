package alioss

import (
	"cvgo/provider/config"
	"cvgo/provider/core"
	"cvgo/provider/core/types"
	"github.com/redis/go-redis/v9" // https://redis.uptrace.dev/zh/guide
	"github.com/spf13/cast"
	"sync"
)

type Service interface {
	init()
	GetConnPool(connName ...string) *redis.Client
}

// android：https://help.aliyun.com/zh/oss/developer-reference/authorize-access?spm=a2c4g.11186623.0.0.34485e0fRLKejO
// sts：https://help.aliyun.com/zh/oss/developer-reference/use-temporary-access-credentials-provided-by-sts-to-access-oss?spm=a2c4g.11186623.0.0.6c783f8cb2tU2S#section-7tz-fgu-oji

type AliossService struct {
	Service
	c           core.Container
	clients     map[string]*redis.Client
	lock        sync.Mutex
	cfgSvc      config.Service
	redisConfig map[string]types.RedisConfig
}

// 如果能获取到配置文件则进行连接，
// 上锁防止调用端在协程中并发初始化时出错
func (self *AliossService) init() {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.redisConfig = self.cfgSvc.GetRedis()
	for connName, configItem := range self.redisConfig {
		rdb := redis.NewClient(&redis.Options{
			Addr:     configItem.Host + ":" + cast.ToString(configItem.Port),
			Password: configItem.Auth,
			DB:       configItem.Db,
		})
		// 挂到 map 中
		if self.clients == nil {
			self.clients = make(map[string]*redis.Client)
		}
		self.clients[connName] = rdb
	}
}

func (self *AliossService) GetConnPool(connName ...string) *redis.Client {
	if self.clients == nil {
		self.init()
	}
	var key string
	if len(connName) > 0 {
		key = connName[0]
	} else {
		for _, v := range self.redisConfig {
			key = v.DefaultConn
			break
		}
	}
	return self.clients[key]
}

package provider

import (
	"cvgo/provider/config"
	"cvgo/provider/core"
	"cvgo/provider/i18n"
	"cvgo/provider/localcache"
	"cvgo/provider/orm"
	"cvgo/provider/plog"
	"cvgo/provider/redis"
)

var Services = core.NewContainer()
var Plog plog.Service

func init() {
	Services.Bind(&config.ConfigProvider{})
	Services.Bind(&plog.PlogProvider{})
	Services.Bind(&orm.OrmProvider{})
	Services.Bind(&redis.RedisProvider{})
	Services.Bind(&i18n.I18nProvider{})
	Services.Bind(&localcache.LocalCacheProvider{})

	Plog = Services.NewSingle(plog.Name).(plog.Service)
}

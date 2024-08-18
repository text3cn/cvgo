package provider

import (
	"cvgo/provider/clog"
	"cvgo/provider/config"
	"cvgo/provider/core"
	"cvgo/provider/i18n"
	"cvgo/provider/localcache"
	"cvgo/provider/orm"
	"cvgo/provider/redis"
)

var Services = core.NewContainer()
var log clog.Service

func init() {
	Services.Bind(&config.ConfigProvider{})
	Services.Bind(&clog.ClogProvider{})
	Services.Bind(&orm.OrmProvider{})
	Services.Bind(&redis.RedisProvider{})
	Services.Bind(&i18n.I18nProvider{})
	Services.Bind(&localcache.LocalCacheProvider{})

	log = Services.NewSingle(clog.Name).(clog.Service)
}

func Clog() clog.Service {
	return log
}

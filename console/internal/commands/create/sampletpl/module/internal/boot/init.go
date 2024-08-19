package boot

import (
	"cvgo/app"
	"cvgo/provider"
	"cvgo/provider/clog"
	"cvgo/provider/config"
)

var HttpServerPort string

func init() {

	app.Config = provider.Services.NewSingle(config.Name).(config.Service)
	HttpServerPort = app.Config.GetHttpPort()

	app.Log = provider.Services.NewSingle(clog.Name).(clog.Service)

}

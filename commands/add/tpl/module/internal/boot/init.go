package boot

import (
	"cvgo/app"
	"github.com/textthree/provider"
	"github.com/textthree/provider/clog"
	"github.com/textthree/provider/config"
)

var HttpServerPort string

func init() {

	app.Config = provider.Services.NewSingle(config.Name).(config.Service)
	HttpServerPort = app.Config.GetHttpPort()

	app.Log = provider.Services.NewSingle(clog.Name).(clog.Service)

}

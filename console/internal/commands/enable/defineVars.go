package enable

import (
	"cvgo/provider"
	"cvgo/provider/clog"
)

var httpPort string
var routingFile string
var log = provider.Services.NewSingle(clog.Name).(clog.Service)
var pwd string

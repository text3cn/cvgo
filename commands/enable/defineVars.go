package enable

import (
	"github.com/textthree/provider"
	"github.com/textthree/provider/clog"
)

var routingFile string
var log = provider.Services.NewSingle(clog.Name).(clog.Service)
var pwd string

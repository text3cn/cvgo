package console

import (
	"github.com/textthree/provider"
	"github.com/textthree/provider/config"
)

func GetHttpPort() string {
	cfg := provider.Services.NewSingle(config.Name).(config.Service)
	return cfg.GetHttpPort()
}

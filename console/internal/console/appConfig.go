package console

import (
	"cvgo/provider"
	"cvgo/provider/config"
)

func GetHttpPort() string {
	cfg := provider.Services.NewSingle(config.Name).(config.Service)
	return cfg.GetHttpPort()
}

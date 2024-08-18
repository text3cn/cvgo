package clog

import (
	"cvgo/provider/config"
	"cvgo/provider/core"
	"strings"
)

const Name = "clog"

type ClogProvider struct {
	core.ServiceProvider // 显示的写上实现了哪个接口主要是为了代码可读性以及 IDE 友好
	level                byte
}

func (self *ClogProvider) Name() string {
	return Name
}

// 日志服务不需要延迟初始化，启动程序就需要打印日志了
func (*ClogProvider) InitOnBind() bool {
	return true
}

// 往服务中心注册自己前的操作
func (self *ClogProvider) BeforeInit(c core.Container) error {
	var level byte
	configSvs := c.NewSingle(config.Name).(config.Service)
	config := configSvs.GetCLog()

	switch strings.ToLower(config.Level) {
	case "trace":
		level = 0
	case "debug":
		level = 1
	case "info":
		level = 2
	case "warn":
		level = 3
	case "error":
		level = 4
	case "fatal":
		level = 5
	case "off":
		level = 6
	default:
		level = 0
	}
	self.level = level
	return nil
}

func (sp *ClogProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (self *ClogProvider) RegisterProviderInstance(c core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		// 这里需要将参数展开，将配置注入到日志类，例如日志开关等
		//c := params[0].(core.Container)
		//if clogSvc != nil {
		//	return clogSvc, nil
		//}
		clogSvc = &ClogService{c: c, level: self.level}
		return clogSvc, nil
	}

}

func (*ClogProvider) AfterInit(instance any) error {
	return nil
}

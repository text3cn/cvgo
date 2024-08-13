package i18n

import (
	"cvgo/provider/core"
	"github.com/spf13/viper"
	"sync"
)

//var pkgs = make(map[string]*viper.Viper)

type Service interface {
	SetLngCode(string)
	GetLngCode() string
	LoadedPackage(string) bool
	SetLanguagePackage(string, *viper.Viper)
	Get(string) string
}

type I18nService struct {
	Service // 实现接口，显示标记
	c       core.Container
	lngCode string
	pkgs    map[string]*viper.Viper
	lock    sync.RWMutex
}

func (self *I18nService) SetLngCode(lng string) {
	self.lngCode = lng
}

func (self *I18nService) GetLngCode() string {
	return self.lngCode
}

// 判断语言包是否已加载到内存
func (self *I18nService) LoadedPackage(lngCode string) bool {
	return self.pkgs[lngCode] != nil
}

// 从磁盘加载到内存。
func (self *I18nService) SetLanguagePackage(lng string, pkg *viper.Viper) {
	self.lock.Lock()
	self.lock.Unlock()
	self.pkgs[lng] = pkg
}

func (self *I18nService) Get(key string) string {
	return self.pkgs[self.lngCode].GetString(key)
}

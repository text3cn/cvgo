package config

import (
	"cvgo/kit/castkit"
	"cvgo/kit/filekit"
	"cvgo/kit/strkit"
	"cvgo/provider/core"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func Instance() *ConfigService {
	if instance == nil {
		file, _ := os.Getwd()
		instance = &ConfigService{currentPath: file + "/"}
	}
	return instance
}

type ConfigService struct {
	Service
	container   core.Container
	currentPath string       // 二进制 main 程序的绝对路径
	lock        sync.RWMutex // 配置文件读写锁
}

type Service interface {
	IsDebug() bool
	LoadConfig(filename string) (*viper.Viper, error)
	Get(key string) *castkit.GoodleVal
	GetHttpPort() string
	GetRuntimePath() string
	GetDatabase() (dbsCfg map[string]core.DBConfig)
	GetRedis() map[string]core.RedisConfig
	GetPLog() core.Plog
	SetCurrentPath(path string)
}

// 设置(篡改)当前工作路径，以便特殊路径在运行程序时按规则找配置文件。
func (self *ConfigService) SetCurrentPath(path string) {
	self.currentPath = path
}

func (self *ConfigService) LoadConfig(filename string) (*viper.Viper, error) {
	// fmt.Println("self.currentPath: ", self.currentPath)
	if configs[filename] != nil {
		return configs[filename], nil
	}
	var retConfig, commonConfig, appConfig, internalConfig, localConfig *viper.Viper
	seg := strings.Split(filename, ".")
	fName := seg[0]
	fType := seg[1]
	cfgFile := fName + "." + fType

	// 获取公共配置
	var commonConfigPath string
	parentDir := filepath.Dir(self.currentPath)
	parentDir = filepath.Dir(parentDir)
	parentDir = filepath.Dir(parentDir)
	commonConfigPath = filepath.Join(parentDir, "config")
	if exists, _ := filekit.PathExists(commonConfigPath); exists {
		commonConfig = loadConfigFile(commonConfigPath, fName, fType)
		retConfig = commonConfig
	}

	// 在可执行文件当前目录找配置文件
	if exists, _ := filekit.PathExists(filepath.Join(self.currentPath, cfgFile)); exists {
		appConfig = loadConfigFile(self.currentPath, fName, fType)
	}
	// 用 app 级覆盖 common 级配项
	if appConfig != nil {
		if retConfig == nil {
			retConfig = appConfig
		} else {
			allKeys := appConfig.AllKeys()
			for _, v := range allKeys {
				retConfig.Set(v, appConfig.Get(v))
			}
		}
	}

	// 在 ./internal/config 目录中找
	path := filepath.Join(self.currentPath, "internal", "config")
	file := filepath.Join(path, fName+"."+fType)
	if exists, _ := filekit.PathExists(file); exists {
		internalConfig = loadConfigFile(path, fName, fType)
	}
	// 再用 internalConfig 覆盖
	if internalConfig != nil {
		if retConfig == nil {
			retConfig = internalConfig
		} else {
			allKeys := internalConfig.AllKeys()
			for _, v := range allKeys {
				retConfig.Set(v, internalConfig.Get(v))
			}
		}
	}

	// 获取 local 级配置
	localConfigDir := filepath.Join(filepath.Dir(self.currentPath), "internal", "config", "local")
	localConfigFile := filepath.Join(localConfigDir, cfgFile)
	if exists, _ := filekit.PathExists(localConfigFile); exists {
		localConfig = loadConfigFile(localConfigDir, fName, fType)
		// 再用 local 级的配置项覆盖
		if retConfig == nil {
			retConfig = localConfig
		} else {
			allKeys := localConfig.AllKeys()
			for _, v := range allKeys {
				retConfig.Set(v, localConfig.Get(v))
			}
		}
	}

	// 没有配置文件
	if retConfig == nil {
		err := errors.New("Unable to find configuration file " + filename +
			" The configuration file should be placed in any of the following paths:\n" +
			self.currentPath + filename + "\n" +
			self.currentPath + "internal/config/" + filename + "\n" +
			self.currentPath + "internal/config/local/" + filename + "\n" +
			commonConfigPath + "/" + filename,
		)
		return nil, err
	}

	// 缓存起来，不用每次读硬盘
	configs[fName] = retConfig
	return retConfig, nil
}

func (self *ConfigService) Get(key string) *castkit.GoodleVal {
	seg := strkit.Explode(".", key)
	if len(seg) == 1 {
		return &castkit.GoodleVal{}
	}
	cfg := configs[seg[0]]
	itemKey := strings.Replace(key, seg[0]+".", "", 1)
	return &castkit.GoodleVal{cfg.Get(itemKey)}
}

func (self *ConfigService) IsDebug() bool {
	key := "debug"
	if cfg, _ := self.getDefaultConfig(); cfg != nil {
		if cfg.IsSet(key) {
			if val, ok := cfg.Get(key).(bool); !ok {
				panic("The configuration of " + key + " is not a valid value")
			} else {
				return cast.ToBool(val)
			}
		}
	}
	return false
}

// http 服务监听段口
func (self *ConfigService) GetHttpPort() (port string) {
	key := "server.http-port"
	if cfg, _ := self.getDefaultConfig(); cfg != nil {
		if value, ok := cfg.Get(key).(int); !ok {
			panic("The configuration of " + key + " is not a valid value")
		} else {
			port = cast.ToString(value)
		}
	}
	return
}

// runtime 目录
func (self *ConfigService) GetRuntimePath() string {
	key := "runtime.path"
	if cfg, _ := self.getDefaultConfig(); cfg != nil {
		if cfg.IsSet(key) {
			if val, ok := cfg.Get(key).(string); !ok {
				panic("The configuration of " + key + " is not a valid value")
			} else {
				return cast.ToString(val)
			}
		}
	}
	return self.currentPath
}

func (self *ConfigService) GetDatabase() (dbsCfg map[string]core.DBConfig) {
	dbsCfg = make(map[string]core.DBConfig)
	cfg, err := self.LoadConfig("database.yaml")
	if err != nil {
		panic(err)
		return
	}
	cfgNodes := mergerLevel2(cfg)
	for k, v := range cfgNodes {
		item := core.DBConfig{}
		mapstructure.Decode(v, &item)
		dbsCfg[k] = item
	}
	return
}

func (self *ConfigService) GetRedis() (configs map[string]core.RedisConfig) {
	configs = make(map[string]core.RedisConfig)
	cfg, _ := self.LoadConfig("redis.yaml")
	if cfg != nil {
		cfgNodes := mergerLevel2(cfg)
		for k, v := range cfgNodes {
			item := core.RedisConfig{}
			mapstructure.Decode(v, &item)
			configs[k] = item
		}
	}
	return
}

func (self *ConfigService) GetPLog() (config core.Plog) {
	if cfg, _ := self.getDefaultConfig(); cfg != nil {
		value := cfg.Get("plog")
		mapstructure.Decode(value, &config)
	}
	return
}

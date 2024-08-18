package config

import (
	"cvgo/kit/castkit"
	"cvgo/kit/filekit"
	"cvgo/kit/strkit"
	"cvgo/provider/core"
	"cvgo/provider/core/types"
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
	GetDatabase() (dbsCfg map[string]types.DBConfig)
	GetRedis() map[string]types.RedisConfig
	GetCLog() types.Clog
	SetCurrentPath(path string)
	GetSwagger() types.SwaggerConfig
	GetEtcd() types.EtcdConfig
	GetFileServer() types.FileSeverConfig
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
	commonConfigDir := filepath.Join(parentDir, "config")
	commonConfigFile := filepath.Join(parentDir, "config", filename)
	if exists, _ := filekit.PathExists(commonConfigFile); exists {
		commonConfig = loadConfigFile(commonConfigDir, fName, fType)
		retConfig = commonConfig
	}

	// 在可执行文件当前目录找配置文件（生产部署时通常最顶级配置为当前目录）
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
	if cfg, _ := self.getAppConfig(); cfg != nil {
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
	if cfg, _ := self.getAppConfig(); cfg != nil {
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
	if cfg, _ := self.getAppConfig(); cfg != nil {
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

func (self *ConfigService) GetDatabase() (dbsCfg map[string]types.DBConfig) {
	dbsCfg = make(map[string]types.DBConfig)
	cfg, err := self.LoadConfig("database.yaml")
	if err != nil {
		panic(err)
		return
	}
	cfgNodes := mergerLevel2(cfg)
	for k, v := range cfgNodes {
		item := types.DBConfig{}
		mapstructure.Decode(v, &item)
		dbsCfg[k] = item
	}
	return
}

func (self *ConfigService) GetRedis() (configs map[string]types.RedisConfig) {
	configs = make(map[string]types.RedisConfig)
	cfg, _ := self.LoadConfig("redis.yaml")
	if cfg != nil {
		cfgNodes := mergerLevel2(cfg)
		for k, v := range cfgNodes {
			item := types.RedisConfig{}
			mapstructure.Decode(v, &item)
			configs[k] = item
		}
	}
	return
}

func (self *ConfigService) GetCLog() (config types.Clog) {
	if cfg, _ := self.getAppConfig(); cfg != nil {
		value := cfg.Get("clog")
		mapstructure.Decode(value, &config)
	} else {
		config.Level = "trace" // 默认输出所有级别日志
	}
	return
}

func (self *ConfigService) GetSwagger() (config types.SwaggerConfig) {
	if cfg, _ := self.getAppConfig(); cfg != nil {
		value := cfg.Get("swagger")
		mapstructure.Decode(value, &config)
	}
	return
}

func (self *ConfigService) GetEtcd() (config types.EtcdConfig) {
	if cfg, _ := self.getAppConfig(); cfg != nil {
		value := cfg.Get("etcd")
		mapstructure.Decode(value, &config)
	}
	return
}

func (self *ConfigService) GetFileServer() (config types.FileSeverConfig) {
	if cfg, _ := self.getAppConfig(); cfg != nil {
		value := cfg.Get("fileServer")
		mapstructure.Decode(value, &config)
	}
	return
}

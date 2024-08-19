package console

import (
	"cvgo/provider"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type kvStorage struct {
	viper    *viper.Viper
	rootPath string
}

func NewKvStorage(rootPath string) kvStorage {
	// 初始化 Viper
	instance := kvStorage{
		viper:    viper.New(),
		rootPath: rootPath,
	}
	dataFileDir := filepath.Join(rootPath, "console")
	instance.viper.AddConfigPath(dataFileDir)
	instance.viper.SetConfigName("runtime") // 配置文件名称（不带扩展名）
	instance.viper.SetConfigType("json")    // 配置文件类型
	instance.viper.AddConfigPath(".")       // 配置文件路径
	return instance
}

// 保存配置到文件
func (this kvStorage) saveData() error {
	dataFile := filepath.Join(this.rootPath, "console", "runtime.json")

	// 打开文件，如果不存在则创建
	file, err := os.OpenFile(dataFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入配置
	err = this.viper.WriteConfigAs(dataFile)
	if err != nil {
		return err
	}

	return nil
}

func (this kvStorage) Set(key string, value any) error {
	// 读取旧配置
	if err := this.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			provider.Clog().Error("读取配置文件时发生错误: %v\n", err)
			return err
		}
	}
	this.viper.Set(key, value)
	if err := this.saveData(); err != nil {
		provider.Clog().Error("保存配置文件时发生错误: %v\n", err)
		return err
	}
	return nil
}

func (this kvStorage) GetBool(key string) (bool, error) {
	if err := this.viper.ReadInConfig(); err != nil {
		return false, errors.New(fmt.Sprintf("读取配置文件时发生错误: %v\n", err))
	}
	return this.viper.GetBool(key), nil
}

func (this kvStorage) GetString(key string) (string, error) {
	if err := this.viper.ReadInConfig(); err != nil {
		return "", errors.New(fmt.Sprintf("读取配置文件时发生错误: %v\n", err))
	}
	return this.viper.GetString(key), nil
}

func (this kvStorage) GetInt(key string) (int, error) {
	if err := this.viper.ReadInConfig(); err != nil {
		return DefaultHttpPort, errors.New(fmt.Sprintf("读取配置文件时发生错误: %v\n", err))
	}
	port := this.viper.GetInt(key)
	if port == 0 {
		port = DefaultHttpPort
	}
	return port, nil
}

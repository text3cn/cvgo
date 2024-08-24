package kvs

import (
	"cvgo/config"
	"cvgo/kvs/kvsKey"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/cvgokit/gokit"
	"os"
	"path/filepath"
)

type KvStorage struct {
	viper         *viper.Viper
	workspacePath string
	workspaceName string
	dataFile      string
}

func Instance(workspacePath ...string) KvStorage {
	var rootPath, workspaceName string
	if len(workspacePath) > 0 {
		rootPath = workspacePath[0]
	} else {
		rootPath, workspaceName = findWorkspacePath()
	}
	dataFileDir := filepath.Join(rootPath, "./scripts/cvgo")
	filekit.EnsureDirExists(dataFileDir)
	instance := KvStorage{
		viper:         viper.New(),
		workspacePath: rootPath,
		workspaceName: workspaceName,
		dataFile:      filepath.Join(dataFileDir, "cvg.json"),
	}
	instance.viper.AddConfigPath(dataFileDir)
	instance.viper.SetConfigName("cvg")  // 配置文件名称（不带扩展名）
	instance.viper.SetConfigType("json") // 配置文件类型
	instance.viper.AddConfigPath(".")    // 配置文件路径

	return instance
}

func findWorkspacePath() (modPath, modName string) {
	found := false
	findGoWork := func(path string) bool {
		modName, _ = gokit.GetModuleName()
		if exists, _ := filekit.PathExists(filepath.Join(path, "go.work")); exists {
			found = true
			return true
		}
		return false
	}
	modPath = filekit.Getwd()
	orignPath := modPath
	for i := 0; i < 7; i++ {
		if !findGoWork(modPath) {
			modPath = filepath.Dir(modPath)
			os.Chdir(modPath)
		} else {
			break
		}
	}
	if !found {
		panic(errors.New("找不到 go.work 文件"))
	}
	os.Chdir(orignPath)
	return
}

func (this KvStorage) findModule() (modPath, modName string) {
	found := false
	workspaceName := this.workspaceName
	findGoWork := func(path string) bool {
		modName, _ = gokit.GetModuleName()
		//clog.GreenPrintln("在", modPath, "找模块 go.mod", "modName="+modName, "workspaceName="+workspaceName)
		if exists, _ := filekit.PathExists(filepath.Join(path, "go.mod")); exists && (modName != workspaceName && modName != "") {
			found = true
			return true
		}
		return false
	}
	modPath = filekit.Getwd()
	for i := 0; i < 5; i++ {
		if !findGoWork(modPath) {
			modPath = filepath.Dir(modPath)
		} else {
			break
		}
	}
	if !found {
		panic(errors.New("请在模块目录下执行此命令"))
	}
	return
}

// Save the configuration to a file
func (this KvStorage) saveData() error {
	// Open the file and create it if it does not exist
	file, err := os.OpenFile(this.dataFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write configuration
	err = this.viper.WriteConfigAs(this.dataFile)
	if err != nil {
		return err
	}

	return nil
}

// 判断是否在模块目录下
func (this KvStorage) CheckInModuleDir() {
	this.findModule() // 找不到 go.mod 会 panic
}

func (this KvStorage) Set(key string, value any) {
	if err := this.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			filekit.FilePutContents(this.dataFile, "{}")
		}
	}
	this.viper.Set(key, value)
	if err := this.saveData(); err != nil {
		panic(err)
	}
}

//func (this KvStorage) SetProjectRootPath() {
//	this.Set("projectRootPath", this.rootPath)
//}

func (this KvStorage) GetRootPath() string {
	if err := this.viper.ReadInConfig(); err != nil {
		panic(err)
	}
	return this.viper.GetString(kvsKey.WorkspacePath)
}

func (this KvStorage) GetBool(key string) (bool, error) {
	if err := this.viper.ReadInConfig(); err != nil {
		return false, errors.New(fmt.Sprintf("读取配置文件时发生错误: %v\n", err))
	}
	return this.viper.GetBool(key), nil
}

func (this KvStorage) GetString(key string) (string, error) {
	if err := this.viper.ReadInConfig(); err != nil {
		return "", errors.New(fmt.Sprintf("读取配置文件时发生错误: %v\n", err))
	}
	return this.viper.GetString(key), nil
}

func (this KvStorage) GetStringSlice(key string) ([]string, error) {
	if err := this.viper.ReadInConfig(); err != nil {
		return nil, errors.New(fmt.Sprintf("读取配置文件时发生错误: %v\n", err))
	}
	return this.viper.GetStringSlice(key), nil
}

func (this KvStorage) GetAllocatedPort() (int, error) {
	if err := this.viper.ReadInConfig(); err != nil {
		return config.DefaultHttpPort, errors.New(fmt.Sprintf("读取配置文件时发生错误: %v\n", err))
	}
	port := this.viper.GetInt(kvsKey.AllocatedPort)
	if port == 0 {
		port = config.DefaultHttpPort
	}
	return port, nil
}

// 获取使用的 web 框架类型
func (this KvStorage) GetWebFramework() string {
	_, modName := this.findModule()
	ret, err := this.GetString(kvsKey.ModuleWebFramework(modName))
	if err != nil {
		panic(err)
	}
	return ret
}

// 获取是否支持 swagger
func (this KvStorage) GetSwagger() (bool, error) {
	modName, _ := gokit.GetModuleName()
	ret, _ := this.GetBool(kvsKey.ModuleSwaggerEnable(modName))
	return ret, nil
}

// 获取工作区路径
func (this KvStorage) GetWorkspacePath() string {
	ret, _ := this.GetString(kvsKey.WorkspacePath)
	return ret
}

func (this KvStorage) GetWorkspaceName() string {
	return this.workspaceName
}

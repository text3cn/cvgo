package enable

import (
	"cvgo/console/internal/console"
	"cvgo/console/internal/paths"
	"cvgo/kit/filekit"
	"cvgo/kit/gokit"
	"cvgo/provider/clog"
)

// 在模块目录中执行
// cd ../../../ && go build -o $GOPATH/bin/cvg ./console && cd app/modules/fiber && cvg add mysql
func addGorm() {
	if !paths.CheckRunAtModuleRoot() {
		return
	}
	modName, err := gokit.GetModuleName()
	if err != nil {
		panic(err)
	}
	kv := console.NewKvStorage(filekit.GetParentDir(3))
	mysqlStorageKey := modName + "." + "mysql"
	commonEntityInited := false // 公共代码是否已经初始化
	if val, _ := kv.GetBool(mysqlStorageKey); val {
		log.Info("mysql 已经添加过了，无法重复执行。")
		return
	}

	path := paths.NewPathForModule()

	// 添加 base.go 到 entity
	err = filekit.CopyFile(path.MysqlBaseEntityTpl(), path.AppEntityMysqlBaseGoFile())
	if err != nil {
		log.Error(err)
		commonEntityInited = true
	}

	// 添加自动迁移
	err = filekit.CopyFile(path.AutoMigrateTpl(), path.AppAutoMigrate())
	if err != nil {
		log.Error(err)
		commonEntityInited = true
	}

	if !commonEntityInited {
		// instance.go 中定义变量
		filekit.AddContentUnderLine(path.InstanceGo(), "import (", `    "gorm.io/gorm"`)
		filekit.FileAppendContent(path.InstanceGo(), "\n"+`var Db *gorm.DB`)

		// boot -> init.go 中获取实例
		content := `    "cvgo/app/entity"`
		filekit.AddContentUnderLine(path.BootInitGo(), "import (", content)
		content = `
	// 获取 mysql 连接池
	// 文档：https://gorm.io/zh_CN/docs
	database := provider.Services.NewSingle(orm.Name).(orm.Service)
	app.Db = database.GetConnPool()
	if !app.IsDevelop() {
		entity.AutoMigrate(app.Db) // 生产环境自动迁移表结构
	}
`
		filekit.AddContentUnderLine(path.BootInitGo(), "func init() {", content)
	}

	// 拷贝配置文件
	err = filekit.CopyFile(path.DatabaseYamlTpl(), path.DatabaseYaml())
	if err != nil {
		log.Error(err)
	}
	err = filekit.CopyFile(path.DatabaseAlphaYamlTpl(), path.DatabaseAlphaYaml())
	if err != nil {
		log.Error(err)
	}
	err = filekit.CopyFile(path.DatabaseReleaseYamlTpl(), path.DatabaseReleaseYaml())
	if err != nil {
		log.Error(err)
	}

	// 标识已添加  mysql
	kv.Set(mysqlStorageKey, true)
	clog.GreenPrintln("添加 mysql 支持成功")
}

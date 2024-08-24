package enable

import (
	"cvgo/kvs"
	"cvgo/paths"
	"cvgo/tpl"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/cvgokit/gokit"
	"github.com/textthree/provider/clog"
)

func addMysql() {
	paths.CheckRunAtModuleRoot()

	modName, err := gokit.GetModuleName()
	if err != nil {
		panic(err)
	}
	kv := kvs.Instance()
	mysqlStorageKey := modName + "." + "mysql"
	copyTplSuccess := true // 公共代码是否已经初始化
	if val, _ := kv.GetBool(mysqlStorageKey); val {
		log.Info("mysql 已经添加过了，无法重复执行。")
		return
	}

	modulePath := paths.NewModulePath()
	workPath := paths.NewWorkPath()
	cvgoPath := paths.NewCvgoPath()

	// 添加 base.go 到 entity
	err = tpl.CopyFileFromEmbed(tpl.EntityBase, cvgoPath.MysqlBaseEntityTpl(), workPath.AppEntityMysqlBaseGoFile())
	if err != nil {
		log.Error(err)
		copyTplSuccess = false
	}

	// 添加自动迁移
	err = tpl.CopyFileFromEmbed(tpl.AutoMigrate, cvgoPath.AutoMigrateTpl(), workPath.AppAutoMigrate())
	if err != nil {
		log.Error(err)
		copyTplSuccess = false
	}

	// 添加 entity 注册表到 scripts
	err = tpl.CopyFileFromEmbed(tpl.EntityRegistry, cvgoPath.EntiryRegistryTpl(), workPath.EntityRegistry())
	if err != nil {
		log.Error(err)
		copyTplSuccess = false
	}

	// 拷贝配置文件
	err = tpl.CopyFileFromEmbed(tpl.Database, cvgoPath.DatabaseYamlTpl(), workPath.DatabaseYaml())
	if err != nil {
		log.Error(err)
	}
	err = tpl.CopyFileFromEmbed(tpl.DatabaseAlpha, cvgoPath.DatabaseAlphaYamlTpl(), workPath.DatabaseAlphaYaml())
	if err != nil {
		log.Error(err)
	}
	err = tpl.CopyFileFromEmbed(tpl.DatabaseRelease, cvgoPath.DatabaseReleaseYamlTpl(), workPath.DatabaseReleaseYaml())
	if err != nil {
		log.Error(err)
	}

	// CURD 代码生成脚本
	err = tpl.CopyFileFromEmbed(tpl.CurdGen, cvgoPath.CurdGenScript(), workPath.CurdGenScript())
	if err != nil {
		log.Error(err)
	}

	//return

	if copyTplSuccess {
		// instance.go 中定义变量
		filekit.AddContentUnderLine(workPath.InstanceGo(), "import (", `    "gorm.io/gorm"`)
		filekit.FileAppendContent(workPath.InstanceGo(), "\n"+`var Db *gorm.DB`)

		// boot -> init.go 中获取实例
		content := `    "` + kv.GetWorkspaceName() + `/entity"
	"github.com/textthree/provider/orm"
`
		filekit.AddContentUnderLine(modulePath.BootInitGo(), "import (", content)
		content = `
	// 获取 mysql 连接池
	// 文档：https://gorm.io/zh_CN/docs
	database := provider.Services.NewSingle(orm.Name).(orm.Service)
	app.Db = database.GetConnPool()
	if !app.IsDevelop() {
		entity.AutoMigrate(app.Db) // 生产环境自动迁移表结构
	}
`
		filekit.AddContentUnderLine(modulePath.BootInitGo(), "func init() {", content)
	}

	// 标识已添加  mysql
	kv.Set(mysqlStorageKey, true)
	clog.GreenPrintln("添加 mysql 支持成功")
}

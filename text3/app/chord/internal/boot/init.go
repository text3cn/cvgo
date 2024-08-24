package boot

import (
	"github.com/textthree/provider/orm"
	"text3/entity"

	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/provider"
	"github.com/textthree/provider/clog"
	"github.com/textthree/provider/config"
	"text3/app"
)

func init() {

	// 获取 mysql 连接池
	// 文档：https://gorm.io/zh_CN/docs
	database := provider.Services.NewSingle(orm.Name).(orm.Service)
	app.Db = database.GetConnPool()
	if !app.IsDevelop() {
		entity.AutoMigrate(app.Db) // 生产环境自动迁移表结构
	}

	app.Config = provider.Services.NewSingle(config.Name).(config.Service)
	app.Log = provider.Services.NewSingle(clog.Name).(clog.Service)
	clog.CyanPrintln("  Current Path: " + filekit.Getwd())
}

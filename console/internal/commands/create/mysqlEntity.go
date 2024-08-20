package create

import (
	"cvgo/console/internal/paths"
	"cvgo/kit/filekit"
	"cvgo/kit/strkit"
	"path/filepath"
	"strings"
)

// 创建 msyql entity，在工程根目录执行：
// cd ../../../ && go build -o $GOPATH/bin/cvg ./console && cd app/modules/chord && cvg create mysqlEntity user 用户表
func createMysqlEntity(tableName, comment string) {
	path := paths.NewPathForApp()
	if !paths.CheckRunAtModuleRoot() {
		log.Error("请在模块根目录下执行此命令")
		return
	}
	entityFile := filepath.Join(path.AppEntityMysql(), tableName+".ent.go")
	if exists, _ := filekit.PathExists(entityFile); exists {
		log.Error(entityFile + " 已经存在，无法创建。")
		return
	}

	entityName := strkit.Ucfirst(tableName) + "Entity"
	content := `package mysql

import "gorm.io/gorm"

// 模型定义文档：https://gorm.io/zh_CN/docs/models.html
type ` + entityName + ` struct {
	CommonField
}

// ` + comment + `
func (` + entityName + `) TableName() string {
	return "` + tableName + `"
}

// 添加表注释
func (this ` + entityName + `) SetTableComment(db *gorm.DB) {
	db.Exec("ALTER TABLE ` + tableName + ` COMMENT = '` + comment + `'")
}
`
	filekit.FilePutContents(entityFile, content)

	// 加入自动迁移，执行语句添加表注释
	content = `
	entities = append(entities, &mysql.` + entityName + `{})
	mysql.` + entityName + `{}.SetTableComment(db)`

	err := filekit.AddContentUnderLine(path.AppAutoMigrate(), "var entities = []interface{}{}", content)
	if err != nil {
		log.Error("无法在 "+path.AppAutoMigrate(), "中找到 var entities = []interface{}{} 这行代码，因此无法将添加自动迁移代码。")
	}

	// 如果还没导包还需要导包
	autoMigrateContent, _ := filekit.FileGetContents(path.AppAutoMigrate())
	if !strings.Contains(autoMigrateContent, "cvgo/app/entity/mysql") {
		content = `    "cvgo/app/entity/mysql"`
		filekit.AddContentUnderLine(path.AppAutoMigrate(), "import (", content)
	}

}

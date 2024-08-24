package add

import (
	"cvgo/commands/add/addgencode"
	"cvgo/console/internal/paths"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/cvgokit/strkit"
	"os"
	"path/filepath"
)

// 创建 msyql entity，可在工程下任意路径执行
// go build -o $GOPATH/bin/cvg ./console && cvg add table user_article
func createMysqlEntity(tableName, comment string) {
	if !paths.CheckRunAtProjectRoot() {
		err := os.Chdir(paths.GetProjectRootPathFromKv())
		if err != nil {
			panic(err)
		}

	}

	path := paths.NewPathForProjectRoot()
	entityFile := filepath.Join(path.AppEntityMysql(), tableName+".ent.go")
	if exists, _ := filekit.PathExists(entityFile); exists {
		log.Error(entityFile + " 已经存在，无法创建。")
		return
	}

	entityName := strkit.SnakeToPascalCase(tableName) + "Entity"
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
	content = `    entities = append(entities, &mysql.` + entityName + `{})`
	err := filekit.AddContentUnderLine(path.AppAutoMigrate(), "var entities = []interface{}{}", content)
	if err != nil {
		log.Error("无法在 "+path.AppAutoMigrate(), "中找到 var entities = []interface{}{} 这行代码，因此无法将添加自动迁移代码。")
	}
	content = `    mysql.` + entityName + `{}.SetTableComment(db)`
	err = filekit.AddContentUnderLine(path.AppAutoMigrate(), "func addTableComment(db *gorm.DB) {", content)
	if err != nil {
		panic(err)
	}
	// 加入全局 EntityRegistry
	content = `    EntityRegistry["` + tableName + `"] =  reflect.TypeOf(mysql.` + entityName + `{})`
	err = filekit.AddContentUnderLine(path.AppAutoMigrate(), "func init() {", content)
	if err != nil {
		panic(err)
	}

	// 如果还没导包还需要导包
	addgencode.ImportPackageIfNotImport(path.AppAutoMigrate(), "cvgo/app/entity/mysql")
}

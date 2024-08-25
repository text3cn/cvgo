package scripts

import (
	"cvgo/kvs"
	"cvgo/paths"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/provider/clog"
	"path/filepath"
)

func CreateScript(script string) {
	paths.CdToWorkspacePath()
	switch script {
	case "migrate_swagger":
		migrateAndSwagger()
	}
}

// go build -o $GOPATH/bin/cvg
// cvg codegen script migrate_swagger
func migrateAndSwagger() {
	workPath := paths.NewWorkPath()
	workspaceName := kvs.Instance().GetWorkspaceName()
	content := `package main

import (
	"bytes"
	"` + workspaceName + `/entity"
	"fmt"
	"github.com/textthree/provider"
	"github.com/textthree/provider/config"
	"github.com/textthree/provider/orm"
	"os"
	"os/exec"
	"path/filepath"
)

var log = provider.Clog()

// 生成 swagger 文档、gorm 自动迁移
// 在工程根目录下执行：go run scripts/migrate_swagger.go
func main() {
	mysqlMigrate("client")
	docGenerate("client")
}

// 迁移表结构
func mysqlMigrate(moduleName string) {
	currentDir, _ := os.Getwd()
	cfg := provider.Svc().NewSingle(config.Name).(config.Service)
	cfg.SetCurrentPath(filepath.Join(currentDir, "app", moduleName) + string(os.PathSeparator))
	database := provider.Svc().NewSingle(orm.Name).(orm.Service)
	conn := database.GetConnPool()
	entity.AutoMigrate(conn)
}

// 使用 swag init 命令生成 swagger 文档
func docGenerate(app string) {
	currentDir, err := os.Getwd()
	sep := string(os.PathSeparator)
	os.Chdir(currentDir + sep + "app" + sep + app)
	_, err = syskit.ExecWithOutput("swag", "init", "--parseDependency", "--propertyStrategy", "pascalcase")
	if err != nil {
		pwd, _ := os.Getwd()
		log.Error("生成 Swagger 文档出错。请在 "+pwd+" 目录下手动执行 swag init --parseDependency 命令查看错误信息。", err)
		return
	}
	os.Chdir(currentDir)
}
`
	err := filekit.FilePutContents(workPath.MigrateSwaggerScript(), content)
	if err != nil {
		clog.RedPrintln(err)
	}
	filekit.DeleteFile(filepath.Join(workPath.ScriptsDir(), ".gitkeep"))
}

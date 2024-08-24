package main

import (
	"bytes"
	"fmt"
	"github.com/textthree/provider"
	"github.com/textthree/provider/config"
	"github.com/textthree/provider/orm"
	"os"
	"os/exec"
	"path/filepath"
	"text3/entity"
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
	cfg := provider.Services.NewSingle(config.Name).(config.Service)
	cfg.SetCurrentPath(filepath.Join(currentDir, "app", "modules", moduleName) + string(os.PathSeparator))
	database := provider.Services.NewSingle(orm.Name).(orm.Service)
	conn := database.GetConnPool()
	entity.AutoMigrate(conn)
}

// 使用 swag init 命令生成 swagger 文档
func docGenerate(app string) {
	currentDir, err := os.Getwd()
	sep := string(os.PathSeparator)
	os.Chdir(currentDir + sep + "app" + sep + "modules" + sep + app)
	cmd := exec.Command("swag", "init", "--parseDependency", "--propertyStrategy", "pascalcase")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		pwd, _ := os.Getwd()
		log.Error("生成 Swagger 文档出错。请在 " + pwd + " 目录下手动执行 swag init --parseDependency 命令查看错误信息。")
		return
	}
	fmt.Println(out.String())
	os.Chdir(currentDir)
}

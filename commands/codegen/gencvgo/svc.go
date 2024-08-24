package gencvgo

import (
	"cvgo/commands/codegen/common"
	"cvgo/kvs"
	"cvgo/kvs/kvsKey"
	"cvgo/paths"
	"fmt"
	"github.com/textthree/cvgokit/arrkit"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/cvgokit/gokit"
	"github.com/textthree/cvgokit/strkit"
	"github.com/textthree/cvgokit/syskit"
	"github.com/textthree/provider/clog"
	"os"
	"path/filepath"
	"strings"
)

// 模块目录下执行：
// go build -o $GOPATH/bin/cvg
// cvg add svc user/Userinfo c --table=user --cursor
func GenService(fileName, funcName string, curdType string, tableName string, cursorPaging bool) {
	kv := kvs.Instance()
	paths.CheckRunAtModuleRoot()
	modulePath := paths.NewModulePath()
	modName, _ := gokit.GetModuleName()
	kvKey := kvsKey.ModuleSvc(modName)
	fileAndFunc := fileName + "/" + funcName
	oldSvcs, _ := kv.GetStringSlice(kvKey)
	if arrkit.InArray(fileAndFunc, oldSvcs) {
		clog.RedPrintln("Service", fileAndFunc, "已存在")
		return
	}

	// 创建 service
	fileNameLower := strings.ToLower(fileName)
	fileNamePascalCase := strkit.Ucfirst(fileName)
	svcFile := filepath.Join(modulePath.ModuleServiceDir(), fileNameLower+".svc.go")
	err := filekit.CreatePath(svcFile)
	if err == nil {
		content := `package service
	
import (
	"github.com/textthree/cvgoweb"
	"sync"
)
	
var ` + fileNameLower + `ServiceInstance *` + fileNamePascalCase + `Service
var ` + fileNameLower + `ServiceOnce sync.Once
	
type ` + fileNamePascalCase + `Service struct {
	ctx *httpserver.Context
	uid int64
}
	
func ` + fileNamePascalCase + `Svc(ctx *httpserver.Context) *` + fileNamePascalCase + `Service {
	` + fileNameLower + `ServiceOnce.Do(func() {
	` + fileNameLower + `ServiceInstance = &` + fileNamePascalCase + `Service{
			ctx: ctx,
			uid: ctx.GetVal("uid").ToInt64(),
		}
	})
	return ` + fileNameLower + `ServiceInstance
}
`
		filekit.FilePutContents(svcFile, content)
		filekit.DeleteFile(filepath.Join(modulePath.ModuleServiceDir(), ".gitkeep"))
	}
	var content string
	if curdType == "" {
		content = `
	// ` + funcName + `
	func (self *` + fileNamePascalCase + `Service) ` + funcName + `() {
	
	}
	`
	} else {
		content = createFuncWithCurd(curdType, tableName, svcFile, cursorPaging, funcName, fileNamePascalCase)
	}
	err = filekit.FileAppendContent(svcFile, content)
	if err != nil {
		clog.RedPrintln(err)
	}
	// 完成
	kv.Set(kvKey, append(oldSvcs, fileAndFunc))
	clog.GreenPrintln("生成方法 " + fileNameLower + ".svc.go -> " + funcName)
}

func createFuncWithCurd(curdType, tableName, svcFile string, cursorPaging bool, funcName, fileNamePascalCase string) string {
	code := ""
	originPath := filekit.Getwd()
	paths.CdToWorkspacePath()
	workspaceName := kvs.Instance().GetWorkspaceName()
	switch curdType {
	case "c":
		//	code = common.CurdCreate(tableName, funcName, fileNamePascalCase)
		code = syskit.ExecCmdText(fmt.Sprintf("go run scripts/cvgo/codegen/curdl.go c %s %s %s %s", tableName, funcName, fileNamePascalCase, ""))
	case "u":
		code = syskit.ExecCmdText(fmt.Sprintf("go run scripts/cvgo/codegen/curdl.go u %s %s %s %s", tableName, funcName, fileNamePascalCase, ""))
	case "r":
		code = syskit.ExecCmdText(fmt.Sprintf("go run scripts/cvgo/codegen/curdl.go r %s %s %s %s", tableName, funcName, fileNamePascalCase, ""))
	case "d":
		code = syskit.ExecCmdText(fmt.Sprintf("go run scripts/cvgo/codegen/curdl.go d %s %s %s %s", tableName, funcName, fileNamePascalCase, ""))
	case "l":
		code = syskit.ExecCmdText(fmt.Sprintf("go run scripts/cvgo/codegen/curdl.go l %s %s %s %v", tableName, funcName, fileNamePascalCase, cursorPaging))
	}
	os.Chdir(originPath)

	// app 包与 msyql 包，如果还没导包还需要导包
	common.ImportPackageIfNotImport(svcFile, workspaceName+"/app")
	common.ImportPackageIfNotImport(svcFile, workspaceName+"/entity/mysql")
	return code
}

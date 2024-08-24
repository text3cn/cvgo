package gencvgo

import (
	"cvgo/console/internal/commands/add/addgencode"
	"cvgo/console/internal/console"
	"cvgo/console/internal/paths"
	"cvgo/kit/arrkit"
	"cvgo/kit/filekit"
	"cvgo/kit/gokit"
	"cvgo/kit/strkit"
	"cvgo/provider/clog"
	"fmt"
	"path/filepath"
	"strings"
)

// 模块目录下执行：
// cd ../../../ && go build -o $GOPATH/bin/cvg ./console && cd app/modules/chord && cvg add svc user/Userinfo l --table=user --cursor
func GenService(fileName, funcName string, curdType string, tableName string, cursorPaging bool) {
	path := paths.NewPathForModule()
	modName, _ := gokit.GetModuleName()
	kv := console.NewKvStorage(filekit.GetParentDir(3))
	kvKey := modName + ".services"
	fileAndFunc := fileName + "/" + funcName
	oldSvcs, _ := kv.GetStringSlice(kvKey)
	if arrkit.InArray(fileAndFunc, oldSvcs) {
		clog.RedPrintln("Service", fileAndFunc, "已存在")
		return
	}

	// 创建 service
	fileNameLower := strings.ToLower(fileName)
	fileNamePascalCase := strkit.Ucfirst(fileName)
	svcFile := filepath.Join(path.ModuleServiceDir(), fileNameLower+".svc.go")
	err := filekit.CreatePath(svcFile)
	if err == nil {
		content := `package service

import (
	"cvgo/provider/httpserver"
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
		filekit.DeleteFile(filepath.Join(path.ModuleServiceDir(), ".gitkeep"))
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
	filekit.FileAppendContent(svcFile, content)
	// 完成
	err = kv.Set(kvKey, append(oldSvcs, fileAndFunc))
	if err != nil {
		fmt.Println(err)
	}
	clog.GreenPrintln("生成 Service 成功")
}

func createFuncWithCurd(curdType, tableName, svcFile string, cursorPaging bool, funcName, fileNamePascalCase string) string {
	code := ""
	switch curdType {
	case "c":
		code = addgencode.CurdCreate(tableName, funcName, fileNamePascalCase)
	case "u":
		code = addgencode.CurdUpdate(tableName, funcName, fileNamePascalCase)
	case "r":
		code = addgencode.CurdGet(tableName, funcName, fileNamePascalCase)
	case "d":
		code = addgencode.CurdDelete(tableName, funcName, fileNamePascalCase)
	case "l":
		code = addgencode.CurdList(tableName, funcName, fileNamePascalCase, cursorPaging)
	}
	// app 包与 msyql 包，如果还没导包还需要导包
	addgencode.ImportPackageIfNotImport(svcFile, "cvgo/app")
	addgencode.ImportPackageIfNotImport(svcFile, "cvgo/app/entity/mysql")
	return code
}

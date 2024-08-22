package gencvgo

import (
	"cvgo/app/entity"
	"cvgo/console/internal/console"
	"cvgo/console/internal/paths"
	"cvgo/kit/arrkit"
	"cvgo/kit/filekit"
	"cvgo/kit/gokit"
	"cvgo/kit/strkit"
	"cvgo/provider/clog"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// 模块目录下执行：
// cd ../../../ && go build -o $GOPATH/bin/cvg ./console && cd app/modules/chord && cvg add svc user/GetUserinfo
func GenService(fileName, funcName string, curdType string, tableName string) {
	path := paths.NewPathForModule()
	modName, _ := gokit.GetModuleName()
	kv := console.NewKvStorage(filekit.GetParentDir(3))
	kvKey := modName + ".services"
	fileAndFunc := fileName + "/" + funcName
	oldSvcs, _ := kv.GetStringSlice(kvKey)
	if arrkit.InArray(fileAndFunc, oldSvcs) {
		clog.RedPrintln("Service", fileAndFunc, "已存在")
		//return
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

	if curdType == "" {

		content := `
// ` + funcName + `
func (self *` + fileNamePascalCase + `Service) ` + funcName + `() {
 
}
`
		filekit.FileAppendContent(svcFile, content)
	} else {
		createFuncWithCurd(tableName)
	}
	// 完成
	err = kv.Set(kvKey, append(oldSvcs, fileAndFunc))
	if err != nil {
		fmt.Println(err)
	}
	clog.GreenPrintln("生成 Service 成功")
}

func createFuncWithCurd(tableName string) {

}

// Create 基础代码
func create(tableName string) string {
	entityType, found := entity.EntityRegistry[tableName]
	if !found {
		clog.RedPrintln("类型", strkit.SnakeToPascalCase(tableName), "不存在")
		return ""
	}
	// 创建实体实例
	entityValue := reflect.New(entityType).Elem()
	content := generateGormCreateCode(entityValue, "mysql", "UserEntity")
	fmt.Println(content)
	return content
}

// 根据 modelValue 生成 Gorm 的创建代码
func generateGormCreateCode(entityValue reflect.Value, packageName, structName string) string {
	modelType := entityValue.Type()
	code := fmt.Sprintf("%s.Db.Create(&%s.%s{\n", "app", packageName, structName)

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		value := entityValue.Field(i)

		// 获取字段默认值的字符串表示形式
		var fieldValue string
		switch value.Kind() {
		case reflect.String:
			fieldValue = `""`
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldValue = "0"
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldValue = "0"
		case reflect.Float32, reflect.Float64:
			fieldValue = "0.0"
		case reflect.Bool:
			fieldValue = "false"
		case reflect.Struct:
			// 检查是否为 time.Time 类型
			if field.Type == reflect.TypeOf(time.Time{}) {
				fieldValue = "time.Time{}"
			} else {
				fieldValue = fmt.Sprintf("%s{}", field.Type.Name())
			}
		default:
			fieldValue = "nil"
		}

		code += fmt.Sprintf("\t%s: %s,\n", field.Name, fieldValue)
	}

	code += "})"

	return code
}
